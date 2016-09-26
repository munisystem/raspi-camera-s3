package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type AWSConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	Bucket          string
}

type Camera struct {
	Width  string
	Height string
}

func NewAWSConfig() *AWSConfig {
	return &AWSConfig{
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		Region:          os.Getenv("AWS_DEFAULT_REGION"),
		Bucket:          os.Getenv("AWS_S3_BUCKET"),
	}
}

func (conf *AWSConfig) UploadToS3(path, key string) {
	file, err := os.Open(path)
	if err != nil {
		log.Println(err.Error())
	}
	defer file.Close()

	cli := s3.New(session.New(), &aws.Config{
		Credentials: credentials.NewStaticCredentials(conf.AccessKeyID, conf.SecretAccessKey, ""),
		Region:      aws.String(conf.Region),
	})

	_, err = cli.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(conf.Bucket),
		Key:    aws.String(key),
		Body:   file,
	})

	if err != nil {
		log.Println(err.Error())
	}
}

func NewCamera() *Camera {
	var err error

	var width string
	sw := os.Getenv("IMAGE_WIDTH")
	if sw == "" {
		width = "1920"
	} else {
		if _, err = strconv.Atoi(sw); err != nil {
			log.Fatal("$IMAGE_WIDTH must be integer.")
		}
	}

	var height string
	sh := os.Getenv("IMAGE_HEIGHT")
	if sh == "" {
		height = "1080"
	} else {
		if _, err = strconv.Atoi(sh); err != nil {
			log.Fatal("$IMAGE_HEIGHT must be integer.")
		}
	}

	return &Camera{
		Width:  width,
		Height: height,
	}
}

func (c *Camera) CaptureStillFrame(path string) {
	err := exec.Command("raspistill", "-o", path, "-w", c.Width, "-h", c.Height, "-vf").Run()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// REF: http://stackoverflow.com/questions/20234104/how-to-format-current-time-using-a-yyyymmddhhmmss-format
	t := time.Now()
	ft := t.Format("20060102150405")

	// yyyyMMdd
	s3Dir := ft[0:8]

	dir, err := ioutil.TempDir("", "images")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	name := ft + "jpg"
	path := filepath.Join(dir, name)

	// Generate image file
	c := NewCamera()
	c.CaptureStillFrame(path)

	conf := NewAWSConfig()
	key := filepath.Join(s3Dir, name)
	conf.UploadToS3(path, key)
}

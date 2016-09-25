package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Config struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	Bucket          string
}

func NewConfig() *Config {
	return &Config{
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		Region:          os.Getenv("AWS_DEFAULT_REGION"),
		Bucket:          os.Getenv("AWS_S3_BUCKET"),
	}
}

func (conf *Config) UploadToS3(path, key string) {
	file, err := os.Open(path)
	if err != nil {
		log.Println(err.Error())
	}
	defer file.Close()

	cli := s3.New(session.New(), &aws.Config{
		Credentials: credentials.NewStaticCredentials(conf.AccessKeyID, conf.SecretAccessKey, ""),
		Region:      aws.String(conf.Region),
	})

	resp, err := cli.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(conf.Bucket),
		Key:    aws.String(key),
		Body:   file,
	})

	if err != nil {
		log.Println(err.Error())
	}

	log.Println(awsutil.StringValue(resp))
}

func main() {
}

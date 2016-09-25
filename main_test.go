package main

import (
	"path/filepath"
	"testing"
)

func TestUploadToS3(t *testing.T) {
	conf := NewConfig()
	path := filepath.Join("test", "god.jpg")
	conf.UploadToS3(path, path)
}

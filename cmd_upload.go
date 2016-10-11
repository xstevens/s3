package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	uploadCommand   = kingpin.Command("upload", "Upload file(s) to S3")
	uploadBucket    = uploadCommand.Flag("bucket", "S3 bucket").Required().String()
	uploadPrefix    = uploadCommand.Flag("prefix", "S3 key prefix").Required().String()
	uploadSourceDir = uploadCommand.Flag("sourcedir", "Source directory").Required().String()
)

func runUpload() {
	sourceDir, err := os.Stat(*uploadSourceDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	files, err := ioutil.ReadDir(*uploadSourceDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	for _, fileInfo := range files {
		sourcePath := fmt.Sprintf("%s/%s", *uploadSourceDir, fileInfo.Name())
		destKey := fmt.Sprintf("%s/%s/%s", *uploadPrefix, sourceDir.Name(), fileInfo.Name())
		fmt.Printf("Uploading %s to %s ... ", sourcePath, destKey)
		file, err := os.Open(sourcePath)
		if err != nil {
			fmt.Printf("Failed to open file: %q\n", err)
			break
		}
		defer file.Close()

		uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String(*region)}))
		_, err = uploader.Upload(&s3manager.UploadInput{
			Body:                 bufio.NewReader(file),
			Bucket:               aws.String(*uploadBucket),
			Key:                  aws.String(destKey),
			ServerSideEncryption: aws.String(s3.ServerSideEncryptionAes256),
		})
		if err != nil {
			fmt.Printf("[ FAILED ]\n\t=> %q\n", err)
		} else {
			fmt.Println("[ SUCCESS ]")
		}
	}
}

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"gopkg.in/alecthomas/kingpin.v2"
)

type UploadCommand struct {
	Region    string
	Bucket    string
	Prefix    string
	SourceDir string
	MFASerial string
	RoleARN   string
}

func configureUploadCommand(app *kingpin.Application) {
	uc := &UploadCommand{}
	upload := app.Command("upload", "Upload file(s) to S3").Action(uc.runUpload)
	upload.Flag("region", "S3 region").Default("us-east-1").StringVar(&uc.Region)
	upload.Flag("bucket", "S3 bucket").Required().StringVar(&uc.Bucket)
	upload.Flag("prefix", "S3 prefix").Required().StringVar(&uc.Prefix)
	upload.Flag("sourcedir", "Source directory").Required().StringVar(&uc.SourceDir)
	upload.Flag("serial", "IAM MFA device ARN").StringVar(&uc.MFASerial)
	upload.Flag("role", "IAM Role ARN to assume").StringVar(&uc.RoleARN)
}

func (uc *UploadCommand) runUpload(ctx *kingpin.ParseContext) error {
	sourceDir, err := os.Stat(uc.SourceDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	files, err := ioutil.ReadDir(uc.SourceDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	for _, fileInfo := range files {
		sourcePath := fmt.Sprintf("%s/%s", uc.SourceDir, fileInfo.Name())
		destKey := fmt.Sprintf("%s/%s/%s", uc.Prefix, sourceDir.Name(), fileInfo.Name())
		fmt.Printf("Uploading %s to %s ... ", sourcePath, destKey)
		file, err := os.Open(sourcePath)
		if err != nil {
			fmt.Printf("Failed to open file: %q\n", err)
			break
		}
		defer file.Close()

		config := aws.NewConfig().WithRegion(uc.Region)
		sess, err := newSession(config, &uc.MFASerial, &uc.RoleARN)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open session: %v\n", err.Error())
			return err
		}

		uploader := s3manager.NewUploader(sess)
		_, err = uploader.Upload(&s3manager.UploadInput{
			Body:                 bufio.NewReader(file),
			Bucket:               aws.String(uc.Bucket),
			Key:                  aws.String(destKey),
			ServerSideEncryption: aws.String(s3.ServerSideEncryptionAes256),
		})
		if err != nil {
			fmt.Printf("[ FAILED ]\n\t=> %q\n", err)
		} else {
			fmt.Println("[ SUCCESS ]")
		}
	}

	return nil
}

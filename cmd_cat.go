package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

type CatCommand struct {
	Region    string
	Bucket    string
	Prefix    string
	MFASerial string
	RoleARN   string
}

func configureCatCommand(app *kingpin.Application) {
	cc := &CatCommand{}
	cat := app.Command("cat", "Reads all keys with specified prefix and writes them to stdout.").Action(cc.runCat)
	cat.Flag("region", "S3 region").Default("us-east-1").StringVar(&cc.Region)
	cat.Flag("bucket", "S3 bucket").Required().StringVar(&cc.Bucket)
	cat.Flag("prefix", "S3 prefix").Required().StringVar(&cc.Prefix)
	cat.Flag("serial", "IAM MFA device ARN").StringVar(&cc.MFASerial)
	cat.Flag("role", "IAM Role ARN to assume").StringVar(&cc.RoleARN)
}

func (cc *CatCommand) runCat(ctx *kingpin.ParseContext) error {
	config := aws.NewConfig().WithRegion(cc.Region)
	sess, err := newSession(config, &cc.MFASerial, &cc.RoleARN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open session: %v\n", err.Error())
		return err
	}
	s3Client := s3.New(sess, config)

	objKeys := make([]string, 0, 10000)
	listParams := &s3.ListObjectsV2Input{
		Bucket:  aws.String(cc.Bucket),
		Prefix:  aws.String(cc.Prefix),
		MaxKeys: aws.Int64(1000),
	}

	// this allows us to see errors that otherwise are hidden by list object pages
	_, err = s3Client.ListObjectsV2(listParams)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list objects: %v\n", err.Error())
		return err
	}

	s3Client.ListObjectsV2Pages(listParams, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range page.Contents {
			objKeys = append(objKeys, *obj.Key)
		}
		return !lastPage
	})

	for _, k := range objKeys {
		getParams := &s3.GetObjectInput{
			Bucket: aws.String(cc.Bucket),
			Key:    aws.String(k),
		}
		getResp, err := s3Client.GetObject(getParams)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get object with key: %s %v\n", k, err.Error())
			continue
		}

		var reader io.ReadCloser
		if strings.HasSuffix(k, ".gz") || strings.HasSuffix(k, ".tgz") {
			reader, err = gzip.NewReader(getResp.Body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to initialize a gzip reader: %v\n", err.Error())
			}
		} else {
			reader = getResp.Body
		}
		defer reader.Close()

		fmt.Fprintf(os.Stderr, "Reading key: %s\n", k)
		data, err := ioutil.ReadAll(reader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read object with key: %s %v\n", k, err.Error())
			continue
		}

		fmt.Println(string(data))
	}
	// downloader := s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String(*region)}))
	// buf := make([]byte, 1024*1024*4)
	// wabuf := aws.NewWriteAtBuffer(buf)
	// for _, k := range objKeys {
	// 	fmt.Println(k)
	// 	input := &s3.GetObjectInput{
	// 		Bucket: catBucket,
	// 		Key:    aws.String(k),
	// 	}

	// 	n, err := downloader.Download(wabuf, input)
	// 	if err != nil {
	// 		fmt.Fprintf(os.Stderr, "Download failure: %v\n", err.Error())
	// 	}
	// 	fmt.Fprintf(os.Stderr, "Bytes read: %d\n", n)

	// 	if strings.HasSuffix(k, ".gz") {
	// 		reader, err := gzip.NewReader(bytes.NewReader(wabuf.Bytes()))
	// 		if err != nil {
	// 			fmt.Fprintf(os.Stderr, "Failed to initialize gzip reader: %v\n", err.Error())
	// 			continue
	// 		}
	// 		defer reader.Close()

	// 		data, err := ioutil.ReadAll(reader)
	// 		if err != nil {
	// 			fmt.Fprintf(os.Stderr, "Failed to read object with key: %s, %v\n", k, err.Error())
	// 			continue
	// 		}
	// 		fmt.Println(string(data))
	// 	} else {
	// 		fmt.Println(string(wabuf.Bytes()))
	// 	}
	// }

	return nil
}

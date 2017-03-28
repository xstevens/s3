package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

type MetaCommand struct {
	Region    string
	Bucket    string
	Prefix    string
	MFASerial string
	RoleARN   string
}

func configureMetaCommand(app *kingpin.Application) {
	mc := &MetaCommand{}
	meta := app.Command("meta", "Reads all keys with specified prefix and writes their metadata to stdout.").Action(mc.runMeta)
	meta.Flag("region", "S3 region").Default("us-east-1").StringVar(&mc.Region)
	meta.Flag("bucket", "S3 bucket").Required().StringVar(&mc.Bucket)
	meta.Flag("prefix", "S3 prefix").Required().StringVar(&mc.Prefix)
	meta.Flag("serial", "IAM MFA device ARN").StringVar(&mc.MFASerial)
	meta.Flag("role", "IAM Role ARN to assume").StringVar(&mc.RoleARN)
}

func (cc *MetaCommand) runMeta(ctx *kingpin.ParseContext) error {
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
	s3Client.ListObjectsV2Pages(listParams, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range page.Contents {
			objKeys = append(objKeys, *obj.Key)
			// the code below does not get all the metadata we need, such as
			// SSE, ReplicationStatus, etc.
			// objJson, err := json.Marshal(obj)
			// if err != nil {
			// 	fmt.Printf("JSON serialization error: %v\n", err.Error())
			// 	return true
			// }
			// fmt.Println(string(objJson))
		}
		return !lastPage
	})

	// unfortunately we have to do a full get object for all the metadata we need
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

		meta := getResp.Metadata
		meta["Key"] = &k
		meta["ReplicationStatus"] = getResp.ReplicationStatus
		meta["ServerSideEncryption"] = getResp.ServerSideEncryption
		meta["SSECustomerAlgorithm"] = getResp.SSECustomerAlgorithm
		meta["SSECustomerKeyMD5"] = getResp.SSECustomerKeyMD5
		meta["SSEKMSKeyId"] = getResp.SSEKMSKeyId
		metaJson, err := json.Marshal(meta)
		if err != nil {
			fmt.Printf("JSON serialization error: %v\n", err.Error())
			return err
		}
		fmt.Println(string(metaJson))
	}

	return nil
}

package main

import (
	"github.com/CastyLab/cdn.service/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func GetS3Bucket() *s3.S3 {
	opts := aws.NewConfig()
	opts.WithEndpoint(config.Map.Secrets.ObjectStorage.Endpoint)
	opts.WithRegion(config.Map.Secrets.ObjectStorage.Region)
	opts.WithCredentials(credentials.NewStaticCredentialsFromCreds(credentials.Value{
		AccessKeyID: config.Map.Secrets.ObjectStorage.AccessKey,
		SecretAccessKey: config.Map.Secrets.ObjectStorage.SecretKey,
	}))
	sess := session.Must(session.NewSession(opts))
	return s3.New(sess)
}

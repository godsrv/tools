package awsS3

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
)

// @author: lipper
// @object: *s3.Client
// @function: NewS3Client
// @description: 实例s3
// @return: *s3.Client
func NewS3Client(ctx context.Context, conf AwsConf) *s3.Client {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: conf.Endpoint,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{AccessKeyID: conf.AccessKey, SecretAccessKey: conf.SecretKey},
		}),
	)
	if err != nil {
		logrus.Panicf("init s3 client with config: %+v err: %v", conf, err)
	}

	Client := s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.UsePathStyle = true
	})

	return Client
}

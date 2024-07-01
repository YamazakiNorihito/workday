package awsConfig

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/translate"
)

type SNSTopicClient struct {
	snsClient *sns.Client
	topicArn  string
}

func (c *AwsConfig) NewDynamodbClient() *dynamodb.Client {
	client := dynamodb.NewFromConfig(c.cfg, func(o *dynamodb.Options) {
		if endpoint := os.Getenv("DYNAMODB_ENDPOINT"); endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
	})
	return client
}

func (c *AwsConfig) NewS3Client() *s3.Client {
	client := s3.NewFromConfig(c.cfg, func(o *s3.Options) {
		o.Credentials = credentials.NewStaticCredentialsProvider("8o2RS265xUkhAQPsmpYy", "qpfrBNwoBSs92UtAMtblncGVsvQyrMyylWEjfHRo", "")
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(os.Getenv("S3_ENDPOINT"))
	})
	return client
}

func (c *AwsConfig) NewSnsClient() *sns.Client {
	client := sns.NewFromConfig(c.cfg, func(o *sns.Options) {
		if endpoint := os.Getenv("SNS_ENDPOINT"); endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
	})
	return client
}

func (c *AwsConfig) NewTranslateClient() *translate.Client {
	client := translate.NewFromConfig(c.cfg)
	return client
}

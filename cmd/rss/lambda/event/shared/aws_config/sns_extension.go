package awsConfig

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func NewSnsTopicClient(client *sns.Client, topicArn string) *SNSTopicClient {
	return &SNSTopicClient{client, topicArn}
}

func (c *SNSTopicClient) Publish(ctx context.Context, message string) error {
	pubInput := &sns.PublishInput{
		TopicArn: aws.String(c.topicArn),
		Message:  aws.String(message),
	}

	_, err := c.snsClient.Publish(ctx, pubInput)
	if err != nil {
		return err
	}
	return nil
}

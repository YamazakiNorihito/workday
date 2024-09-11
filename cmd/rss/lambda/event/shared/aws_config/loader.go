package awsConfig

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

func LoadConfig(ctx context.Context) AwsConfig {
	var optFns []func(*config.LoadOptions) error

	if region := os.Getenv("AWS_REGION_CODE"); region != "" {
		optFns = append(optFns, config.WithRegion(region))
	}

	if accessKey := os.Getenv("AWS_ACCESS_KEY_ID"); accessKey != "" {
		secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
		optFns = append(optFns, func(o *config.LoadOptions) error {
			o.Credentials = credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")
			return nil
		})
	}

	cfg, err := config.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		panic("unable to load AWS SDK config: " + err.Error())
	}
	return AwsConfig{cfg}
}

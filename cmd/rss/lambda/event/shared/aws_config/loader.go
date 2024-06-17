package awsConfig

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
)

func LoadConfig(ctx context.Context) AwsConfig {
	var optFns []func(*config.LoadOptions) error

	if region := os.Getenv("AWS_REGION_CODE"); region != "" {
		optFns = append(optFns, config.WithRegion(region))
	}

	cfg, err := config.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		panic("unable to load AWS SDK config: " + err.Error())
	}
	return AwsConfig{cfg}
}

package aws

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
)

const awsInitTimeout = time.Second * 10

var cw *cloudwatch.Client
var onceCW sync.Once

func CloudWatch() *cloudwatch.Client {
	if cw == nil {
		onceCW.Do(func() {
			log.Println("AWS Init")
			ctx, cancel := context.WithTimeout(context.Background(), awsInitTimeout)
			defer cancel()
			cw = cloudwatch.NewFromConfig(mustLoadAWSConfig(ctx))
		})
	}
	return cw
}

func mustLoadAWSConfig(ctx context.Context) aws.Config {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(fmt.Errorf("unable to load SDK config: %w ", err))
	}
	return cfg
}

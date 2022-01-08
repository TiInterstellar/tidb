package awss3

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pingcap/tidb/config"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAll(t *testing.T) {
	cfg := config.NewConfig()
	cfg.LoadEnv()
	cfg.Aws.Region = "ap-northeast-2"
	config.StoreGlobalConfig(cfg)
	cli, err := CreateS3Client()
	require.NoError(t, err)

	result, err := cli.ListBuckets(&s3.ListBucketsInput{})
	require.NoError(t, err)

	for _, bucket := range result.Buckets {
		DeleteBucketForTablePartition(cli, *bucket.Name)
	}
}

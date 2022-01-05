package remotestorage

import (
	"fmt"
	"github.com/pingcap/log"
	"github.com/pingcap/tidb/config"
	"go.uber.org/zap"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3 struct {
	client *s3.S3
	bucket string
}

func NewS3() *S3 {
	sess := session.Must(session.NewSession())

	// FIXME: we should check if s3 configured before using, skipped to make life easier.
	cfg := config.GetGlobalConfig().RemoteStorage.S3

	svc := s3.New(sess, &aws.Config{Credentials: credentials.NewStaticCredentials(cfg.AccessKeyId, cfg.SecretAccessKey, "")})

	return &S3{
		client: svc,
		bucket: cfg.Bucket,
	}
}

func (s *S3) Write(path string, size int64, r io.Reader) (err error) {
	_, err = s.client.PutObject(&s3.PutObjectInput{
		Body:          aws.ReadSeekCloser(r),
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(path),
		ContentLength: aws.Int64(size),
	})
	if err != nil {
		return fmt.Errorf("write s3: %w", err)
	}
	log.L().Info("write s3", zap.String("path", path))
	return nil
}

package server

import (
	"context"
	"strings"
	"time"

	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/models"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func (s *Server) initObjectStorage() {
	config := s.configFn().GetFile()
	path := "user.server.initObjectStorage"

	// create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	done := make(chan struct{})
	endpoint := strings.TrimPrefix(strings.TrimPrefix(config.GetAmazonS3Endpoint(), "http://"), "https://")

	go func() {
		client, err := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(config.GetAmazonS3AccessKeyId(), config.GetAmazonS3SecretAccessKey(), ""),
			Secure: config.GetAmazonS3Ssl(),
		})
		if err != nil {
			s.errors <- &models.InternalError{Err: err, Msg: "failed to initialize the object storage", Path: path}
			close(done)
			return
		}

		err = client.MakeBucket(ctx, config.GetAmazonS3Bucket(), minio.MakeBucketOptions{Region: config.GetAmazonS3Region()})
		if err != nil {
			exists, errBucket := client.BucketExists(ctx, config.GetAmazonS3Bucket())
			if errBucket == nil && exists {
				s.log.Infof("the bucket: %s, already exists", config.GetAmazonS3Bucket())
			} else {
				s.errors <- &models.InternalError{Err: errBucket, Msg: "failed to check if the bucket exists", Path: path}
			}
		} else {
			s.log.Infof("the bucket: %s, got created successfully!", config.GetAmazonS3Bucket())
		}

		s.objectStorage = client
		close(done)
	}()

	select {
	case <-ctx.Done():
		s.errors <- &models.InternalError{
			Err:  ctx.Err(),
			Msg:  "object storage initialization timed out",
			Path: path,
		}
	case <-done:
		s.log.Info("object storage initialization finished")
	}
}

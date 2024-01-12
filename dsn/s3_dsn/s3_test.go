package s3_dsn_test

import (
	"fmt"
	"testing"

	"github.com/geektheripper/vast-dsn/dsn/s3_dsn"
)

func TestS3DSNParser(t *testing.T) {
	t.Run("amazon s3", func(t *testing.T) {
		config, err := s3_dsn.ParseS3DSN("s3://access_key:secret_key@-?region=us-east-2")
		if err != nil {
			t.Error(err)
			return
		}

		cred, err := config.Credentials.Get()
		if err != nil {
			t.Error(err)
			return
		}

		if cred.AccessKeyID != "access_key" {
			t.Error("access key not parsed")
			return
		}

		if *config.Region != "us-east-2" {
			t.Error("region not parsed")
			return
		}
	})

	t.Run("self sign https minio", func(t *testing.T) {
		config, err := s3_dsn.ParseS3DSN("s3://access_key:secret_key@maggie.minio.geektr.co:9003?region=")
		if err != nil {
			t.Error(err)
			return
		}

		cred, err := config.Credentials.Get()
		if err != nil {
			t.Error(err)
			return
		}

		if cred.AccessKeyID != "access_key" {
			t.Error("access key not parsed")
			return
		}

		if *config.Region != "" {
			t.Error("region not parsed")
			return
		}

		if *config.Endpoint != "https://maggie.minio.geektr.co:9003" {
			t.Error("endpoint not parsed")
			return
		}
	})

	t.Run("s3 bucket dsn", func(t *testing.T) {
		config, bucket, err := s3_dsn.ParseS3BucketDSN("s3://access_key:secret_key@maggie.minio.geektr.co:9003/foobar/path/to/key?region=")
		if bucket != "" || config != nil {
			t.Error("wrong bucket dsn parsed")
			return
		}

		if fmt.Sprint(err) != fmt.Sprintf("invalid s3 bucket dsn: unexpected key: %s", "path/to/key") {
			t.Error(err)
			return
		}

		_, bucket2, _ := s3_dsn.ParseS3BucketDSN("s3://access_key:secret_key@maggie.minio.geektr.co:9003/foobar?region=")

		if bucket2 != "foobar" {
			t.Error("bucket not parsed")
			return
		}
	})

	t.Run("s3 object dsn", func(t *testing.T) {
		_, _, key, err := s3_dsn.ParseS3ObjectDSN("s3://access_key:secret_key@maggie.minio.geektr.co:9003/foobar/path/to/key?region=")
		if err != nil {
			t.Error(err)
			return
		}

		if key != "path/to/key" {
			t.Error("key not parsed")
			return
		}
	})
}

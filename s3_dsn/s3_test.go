package s3_dsn_test

import (
	"fmt"
	"testing"

	"github.com/geektheripper/vast-dsn/s3_dsn"
)

func TestS3DSNParser(t *testing.T) {
	t.Run("amazon s3", func(t *testing.T) {
		opts, err := s3_dsn.ParseS3("s3://access_key:secret_key@-?region=us-east-2")
		if err != nil {
			t.Error(err)
			return
		}

		if opts.Region != "us-east-2" {
			t.Error("region not parsed")
			return
		}

		if opts.AccessKeyID != "access_key" {
			t.Error("access key not parsed")
			return
		}
	})

	t.Run("self sign https minio", func(t *testing.T) {
		opts, err := s3_dsn.ParseS3("s3://access_key:secret_key@maggie.minio.geektr.co:9000?region=")
		if err != nil {
			t.Error(err)
			return
		}

		if opts.Region != "" {
			t.Error("region not parsed")
			return
		}

		if opts.Host != "maggie.minio.geektr.co:9000" {
			t.Error("endpoint not parsed")
			return
		}

		if opts.AccessKeyID != "access_key" {
			t.Error("access key not parsed")
			return
		}
	})

	t.Run("s3 bucket dsn", func(t *testing.T) {
		opts, bucket, err := s3_dsn.ParseS3Bucket("s3://access_key:secret_key@maggie.minio.geektr.co:9000/foobar/path/to/key?region=")
		if bucket != "" || opts != nil {
			t.Error("wrong bucket dsn parsed")
			return
		}

		if fmt.Sprint(err) != fmt.Sprintf("invalid s3 bucket dsn: unexpected key: %s", "path/to/key") {
			t.Error(err)
			return
		}

		_, bucket2, _ := s3_dsn.ParseS3Bucket("s3://access_key:secret_key@maggie.minio.geektr.co:9000/foobar?region=")

		if bucket2 != "foobar" {
			t.Error("bucket not parsed")
			return
		}
	})

	t.Run("s3 object dsn", func(t *testing.T) {
		_, _, key, err := s3_dsn.ParseS3Object("s3://access_key:secret_key@maggie.minio.geektr.co:9000/foobar/path/to/key?region=")
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

package s3_dsn_test

import (
	"testing"

	"github.com/geektheripper/vast-dsn/dsn/s3_dsn"
)

func TestS3DSNParser(t *testing.T) {
	t.Run("amazon s3", func(t *testing.T) {
		config, err := s3_dsn.ParseDSN("s3://access_key:secret_key@-?region=us-east-2")
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
		config, err := s3_dsn.ParseDSN("s3://access_key:secret_key@maggie.minio.geektr.co:9003?region=")
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
}

package s3_dsn

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/geektheripper/vast-dsn/utils"
)

type S3Options struct {
	Protocol        string
	AccessKeyID     string
	SecretAccessKey string
	Host            string
	Region          string
	NoVerifySSL     bool
	UsePathStyle    bool
}

func Parse(dsn string) (opts *S3Options, bucket string, key string, err error) {
	opts = &S3Options{
		Protocol: "https",
		Region:   "us-east-1",

		NoVerifySSL:  false,
		UsePathStyle: false,
	}

	url, err := url.Parse(dsn)
	if err != nil {
		return
	}

	if url.Scheme != "s3" {
		err = errors.New("invalid scheme")
		return
	}

	opts.AccessKeyID = url.User.Username()
	opts.SecretAccessKey, _ = url.User.Password()
	if (opts.AccessKeyID != "") != (opts.SecretAccessKey != "") {
		err = errors.New("invalid credentials: access_key_id and secret_access_key must be provided together")
		return
	}

	if url.Query().Has("protocol") {
		opts.Protocol = url.Query().Get("protocol")
		if opts.Protocol != "http" && opts.Protocol != "https" {
			err = errors.New("invalid protocol")
			return
		}
	}

	if url.Query().Has("region") {
		opts.Region = url.Query().Get("region")
	}

	if url.Query().Has("no-verify-ssl") {
		opts.NoVerifySSL = url.Query().Get("no-verify-ssl") == "true"
	}

	if url.Query().Has("force-path-style") {
		opts.UsePathStyle = url.Query().Get("force-path-style") == "true"
	}

	if url.Query().Has("use-path-style") {
		opts.UsePathStyle = url.Query().Get("use-path-style") == "true"
	}

	path := url.Path
	if path != "" {
		parts := strings.Split(strings.TrimLeft(path, "/"), "/")
		bucket = parts[0]
		key = strings.Join(parts[1:], "/")
	}

	opts.Host = url.Host

	return
}

func MustParse(dsn string, logger ...utils.Logger) (opts *S3Options, bucket string, key string) {
	log := utils.EnsureLogger(logger...)

	opts, bucket, key, err := Parse(dsn)
	if err != nil {
		log.Fatalf("failed to parse s3 dsn: %v", err)
	}
	return opts, bucket, key
}

func ParseS3(dsn string) (opts *S3Options, err error) {
	opts, bucket, _, err := Parse(dsn)
	if bucket != "" {
		return nil, fmt.Errorf("invalid s3 dsn: unexpected bucket: %s", bucket)
	}

	return
}

func MustNewS3(dsn string, logger ...utils.Logger) *S3Options {
	log := utils.EnsureLogger(logger...)

	s3dsn, err := ParseS3(dsn)
	if err != nil {
		log.Fatalf("failed to parse s3 dsn: %v", err)
	}
	return s3dsn
}

func ParseS3Bucket(dsn string) (opts *S3Options, bucket string, err error) {
	opts, bucket, key, err := Parse(dsn)

	if bucket == "" {
		return nil, "", fmt.Errorf("invalid s3 bucket dsn: missing bucket")
	}

	if key != "" {
		return nil, "", fmt.Errorf("invalid s3 bucket dsn: unexpected key: %s", key)
	}

	return
}

func MustParseS3Bucket(dsn string, logger ...utils.Logger) (*S3Options, string) {
	log := utils.EnsureLogger(logger...)

	s3dsn, bucket, err := ParseS3Bucket(dsn)
	if err != nil {
		log.Fatalf("failed to parse s3 bucket dsn: %v", err)
	}
	return s3dsn, bucket
}

func ParseS3Path(dsn string) (opts *S3Options, bucket string, path string, err error) {
	opts, bucket, key, err := Parse(dsn)
	if bucket == "" {
		return nil, "", "", fmt.Errorf("invalid s3 bucket path dsn: missing bucket")
	}

	if !strings.HasSuffix(key, "/") {
		return nil, "", "", fmt.Errorf("invalid s3 bucket path dsn: path must end with a slash")
	}

	path = strings.TrimSuffix(key, "/")

	return
}

func ParseNewS3Path(dsn string, logger ...utils.Logger) (*S3Options, string, string) {
	log := utils.EnsureLogger(logger...)

	s3dsn, bucket, path, err := ParseS3Path(dsn)
	if err != nil {
		log.Fatalf("failed to parse s3 path dsn: %v", err)
	}
	return s3dsn, bucket, path
}

func ParseS3Object(dsn string) (opts *S3Options, bucket string, key string, err error) {
	opts, bucket, key, err = Parse(dsn)
	if bucket == "" {
		return nil, "", "", fmt.Errorf("invalid s3 object dsn: missing bucket")
	}
	if key == "" {
		return nil, "", "", fmt.Errorf("invalid s3 object dsn: missing key")
	}
	return
}

func MustParseS3Object(dsn string, logger ...utils.Logger) (*S3Options, string, string) {
	log := utils.EnsureLogger(logger...)

	s3dsn, bucket, key, err := ParseS3Object(dsn)
	if err != nil {
		log.Fatalf("failed to parse s3 object dsn: %v", err)
	}
	return s3dsn, bucket, key
}

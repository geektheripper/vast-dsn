package s3_dsn

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

func ParseDSN(dsn string) (config *aws.Config, bucket string, key string, err error) {
	config = &aws.Config{}

	url, err := url.Parse(dsn)
	if err != nil {
		return
	}

	if url.Scheme != "s3" {
		err = errors.New("invalid scheme")
		return
	}

	user := url.User.Username()
	pass, passOk := url.User.Password()

	protocol := "https"
	if url.Query().Has("protocol") {
		protocol = url.Query().Get("protocol")
		if protocol != "http" && protocol != "https" {
			err = errors.New("invalid protocol")
			return
		}
	}

	region := "us-east-1"
	if url.Query().Has("region") {
		region = url.Query().Get("region")
	}

	disableSsl := url.Query().Get("disable-ssl") == "true"
	s3ForcePathStyle := url.Query().Get("s3-force-path-style") == "true"

	if user != "" && passOk {
		config.Credentials = credentials.NewStaticCredentials(user, pass, "")
	}

	endpoint := url.Host
	if endpoint != "-" {
		config.Endpoint = aws.String(protocol + "://" + endpoint)
	}

	config.Region = aws.String(region)
	config.DisableSSL = aws.Bool(disableSsl)
	config.S3ForcePathStyle = aws.Bool(s3ForcePathStyle)

	path := url.Path
	if path != "" {
		parts := strings.Split(strings.TrimLeft(path, "/"), "/")
		bucket = parts[0]
		key = strings.Join(parts[1:], "/")
	}

	return
}

func MustParseDSN(dsn string) (config *aws.Config, bucket string, key string, err error) {
	config, bucket, key, err = ParseDSN(dsn)
	if err != nil {
		panic(err)
	}
	return
}

func ParseS3DSN(dsn string) (config *aws.Config, err error) {
	config, bucket, _, err := ParseDSN(dsn)
	if bucket != "" {
		return nil, fmt.Errorf("invalid s3 dsn: unexpected bucket: %s", bucket)
	}

	return
}

func MustParseS3DSN(dsn string) (config *aws.Config, err error) {
	config, err = ParseS3DSN(dsn)
	if err != nil {
		panic(err)
	}
	return
}

func ParseS3BucketDSN(dsn string) (config *aws.Config, bucket string, err error) {
	config, bucket, key, err := ParseDSN(dsn)

	if bucket == "" {
		return nil, "", fmt.Errorf("invalid s3 bucket dsn: missing bucket")
	}

	if key != "" {
		return nil, "", fmt.Errorf("invalid s3 bucket dsn: unexpected key: %s", key)
	}

	return
}

func MustParseS3BucketDSN(dsn string) (config *aws.Config, bucket string, err error) {
	config, bucket, err = ParseS3BucketDSN(dsn)
	if err != nil {
		panic(err)
	}
	return
}

func ParseS3ObjectDSN(dsn string) (config *aws.Config, bucket string, key string, err error) {
	config, bucket, key, err = ParseDSN(dsn)
	if bucket == "" {
		return nil, "", "", fmt.Errorf("invalid s3 object dsn: missing bucket")
	}
	if key == "" {
		return nil, "", "", fmt.Errorf("invalid s3 object dsn: missing key")
	}
	return
}

func MustParseS3ObjectDSN(dsn string) (config *aws.Config, bucket string, key string, err error) {
	config, bucket, key, err = ParseS3ObjectDSN(dsn)
	if err != nil {
		panic(err)
	}
	return
}

package s3_dsn

import (
	"errors"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

func ParseDSN(dsn string) (config *aws.Config, err error) {
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

	return
}

func MustParseDSN(dsn string) *aws.Config {
	config, err := ParseDSN(dsn)
	if err != nil {
		panic(err)
	}

	return config
}

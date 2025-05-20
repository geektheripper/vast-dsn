package s3_client

import (
	"crypto/tls"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/geektheripper/vast-dsn/s3_dsn"
)

func NewS3Client(dsn *s3_dsn.S3Options) (*s3.Client, error) {
	client := s3.NewFromConfig(aws.Config{Region: dsn.Region}, func(o *s3.Options) {
		if dsn.Host != "" {
			o.BaseEndpoint = aws.String(dsn.Host)
		}
		if dsn.Protocol != "" {
			o.BaseEndpoint = aws.String(dsn.Protocol + "://" + dsn.Host)
		}
		if dsn.NoVerifySSL {
			o.HTTPClient = &http.Client{
				Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
			}
		}
		if dsn.AccessKeyID != "" {
			o.Credentials = credentials.NewStaticCredentialsProvider(dsn.AccessKeyID, dsn.SecretAccessKey, "")
		}
		o.EndpointResolverV2 = s3.NewDefaultEndpointResolverV2()
		o.UsePathStyle = dsn.UsePathStyle
	})

	return client, nil
}

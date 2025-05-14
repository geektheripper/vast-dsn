package s3_client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"github.com/geektheripper/vast-dsn/s3_dsn"
)

type Resolver struct {
	Protocol string
}

func (r *Resolver) ResolveEndpoint(ctx context.Context, params s3.EndpointParameters) (smithyendpoints.Endpoint, error) {
	if *(params.Endpoint) == "-" {
		return s3.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	subdomain := *params.Endpoint
	if params.Bucket != nil {
		subdomain = fmt.Sprintf("%s.%s", *params.Bucket, *params.Endpoint)
	}

	u := url.URL{
		Scheme: r.Protocol,
		Host:   subdomain,
		Path:   "",
	}

	if *params.ForcePathStyle && params.Bucket != nil {
		u.Host = *params.Endpoint
		u.Path += "/" + *params.Bucket
	}

	return smithyendpoints.Endpoint{URI: u}, nil
}

func NewS3Client(dsn *s3_dsn.S3Options) (*s3.Client, error) {
	config := aws.Config{BaseEndpoint: aws.String(dsn.Host)}

	if dsn.Region != "" {
		config.Region = dsn.Region
	}

	if dsn.NoVerifySSL {
		httpclient := awshttp.NewBuildableClient()
		httpclient.WithTransportOptions(func(transport *http.Transport) {
			if transport.TLSClientConfig == nil {
				transport.TLSClientConfig = &tls.Config{}
			}
			transport.TLSClientConfig.InsecureSkipVerify = true
		})
		config.HTTPClient = httpclient
	}

	if dsn.AccessKeyID != "" {
		config.Credentials = aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     dsn.AccessKeyID,
				SecretAccessKey: dsn.SecretAccessKey,
			}, nil
		})
	}

	client := s3.NewFromConfig(config, func(options *s3.Options) {
		options.EndpointResolverV2 = &Resolver{Protocol: dsn.Protocol}
		options.UsePathStyle = dsn.UsePathStyle
	})

	return client, nil
}

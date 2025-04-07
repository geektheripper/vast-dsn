package s3_dsn

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"github.com/geektheripper/vast-dsn/utils"
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

func Load(dsn string) (client *s3.Client, bucket string, key string, err error) {
	url, err := url.Parse(dsn)
	if err != nil {
		return
	}

	if url.Scheme != "s3" {
		err = errors.New("invalid scheme")
		return
	}

	access_key_id := url.User.Username()
	secret_access_key, _ := url.User.Password()
	if (access_key_id != "") != (secret_access_key != "") {
		err = errors.New("invalid credentials: access_key_id and secret_access_key must be provided together")
		return
	}

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

	no_verify_ssl := url.Query().Get("no-verify-ssl") == "true"
	use_path_style := false
	_force_path_style := url.Query().Get("force-path-style")
	_use_path_style := url.Query().Get("use-path-style")
	if _use_path_style != "" {
		use_path_style = _use_path_style == "true"
	} else {
		use_path_style = _force_path_style == "true"
	}

	path := url.Path
	if path != "" {
		parts := strings.Split(strings.TrimLeft(path, "/"), "/")
		bucket = parts[0]
		key = strings.Join(parts[1:], "/")
	}

	config := aws.Config{BaseEndpoint: aws.String(url.Host)}

	if region != "" {
		config.Region = region
	}

	if no_verify_ssl {
		httpclient := awshttp.NewBuildableClient()
		httpclient.WithTransportOptions(func(transport *http.Transport) {
			if transport.TLSClientConfig == nil {
				transport.TLSClientConfig = &tls.Config{}
			}
			transport.TLSClientConfig.InsecureSkipVerify = true
		})
		config.HTTPClient = httpclient
	}

	if access_key_id != "" {
		config.Credentials = aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     access_key_id,
				SecretAccessKey: secret_access_key,
			}, nil
		})
	}

	client = s3.NewFromConfig(config, func(options *s3.Options) {
		options.EndpointResolverV2 = &Resolver{Protocol: protocol}
		options.UsePathStyle = use_path_style
	})

	return
}

func MustLoad(dsn string, logger ...utils.Logger) (client *s3.Client, bucket string, key string) {
	log := utils.EnsureLogger(logger...)

	client, bucket, key, err := Load(dsn)
	if err != nil {
		log.Fatalf("failed to parse s3 dsn: %v", err)
	}
	return client, bucket, key
}

func NewS3(dsn string) (client *s3.Client, err error) {
	client, bucket, _, err := Load(dsn)
	if bucket != "" {
		return nil, fmt.Errorf("invalid s3 dsn: unexpected bucket: %s", bucket)
	}

	return
}

func MustNewS3(dsn string, logger ...utils.Logger) *s3.Client {
	log := utils.EnsureLogger(logger...)

	client, err := NewS3(dsn)
	if err != nil {
		log.Fatalf("failed to parse s3 dsn: %v", err)
	}
	return client
}

func NewS3Bucket(dsn string) (client *s3.Client, bucket string, err error) {
	client, bucket, key, err := Load(dsn)

	if bucket == "" {
		return nil, "", fmt.Errorf("invalid s3 bucket dsn: missing bucket")
	}

	if key != "" {
		return nil, "", fmt.Errorf("invalid s3 bucket dsn: unexpected key: %s", key)
	}

	return
}

func MustNewS3Bucket(dsn string, logger ...utils.Logger) (*s3.Client, string) {
	log := utils.EnsureLogger(logger...)

	client, bucket, err := NewS3Bucket(dsn)
	if err != nil {
		log.Fatalf("failed to parse s3 bucket dsn: %v", err)
	}
	return client, bucket
}

func NewS3Object(dsn string) (client *s3.Client, bucket string, key string, err error) {
	client, bucket, key, err = Load(dsn)
	if bucket == "" {
		return nil, "", "", fmt.Errorf("invalid s3 object dsn: missing bucket")
	}
	if key == "" {
		return nil, "", "", fmt.Errorf("invalid s3 object dsn: missing key")
	}
	return
}

func MustNewS3Object(dsn string, logger ...utils.Logger) (client *s3.Client, bucket string, key string) {
	log := utils.EnsureLogger(logger...)

	client, bucket, key, err := NewS3Object(dsn)
	if err != nil {
		log.Fatalf("failed to parse s3 object dsn: %v", err)
	}
	return client, bucket, key
}

# Vast DSN

DSN (data source name) for anything.

## S3

`s3://[<credentials>@]<endpoint>[:(port)]?region=<region>&querys...`

_credentials_:

`<access-key-id>:<secret-access-key>`

_endpoint_:

`-` for aws s3, use default endpoint by sdk
`some-domain` from custom endpoint or s3 alternatives

_querys_:

- protocol: `http|https`, default `https`
- region: `string`, default `aws-east-1`
- disable-ssl: `true|false`, default `false`
- s3-force-path-style: `true|false`, default `false`

S3 Bucket:

`s3://[<credentials>@]<endpoint>[:(port)]/<bucket>[?region=<region>&querys...]`

S3 Object:

`s3://[<credentials>@]<endpoint>[:(port)]/<bucket>/<key>[?region=<region>&querys...]`

## S3 DSN Usage

```go
import "github.com/geektheripper/vast-dsn/dsn/s3_dsn"

config, bucket, key, err = s3_dsn.ParseDSN("s3://minio.vastdns.example.com:9003/foobar/path/to/key?region=")

// error when got unexpected bucket
config, err := s3_dsn.ParseS3DSN("s3://access_key:secret_key@-?region=us-east-2")

// error when got unexpected key
config, bucket, err := s3_dsn.ParseS3BucketDSN("s3://access_key:secret_key@minio.vastdns.example.com:9003/foobar?region=")

config, bucket, key, err = s3_dsn.ParseS3ObjectDSN("s3://minio.vastdns.example.com:9003/foobar/path/to/key?region=")

// panic when got error
config := s3_dsn.MustParseS3DSN("...")
config := s3_dsn.ParseS3BucketDSN("...")
config := s3_dsn.ParseS3ObjectDSN("...")
```

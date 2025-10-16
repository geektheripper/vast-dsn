# S3 DSN

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
- no-verify-ssl: `true|false`, default `false`
- use-path-style: `true|false`, default `false`

S3 Bucket:

`s3://[<credentials>@]<endpoint>[:(port)]/<bucket>[?region=<region>&querys...]`

S3 Object:

`s3://[<credentials>@]<endpoint>[:(port)]/<bucket>/<key>[?region=<region>&querys...]`

## S3 DSN Usage

```go
import "github.com/geektheripper/vast-dsn/dsn/s3_dsn"

opts, bucket, key, err = s3_dsn.Parse("s3://minio.vastdns.example.com:9003/foobar/path/to/key?region=")

// error when got unexpected bucket
opts, err := s3_dsn.ParseS3("s3://access_key:secret_key@-?region=us-east-2")

// error when got unexpected key
opts, bucket, err := s3_dsn.ParseS3Bucket("s3://access_key:secret_key@minio.vastdns.example.com:9003/foobar?region=")

// error when ends with non slash character
opts, bucket, key, err = s3_dsn.ParseS3Object("s3://minio.vastdns.example.com:9003/foobar/path/to/dir/?region=")

opts, bucket, key, err = s3_dsn.ParseS3Object("s3://minio.vastdns.example.com:9003/foobar/path/to/key?region=")

// panic when got error
opts := s3_dsn.MustParseS3("...")
opts, bucket := s3_dsn.MustParseS3Bucket("...")
opts, bucket, path := s3_dsn.MustParseS3Path("...")
opts, bucket, key := s3_dsn.MustParseS3Object("...")
```

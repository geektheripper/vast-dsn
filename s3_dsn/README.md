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

client, bucket, key, err = s3_dsn.Load("s3://minio.vastdns.example.com:9003/foobar/path/to/key?region=")

// error when got unexpected bucket
client, err := s3_dsn.NewS3("s3://access_key:secret_key@-?region=us-east-2")

// error when got unexpected key
client, bucket, err := s3_dsn.NewS3Bucket("s3://access_key:secret_key@minio.vastdns.example.com:9003/foobar?region=")

// error when ends with non slash character
client, bucket, key, err = s3_dsn.NewS3Object("s3://minio.vastdns.example.com:9003/foobar/path/to/dir/?region=")

client, bucket, key, err = s3_dsn.NewS3Object("s3://minio.vastdns.example.com:9003/foobar/path/to/key?region=")

// panic when got error
client := s3_dsn.MustNewS3("...")
client, bucket := s3_dsn.MustNewS3Bucket("...")
client, bucket, path := s3_dsn.MustNewS3Path("...")
client, bucket, key := s3_dsn.MustNewS3Object("...")
```

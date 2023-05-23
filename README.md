# Vast DSN

DSN (data source name) for anything.

## S3

`s3://[<credentials>@]<endpoint>[:(port)]?region=<region>&querys...`

*credentials*:

`<access-key-id>:<secret-access-key>`

*endpoint*:

`-` for aws s3, use default endpoint by sdk
`some-domain` from custom endpoint or s3 alternatives

*querys*:

- protocol: `http|https`, default `https`
- region: `string`, default `aws-east-1`
- disable-ssl: `true|false`, default `false`
- s3-force-path-style: `true|false`, default `false`

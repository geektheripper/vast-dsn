# Redis DSN

`<protocol>://[<credentials>@]<host>[:(port)][/<db>[/<key-prefix>]]?db=<db>[&querys...]`

_protocol_:

`redis|rediss|unix`

_credentials_:

`[<username>:]<password>`

_querys_:

part of queries, access [options.go](https://github.com/redis/go-redis/blob/8fadbef84a3f4e7573f8b38e5023fd469470a8a4/options.go#L452-L501) for more

- client_name: `string`
- max_retries: `int`
- pool_size: `int`
- read_timeout: `duration`
- write_timeout: `duration`

## Redis DSN Usage

```go
import "github.com/geektheripper/vast-dsn/dsn/redis_dsn"

config, err := redis_dsn.Parse("redis://localhost:6379/0/key:prefix:for:biz?client_name=biz&max_retries=10")
```

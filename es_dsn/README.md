# Elasticsearch DSN


## DSN Pattern

`<http|https>://[<credentials>@]<endpoint>[:<port>][?index=<index>|index_prefix=<index_prefix>]`

_credentials_:

`<username>:<password>`

_endpoint_:

`<domain or ip>`

_index_:

`<index>` specify index for application

_index_prefix_:

`<index_prefix>` specify index prefix for application, application should use `<index_prefix>_<index>` as index name


## Package Usage

```go
import "github.com/geektheripper/vast-dsn/es_dsn"

esc, err := es_dsn.Parse("http://username:password@es.vastdns.example.com:9200")

esc.Config // *elasticsearch.Config
esc.Index // string
esc.IndexPrefix // string

// panic if error
esc := es_dsn.MustParse("http://username:password@es.vastdns.example.com:9200")
```

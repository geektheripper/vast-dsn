# OIDC DSN

## DSN Pattern

`oidc://<issuer_url>?client_id=<client_id>&client_secret=<client_secret>`

## Package Usage

```go
import "github.com/geektheripper/vast-dsn/oidc_dsn"

oidc, err := oidc_dsn.Parse("oidc://keycloak.vastdns.example.com:9200/realms/myrealm?client_id=myid&client_secret=mysecret")

oidc.Issuer // string
oidc.ClientID // string
oidc.ClientSecret // string

// panic if error
oidc := oidc_dsn.MustParse("oidc://keycloak.vastdns.example.com:9200/realms/myrealm?client_id=myid&client_secret=mysecret")
```

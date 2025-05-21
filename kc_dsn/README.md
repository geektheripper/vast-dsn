# Keycloak DSN

## DSN Pattern

`keycloak://<username>:<password>@<url>/realms/<realm>?client_id=<client_id>&client_secret=<client_secret>`

## Package Usage

```go
import "github.com/geektheripper/vast-dsn/kc_dsn"

// parse dsn string
kc, err := kc_dsn.Parse("keycloak://username:password@keycloak.vastdns.example.com:9200/realms/realm?client_id=client_id&client_secret=client_secret")

// get fields
kc.URL // string
kc.Realm // string
kc.Username // string
kc.Password // string
kc.ClientID // string
kc.ClientSecret // string

// panic if error
kc := kc_dsn.MustParse("keycloak://username:password@keycloak.vastdns.example.com:9200/realms/realm?client_id=client_id&client_secret=client_secret")

// convert to dsn string
dsn := kc.String()

// extract oidc config
oidc, err := kc.OIDCConfig()
```

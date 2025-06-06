package kc_dsn_test

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/geektheripper/vast-dsn/kc_dsn"
)

var ExpectURL = "https://keycloak.vastdns.example.com:9200"
var ExpectRealm = "myrealm"
var ExpectUsername = "superuser"
var ExpectPassword = "super:@/#*pass"
var ExpectClientID = "myid"
var ExpectClientSecret = "mysecret"
var ExpectIssuer = "https://keycloak.vastdns.example.com:9200/realms/myrealm"

var ExampleDSN = fmt.Sprintf("keycloak://%s:%s@keycloak.vastdns.example.com:9200/realms/%s?client_id=%s&client_secret=%s",
	ExpectUsername,
	url.QueryEscape(ExpectPassword),
	ExpectRealm,
	ExpectClientID,
	ExpectClientSecret,
)

func TestParseDSN(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		kcCfg, err := kc_dsn.Parse(ExampleDSN)
		if err != nil {
			t.Error(err)
			return
		}

		if kcCfg.ClientID != ExpectClientID ||
			kcCfg.ClientSecret != ExpectClientSecret ||
			kcCfg.Username != ExpectUsername ||
			kcCfg.Password != ExpectPassword ||
			kcCfg.URL != ExpectURL ||
			kcCfg.Realm != ExpectRealm {
			t.Error("dsn not parsed")
			return
		}
	})

	t.Run("invalid dsn", func(t *testing.T) {
		_, err := kc_dsn.Parse(strings.Replace(ExampleDSN, "keycloak", "http", 1))
		if err == nil {
			t.Error("error should not be nil")
			return
		}
	})

	t.Run("to string", func(t *testing.T) {
		if kc_dsn.MustParse(ExampleDSN).String() != ExampleDSN {
			t.Error("string not parsed")
			return
		}
	})

	t.Run("extract oidc config", func(t *testing.T) {
		kcCfg, err := kc_dsn.Parse(ExampleDSN)
		if err != nil {
			t.Error(err)
			return
		}

		oidcCfg, err := kcCfg.OIDCConfig()
		if err != nil {
			t.Error(err)
			return
		}

		if oidcCfg.Issuer != ExpectIssuer ||
			oidcCfg.ClientID != ExpectClientID ||
			oidcCfg.ClientSecret != ExpectClientSecret {
			t.Error("oidc config not parsed")
			return
		}
	})
}

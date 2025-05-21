package oidc_dsn_test

import (
	"strings"
	"testing"

	"github.com/geektheripper/vast-dsn/oidc_dsn"
)

var ExampleDSN = "oidc://keycloak.vastdns.example.com:9200/realms/myrealm?client_id=myid&client_secret=mysecret"
var ExpectIssuer = "https://keycloak.vastdns.example.com:9200/realms/myrealm"
var ExpectClientID = "myid"
var ExpectClientSecret = "mysecret"

func TestParseDSN(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		oidcCfg, err := oidc_dsn.Parse(ExampleDSN)
		if err != nil {
			t.Error(err)
			return
		}

		if oidcCfg.ClientID != ExpectClientID || oidcCfg.ClientSecret != ExpectClientSecret {
			t.Error("client credentials not parsed")
			return
		}

		if oidcCfg.Issuer != ExpectIssuer {
			t.Error("issuer url not parsed")
			return
		}
	})

	t.Run("invalid dsn", func(t *testing.T) {
		_, err := oidc_dsn.Parse(strings.Replace(ExampleDSN, "oidc", "http", 1))
		if err == nil {
			t.Error("error should not be nil")
			return
		}
	})

	t.Run("to string", func(t *testing.T) {
		if oidc_dsn.MustParse(ExampleDSN).String() != ExampleDSN {
			t.Error("string not parsed")
			return
		}
	})
}

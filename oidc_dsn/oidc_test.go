package oidc_dsn_test

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/geektheripper/vast-dsn/oidc_dsn"
)

var ExpectIssuer = "https://keycloak.vastdns.example.com:9200/realms/myrealm"
var ExpectClientID = "myid"
var ExpectClientSecret = "mysecret"
var ExpectScopes = []string{"openid", "profile", "email"}
var ExpectHttpEndpoint = "http://1.2.3.4:9876"

var ExampleDSN = fmt.Sprintf("oidc://keycloak.vastdns.example.com:9200/realms/myrealm?client_id=%s&client_secret=%s&scopes=%s&http_endpoint=%s",
	ExpectClientID,
	ExpectClientSecret,
	strings.Join(ExpectScopes, ","),
	url.QueryEscape(ExpectHttpEndpoint),
)

func TestParseDSN(t *testing.T) {
	ExampleDSN = oidc_dsn.MustParse(ExampleDSN).String()
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

		if !reflect.DeepEqual(oidcCfg.Scopes, ExpectScopes) {
			t.Error("scopes not parsed")
			return
		}

		if oidcCfg.HttpEndpoint.String() != ExpectHttpEndpoint {
			t.Error("http_endpoint not parsed")
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
		if oidc_dsn.MustParse(ExampleDSN).String() != strings.ReplaceAll(ExampleDSN, ",", "%2C") {
			t.Error("string not parsed")
			return
		}
	})

	t.Run("use http", func(t *testing.T) {
		oidcCfg, err := oidc_dsn.Parse(ExampleDSN + "&use-http=true")
		if err != nil {
			t.Error(err)
			return
		}

		if oidcCfg.Issuer != strings.Replace(ExpectIssuer, "https", "http", 1) {
			t.Error("issuer url not parsed")
			return
		}
	})
}

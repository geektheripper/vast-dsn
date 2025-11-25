package oidc_dsn

import (
	"errors"
	"log"
	"net/url"
	"strings"
)

type OIDCConfig struct {
	Issuer       string
	ClientID     string
	ClientSecret string
	Scopes       []string
	HttpEndpoint *url.URL
}

func (o *OIDCConfig) String() string {
	result, err := url.Parse(o.Issuer)
	if err != nil {
		return ""
	}

	result.Scheme = "oidc"

	values := url.Values{}
	values.Set("client_id", o.ClientID)
	values.Set("client_secret", o.ClientSecret)
	values.Set("scopes", strings.Join(o.Scopes, ","))
	values.Set("http_endpoint", o.HttpEndpoint.String())

	if result.Scheme == "http" {
		values.Set("use_http", "true")
	}
	result.RawQuery = values.Encode()

	return result.String()
}

func Parse(dsnStr string) (*OIDCConfig, error) {
	dsnUrl, err := url.Parse(dsnStr)
	if err != nil {
		return nil, err
	}

	if dsnUrl.Scheme != "oidc" {
		return nil, errors.New("invalid scheme, should be oidc")
	}

	schema := "https"
	if dsnUrl.Query().Get("use_http") == "true" {
		schema = "http"
	}

	issuerUrl := url.URL{Scheme: schema, Host: dsnUrl.Host, Path: dsnUrl.Path}
	oidcdsn := &OIDCConfig{
		Issuer:       issuerUrl.String(),
		ClientID:     dsnUrl.Query().Get("client_id"),
		ClientSecret: dsnUrl.Query().Get("client_secret"),
		Scopes:       strings.Split(dsnUrl.Query().Get("scopes"), ","),
	}

	http_endpoint := dsnUrl.Query().Get("http_endpoint")
	if http_endpoint != "" {
		oidcdsn.HttpEndpoint, err = url.Parse(http_endpoint)
		if err != nil {
			return nil, errors.New("invalid http_endpoint, failed to parse")
		}
	}

	if oidcdsn.ClientID == "" || oidcdsn.ClientSecret == "" {
		return nil, errors.New("client_id and client_secret are required")
	}

	return oidcdsn, nil
}

func MustParse(dsn string) *OIDCConfig {
	oidcdsn, err := Parse(dsn)
	if err != nil {
		log.Fatalf("failed to parse oidc dsn: %v", err)
	}
	return oidcdsn
}

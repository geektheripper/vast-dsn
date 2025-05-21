package kc_dsn

import (
	"errors"
	"log"
	"net/url"
	"regexp"

	"github.com/geektheripper/vast-dsn/oidc_dsn"
)

type KeycloakConfig struct {
	URL          string
	Realm        string
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
}

func (o *KeycloakConfig) ValidateURL() (*url.URL, error) {
	kcUrl, err := url.Parse(o.URL)
	if err != nil {
		return nil, err
	}

	if kcUrl.Path != "" {
		return nil, errors.New("invalid url, should not have path")
	}

	if kcUrl.Host == "" {
		return nil, errors.New("invalid url, should have host")
	}

	if kcUrl.Scheme != "https" {
		return nil, errors.New("invalid url, should have https scheme")
	}

	return kcUrl, nil
}

func (o *KeycloakConfig) OIDCConfig() (*oidc_dsn.OIDCConfig, error) {
	kcUrl, err := o.ValidateURL()
	if err != nil {
		return nil, err
	}

	oidcUrl := url.URL{
		Scheme: "https",
		Host:   kcUrl.Host,
		Path:   "/realms/" + o.Realm,
	}

	return &oidc_dsn.OIDCConfig{
		Issuer:       oidcUrl.String(),
		ClientID:     o.ClientID,
		ClientSecret: o.ClientSecret,
	}, nil
}

func (o *KeycloakConfig) ToDSN() (string, error) {
	kcUrl, err := o.ValidateURL()
	if err != nil {
		return "", err
	}

	values := url.Values{}
	values.Set("client_id", o.ClientID)
	values.Set("client_secret", o.ClientSecret)

	kcDsn := url.URL{
		Scheme:   "keycloak",
		Host:     kcUrl.Host,
		Path:     "/realms/" + o.Realm,
		RawQuery: values.Encode(),
	}

	kcDsn.User = url.UserPassword(o.Username, o.Password)
	return kcDsn.String(), nil
}

func (o *KeycloakConfig) String() string {
	dsn, err := o.ToDSN()
	if err != nil {
		log.Fatal("failed to convert keycloak config to dsn")
	}
	return dsn
}

var realmRegex = regexp.MustCompile(`^/realms/([^/]+)$`)

func Parse(dsnStr string) (*KeycloakConfig, error) {
	dsnUrl, err := url.Parse(dsnStr)
	if err != nil {
		return nil, err
	}

	if dsnUrl.Scheme != "keycloak" {
		return nil, errors.New("invalid scheme, should be keycloak://")
	}

	if dsnUrl.Host == "" {
		return nil, errors.New("invalid url, should have host part")
	}

	kcUrl := url.URL{Scheme: "https", Host: dsnUrl.Host}
	kcConfig := &KeycloakConfig{URL: kcUrl.String()}

	matches := realmRegex.FindStringSubmatch(dsnUrl.Path)
	if matches == nil {
		return nil, errors.New("invalid path, should be /realms/<realm>")
	}
	kcConfig.Realm = matches[1]

	kcConfig.Username = dsnUrl.User.Username()
	kcConfig.Password, _ = dsnUrl.User.Password()
	if kcConfig.Username == "" || kcConfig.Password == "" {
		return nil, errors.New("invalid credentials, should have username:password")
	}

	kcConfig.ClientID = dsnUrl.Query().Get("client_id")
	kcConfig.ClientSecret = dsnUrl.Query().Get("client_secret")
	if kcConfig.ClientID == "" || kcConfig.ClientSecret == "" {
		return nil, errors.New("invalid credentials, should have client_id and client_secret in query")
	}

	return kcConfig, nil
}

func MustParse(dsn string) *KeycloakConfig {
	kcConfig, err := Parse(dsn)
	if err != nil {
		log.Fatalf("failed to parse keycloak dsn: %v", err)
	}
	return kcConfig
}

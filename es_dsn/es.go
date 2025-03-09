package es_dsn

import (
	"errors"
	"net/url"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/geektheripper/vast-dsn/utils"
)

type EsDSN struct {
	Config      *elasticsearch.Config
	Index       string
	IndexPrefix string
}

func Parse(dsnStr string) (*EsDSN, error) {
	dsnUrl, err := url.Parse(dsnStr)
	if err != nil {
		return nil, err
	}

	if dsnUrl.Scheme != "http" && dsnUrl.Scheme != "https" {
		return nil, errors.New("invalid scheme, should be http or https")
	}

	esdsn := &EsDSN{Config: &elasticsearch.Config{}}

	esAddr := url.URL{Scheme: dsnUrl.Scheme, Host: dsnUrl.Host}
	esdsn.Config.Addresses = []string{esAddr.String()}
	esdsn.Config.Username = dsnUrl.User.Username()
	esdsn.Config.Password, _ = dsnUrl.User.Password()

	esdsn.Index = dsnUrl.Query().Get("index")
	esdsn.IndexPrefix = dsnUrl.Query().Get("index_prefix")

	if esdsn.Index != "" && esdsn.IndexPrefix != "" {
		return nil, errors.New("index and index_prefix cannot be used together")
	}

	return esdsn, nil
}

func MustParse(dsn string, logger ...utils.Logger) *EsDSN {
	log := utils.EnsureLogger(logger...)

	esdsn, err := Parse(dsn)
	if err != nil {
		log.Fatalf("failed to parse elasticsearch dsn: %v", err)
		os.Exit(1)
	}
	return esdsn
}

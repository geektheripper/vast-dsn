package es_dsn

import (
	"errors"
	"net/url"

	"github.com/elastic/go-elasticsearch/v8"
)

type EsDSN struct {
	Config      *elasticsearch.Config
	Index       string
	IndexPrefix string
}

func Parse(dsn string) (*EsDSN, error) {
	url, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	if url.Scheme != "http" && url.Scheme != "https" {
		return nil, errors.New("invalid scheme, should be http or https")
	}

	esdsn := &EsDSN{Config: &elasticsearch.Config{}}

	esdsn.Config.Username = url.User.Username()
	esdsn.Config.Password, _ = url.User.Password()

	esdsn.Index = url.Query().Get("index")
	esdsn.IndexPrefix = url.Query().Get("index_prefix")

	if esdsn.Index != "" && esdsn.IndexPrefix != "" {
		return nil, errors.New("index and index_prefix cannot be used together")
	}

	return esdsn, nil
}

func MustParse(dsn string) *EsDSN {
	esdsn, err := Parse(dsn)
	if err != nil {
		panic(err)
	}
	return esdsn
}

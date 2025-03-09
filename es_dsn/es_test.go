package es_dsn_test

import (
	"testing"

	"github.com/geektheripper/vast-dsn/es_dsn"
)

func TestParseDSN(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		esdsn, err := es_dsn.Parse("http://username:password@es.vastdns.example.com:9200")
		if err != nil {
			t.Error(err)
			return
		}

		if esdsn.Config.Username != "username" || esdsn.Config.Password != "password" {
			t.Error("credentials not parsed")
			return
		}

		if esdsn.Index != "" || esdsn.IndexPrefix != "" {
			t.Error("index and index prefix should be empty")
			return
		}
	})

	t.Run("with index", func(t *testing.T) {
		esdsn, err := es_dsn.Parse("http://username:password@es.vastdns.example.com:9200?index=index")
		if err != nil {
			t.Error(err)
			return
		}

		if esdsn.Index != "index" {
			t.Error("index not parsed")
			return
		}

		if esdsn.IndexPrefix != "" {
			t.Error("index prefix should be empty")
			return
		}
	})

	t.Run("with index and index prefix", func(t *testing.T) {
		_, err := es_dsn.Parse("http://username:password@es.vastdns.example.com:9200?index=index&index_prefix=index_prefix")
		if err == nil {
			t.Error("error should not be nil")
			return
		}
	})
}

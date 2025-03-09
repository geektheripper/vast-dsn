package redis_dsn_test

import (
	"testing"

	"github.com/geektheripper/vast-dsn/redis_dsn"
)

func TestRedisDSNParser(t *testing.T) {
	t.Run("parse", func(t *testing.T) {
		dsn := "redis://user:pass@host:6379/1231/prefix:biz?client_name=biz&max_retries=10"
		redisDSN, err := redis_dsn.Parse(dsn)
		if err != nil {
			t.Fatalf("failed to parse DSN: %v", err)
		}

		if redisDSN.Options.Addr != "host:6379" {
			t.Fatalf("expected addr: %s, got: %s", "host:6379", redisDSN.Options.Addr)
		}

		if redisDSN.KeyPrefix != "prefix:biz" {
			t.Fatalf("expected key prefix: %s, got: %s", "prefix:biz", redisDSN.KeyPrefix)
		}

		if redisDSN.Options.ClientName != "biz" {
			t.Fatalf("expected client name: %s, got: %s", "biz", redisDSN.Options.ClientName)
		}

		if redisDSN.Options.MaxRetries != 10 {
			t.Fatalf("expected max retries: %d, got: %d", 10, redisDSN.Options.MaxRetries)
		}
	})

	t.Run("parse failed", func(t *testing.T) {
		dsn := "s3://user:pass@host:6379/1231/prefix:biz?client_name=biz&max_retries=10"
		redisDSN, err := redis_dsn.Parse(dsn)
		if err == nil {
			t.Fatalf("expected error, got: %v", redisDSN)
		}
	})
}

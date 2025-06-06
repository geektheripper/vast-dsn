package redis_dsn_test

import (
	"testing"

	"github.com/geektheripper/vast-dsn/redis_dsn"
)

func TestRedisDSNParser(t *testing.T) {
	t.Run("parse", func(t *testing.T) {
		dsn := "redis://user:pass@host:6379/1231/prefix:biz?client_name=biz&max_retries=10"
		opts, prefix, err := redis_dsn.Parse(dsn)
		if err != nil {
			t.Fatalf("failed to parse DSN: %v", err)
		}

		if opts.Addr != "host:6379" {
			t.Fatalf("expected addr: %s, got: %s", "host:6379", opts.Addr)
		}

		if prefix != "prefix:biz" {
			t.Fatalf("expected key prefix: %s, got: %s", "prefix:biz", prefix)
		}

		if opts.ClientName != "biz" {
			t.Fatalf("expected client name: %s, got: %s", "biz", opts.ClientName)
		}

		if opts.MaxRetries != 10 {
			t.Fatalf("expected max retries: %d, got: %d", 10, opts.MaxRetries)
		}
	})

	t.Run("parse failed", func(t *testing.T) {
		dsn := "s3://user:pass@host:6379/1231/prefix:biz?client_name=biz&max_retries=10"
		opts, _, err := redis_dsn.Parse(dsn)
		if err == nil {
			t.Fatalf("expected error, got: %v", opts)
		}
	})
}

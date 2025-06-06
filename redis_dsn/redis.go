package redis_dsn

import (
	"errors"
	"log"
	"net/url"
	"regexp"

	"github.com/redis/go-redis/v9"
)

type RedisDSN struct {
	Options   *redis.Options
	KeyPrefix string
}

var pathRegex = regexp.MustCompile(`^/(\d+)(/(.+))?$`)

func Parse(dsn string) (opts *redis.Options, keyPrefix string, err error) {
	dsnUrl, err := url.Parse(dsn)
	if err != nil {
		return nil, "", err
	}

	standardDSN := *dsnUrl
	if standardDSN.Path != "" {
		matches := pathRegex.FindStringSubmatch(standardDSN.Path)
		if len(matches) == 0 {
			return nil, "", errors.New("invalid path")
		}
		standardDSN.Path = "/" + matches[1]
		if len(matches) > 3 {
			keyPrefix = matches[3]
		}
	}

	opts, err = redis.ParseURL(standardDSN.String())
	if err != nil {
		return nil, "", err
	}

	return opts, keyPrefix, nil
}

func MustParse(dsn string) (*redis.Options, string) {
	opts, prefix, err := Parse(dsn)
	if err != nil {
		log.Fatalf("failed to parse redis dsn: %v", err)
	}
	return opts, prefix
}

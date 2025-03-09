package redis_dsn

import (
	"errors"
	"net/url"
	"regexp"

	"github.com/geektheripper/vast-dsn/utils"
	"github.com/redis/go-redis/v9"
)

type RedisDSN struct {
	Options   *redis.Options
	KeyPrefix string
}

var pathRegex = regexp.MustCompile(`^/(\d+)(/(.+))?$`)

func Parse(dsn string) (*RedisDSN, error) {
	result := &RedisDSN{}

	dsnUrl, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	standardDSN := *dsnUrl
	if standardDSN.Path != "" {
		matches := pathRegex.FindStringSubmatch(standardDSN.Path)
		if len(matches) == 0 {
			return nil, errors.New("invalid path")
		}
		standardDSN.Path = "/" + matches[1]
		if len(matches) > 3 {
			result.KeyPrefix = matches[3]
		}
	}

	options, err := redis.ParseURL(standardDSN.String())
	if err != nil {
		return nil, err
	}

	result.Options = options

	return result, nil
}

func MustParse(dsn string, logger ...utils.Logger) *RedisDSN {
	log := utils.EnsureLogger(logger...)

	redisDSN, err := Parse(dsn)
	if err != nil {
		log.Fatalf("failed to parse redis dsn: %v", err)
	}
	return redisDSN
}

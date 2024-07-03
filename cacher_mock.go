package tools

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

// use this when init for ServiceContext, for local test
func MockRedis() (*redis.Client, error) {
	mr, err := miniredis.Run()
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	_, err = miniredis.Run()
	if err != nil {
		return nil, err
	}

	return client, nil
}

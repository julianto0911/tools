package tools

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// this is useful for testing, to predefined behavior of the response

type CacherV2 interface {
	Ping() error
	Get(name string) (string, error)
	Set(name string, value string) error
	SetWithDuration(name string, value string, d time.Duration) error
	Delete(name string) error
	GetKeysWithParam(name string) ([]string, error)
	PrintKeys()
}

func NewCacherV2(rdc *redis.Client, prefix string, expiracy int) CacherV2 {
	return &cacher{
		rdb:      rdc,
		expiracy: time.Duration(expiracy) * time.Second,
		prefix:   prefix,
	}
}

type cacher struct {
	rdb      *redis.Client
	expiracy time.Duration
	prefix   string
}

func (c *cacher) PrintKeys() {
	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = c.rdb.Scan(ctxB, cursor, "", 0).Result()
		if err != nil {
			panic(err)
		}

		for _, key := range keys {
			fmt.Println("key", key)
		}

		if cursor == 0 { // no more keys
			break
		}
	}
}

func (c *cacher) SetWithDuration(name string, value string, d time.Duration) error {
	return c.rdb.Set(ctxB, c.prefix+"_"+name, value, d).Err()
}

func (c *cacher) Set(name string, value string) error {
	return c.rdb.Set(ctxB, c.prefix+"_"+name, value, c.expiracy).Err()
}

func (c *cacher) Get(name string) (string, error) {
	return c.rdb.Get(ctxB, c.prefix+"_"+name).Result()
}

func (c *cacher) Delete(name string) error {
	return c.rdb.Del(ctxB, c.prefix+"_"+name).Err()
}

func (c *cacher) Ping() error {
	return c.rdb.Ping(ctxB).Err()
}

func (c *cacher) GetKeysWithParam(name string) ([]string, error) {
	return c.rdb.Keys(ctxB, c.prefix+"_"+name).Result()
}

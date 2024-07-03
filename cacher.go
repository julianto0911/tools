package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctxB = context.Background()

type RedisConfiguration struct {
	Host     string
	Port     string
	Password string
	Prefix   string
	UseMock  bool
}

func NewRedisClient(url, port, password string, dbIndex int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     url + ":" + port,
		Password: password,
		DB:       dbIndex,
	})
}

type Cacher struct {
	rdb       *redis.Client
	expiracy  time.Duration
	prefix    string
	registers []string
}

func NewCacher(rdc *redis.Client, prefix string, expiracy int) Cacher {
	return Cacher{
		rdb:      rdc,
		expiracy: time.Duration(expiracy) * time.Second,
		prefix:   prefix,
	}
}

func (c *Cacher) PrintKeys() {
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

func (c *Cacher) SetWithDuration(name string, value string, d time.Duration) error {
	c.addRegister(name)
	return c.rdb.Set(ctxB, c.prefix+"_"+name, value, d).Err()
}

func (c *Cacher) Set(name string, value string) error {
	c.addRegister(name)
	return c.rdb.Set(ctxB, c.prefix+"_"+name, value, c.expiracy).Err()
}

func (c *Cacher) Get(name string) (string, error) {
	c.addRegister(name)
	return c.rdb.Get(ctxB, c.prefix+"_"+name).Result()
}

func (c *Cacher) Delete(name string) error {
	c.release()
	return c.rdb.Del(ctxB, c.prefix+"_"+name).Err()
}

func (c *Cacher) Release() {
	c.release()
}

func (c *Cacher) Ping() error {
	ctx := context.Background()
	return c.rdb.Ping(ctx).Err()
}

func (c *Cacher) addRegister(name string) {
	for _, row := range c.registers {
		if row == name {
			return
		}
	}
	c.registers = append(c.registers, name)
}

func (c *Cacher) release() {
	for _, row := range c.registers {
		c.rdb.Del(ctxB, c.prefix+"_"+row)
	}

	c.registers = []string{}
}

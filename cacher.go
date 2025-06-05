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

// this is useful for testing, to predefined behavior of the response
type cacherResponse struct {
	Key   string
	Type  string
	Value string
	Error error
}

type Cacher struct {
	rdb       *redis.Client
	expiracy  time.Duration
	prefix    string
	registers []string
	responses map[string]cacherResponse
	forced    bool  //force to send error for any kind of request
	err       error //error info when forced is true
}

func NewCacher(rdc *redis.Client, prefix string, expiracy int) Cacher {
	return Cacher{
		rdb:       rdc,
		expiracy:  time.Duration(expiracy) * time.Second,
		prefix:    prefix,
		responses: make(map[string]cacherResponse),
	}
}

func (c *Cacher) SetForcedError(err error) {
	c.forced = true
	c.err = err
}

func (c *Cacher) SetResponse(key, tipe, value string, err error) {
	c.responses[key+"_"+tipe] = cacherResponse{
		Key:   key,
		Value: value,
		Error: err,
		Type:  tipe,
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
	if c.forced {
		c.forced = false
		err := c.err
		c.err = nil
		return err
	}

	if v, exist := c.responses[name+"_set"]; exist {
		//remove registered response
		delete(c.responses, name+"_set")
		return v.Error
	}

	c.addRegister(name)
	return c.rdb.Set(ctxB, c.prefix+"_"+name, value, d).Err()
}

func (c *Cacher) Set(name string, value string) error {
	if c.forced {
		c.forced = false
		err := c.err
		c.err = nil
		return err
	}

	if v, exist := c.responses[name+"_set"]; exist {
		//remove registered response
		delete(c.responses, name+"_set")
		return v.Error
	}

	c.addRegister(name)
	return c.rdb.Set(ctxB, c.prefix+"_"+name, value, c.expiracy).Err()
}

func (c *Cacher) Get(name string) (string, error) {
	if c.forced {
		c.forced = false
		err := c.err
		c.err = nil
		return "", err
	}

	if v, exist := c.responses[name+"_get"]; exist {
		//remove registered response
		delete(c.responses, name+"_get")
		return v.Value, v.Error
	}

	c.addRegister(name)
	return c.rdb.Get(ctxB, c.prefix+"_"+name).Result()
}

func (c *Cacher) Delete(name string) error {
	if c.forced {
		c.forced = false
		err := c.err
		c.err = nil
		return err
	}

	if v, exist := c.responses[name+"_delete"]; exist {
		//remove registered response
		delete(c.responses, name+"_delete")
		return v.Error
	}

	return c.rdb.Del(ctxB, c.prefix+"_"+name).Err()
}

func (c *Cacher) Release() {
	c.release()
}

func (c *Cacher) Ping() error {
	if c.forced {
		c.forced = false
		err := c.err
		c.err = nil
		return err
	}
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

func (c *Cacher) GetKeysWithParam(name string) ([]string, error) {
	if c.forced {
		c.forced = false
		err := c.err
		c.err = nil
		return nil, err
	}
	if v, exist := c.responses[name]; exist {
		//remove registered response
		delete(c.responses, name)
		return []string{v.Value}, v.Error
	}

	c.addRegister(name)
	return c.rdb.Keys(ctxB, c.prefix+"_"+name).Result()
}

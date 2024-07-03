package tools

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

func TestCacher(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("fail init mock cacher: %s", err)
	}
	mr.StartAddr("127.0.0.1:6379")

	client := NewRedisClient("127.0.0.1", "6379", "", 0)

	// client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	// _, err = miniredis.Run()
	// if err != nil {
	// 	t.Fatalf("fail run mock cacher: %s", err)
	// }

	r := NewCacher(client, "", 60)
	err = r.SetWithDuration("test", "test", 1*time.Second)
	assert.Nil(t, err, "should nil")

	err = r.Set("test", "test")
	assert.Nil(t, err, "should nil")

	val, err := r.Get("test")
	assert.Nil(t, err, "should nil")
	assert.Equal(t, val, "test", "should have value")

	err = r.Delete("test")
	assert.Nil(t, err, "should nil")

	err = r.Ping()
	assert.Nil(t, err, "should nil")

	r.Release()
}

func TestMockRedis(t *testing.T) {
	rds, err := MockRedis()
	assert.Nil(t, err)

	rds.Set(context.Background(), "test", "test", 5*time.Second)

	value, err := rds.Get(context.Background(), "test").Result()
	assert.Nil(t, err)
	assert.Equal(t, "test", value)
}

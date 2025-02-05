package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedis(db int, addr, username, password string, ttl time.Duration) *Redis {
	if addr == "" {
		log.Println("Redis: NewRedis: No address provided, skipping redis config")
		return nil
	}

	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
		DB:       db,
	})

	status := c.Ping(context.Background())
	if status.Err() != nil {
		log.Fatal("Redis: NewRedis: Error connecting to redis:", status.Err())
	}
	return &Redis{
		client: c,
		ttl:    ttl,
	}
}

func (r *Redis) Get(ctx context.Context, key string) (interface{}, error) {
	v, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return "", err
	}

	var data interface{}
	dec := gob.NewDecoder(bytes.NewReader(v))
	if err := dec.Decode(&data); err != nil {
		log.Println("Redis: Get: Error decoding value:", err)
		return nil, err
	}
	return data, nil
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}) error {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(value); err != nil {
		log.Println("Redis: Set: Error encoding value:", err)
		return err
	}

	return r.client.Set(ctx, key, b.Bytes(), r.ttl).Err()
}

func (r *Redis) RemoveAll(ctx context.Context) error {
	return r.client.FlushDB(ctx).Err()
}

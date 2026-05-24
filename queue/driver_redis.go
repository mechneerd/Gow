package queue

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"

	"github.com/mechneerd/gow/database/redis"
)

// RedisDriver is a queue driver backed by Redis.
type RedisDriver struct {
	client *redis.Client
	queue  string
	ctx    context.Context
}

// NewRedisDriver creates a new Redis queue driver.
func NewRedisDriver(client *redis.Client, defaultQueueName string) *RedisDriver {
	return &RedisDriver{
		client: client,
		queue:  defaultQueueName,
		ctx:    context.Background(),
	}
}

// Push pushes a job onto the Redis queue list.
func (d *RedisDriver) Push(job Job) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	
	// Note: The concrete job type must be registered with gob.Register() elsewhere.
	err := enc.Encode(&job)
	if err != nil {
		return err
	}

	return d.client.RawClient().LPush(d.ctx, d.queue, buf.Bytes()).Err()
}

// Pop pops a job off the Redis queue list.
func (d *RedisDriver) Pop() (Job, error) {
	res, err := d.client.RawClient().BRPop(d.ctx, 0, d.queue).Result()
	if err != nil {
		return nil, err
	}

	if len(res) < 2 {
		return nil, errors.New("empty pop result")
	}

	payload := res[1]

	var job Job
	buf := bytes.NewBuffer([]byte(payload))
	dec := gob.NewDecoder(buf)
	
	err = dec.Decode(&job)
	if err != nil {
		return nil, err
	}

	return job, nil
}


package queue

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const workflowQueue = "workflow:jobs"

type RedisQueue struct {
	client *redis.Client
}

func NewRedisQueue(client *redis.Client) *RedisQueue {
	return &RedisQueue{
		client: client,
	}
}

func (q *RedisQueue) Push(ctx context.Context, payload []byte) error {
	return q.client.LPush(ctx, workflowQueue, payload).Err()
}

func (q *RedisQueue) Pop(ctx context.Context) ([]byte, error) {
	result, err := q.client.BRPop(ctx, 5*time.Second, workflowQueue).Result()
	if err != nil {
		return nil, err
	}

	return []byte(result[1]), nil
}
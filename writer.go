package log

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap/zapcore"
)

var (
	_ zapcore.WriteSyncer = new(RedisWriter)
)

type RedisWriter struct {
	cli     *redis.Client
	listKey string
}

func NewRedisWriter(key string, cli *redis.Client) *RedisWriter {
	return &RedisWriter{
		cli: cli, listKey: key,
	}
}

func (w *RedisWriter) Write(p []byte) (int, error) {
	n, err := w.cli.RPush(context.Background(), w.listKey, p).Result()
	return int(n), err
}

func (w *RedisWriter) Sync() error {
	return nil
}

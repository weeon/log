package log

import (
	"github.com/go-redis/redis/v7"
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
	n, err := w.cli.RPush(w.listKey, p).Result()
	return int(n), err
}

func (w *RedisWriter) Sync() error {
	return nil
}

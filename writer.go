package log

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap/zapcore"
)

var (
	_ zapcore.WriteSyncer = new(RedisWriter)
)

type RedisWriter struct {
	cli *redis.Client
	key string
}

func NewRedisWriter(key string, cli *redis.Client) *RedisWriter {
	return &RedisWriter{
		cli: cli, key: key,
	}
}

func (w *RedisWriter) Write(p []byte) (int, error) {
	n, err := w.cli.Publish(context.Background(), w.key, p).Result()
	if err != nil {
		fmt.Println("log redis publish error:", err)
	}
	return int(n), err
}

func (w *RedisWriter) Sync() error {
	return nil
}

// nats
type NatsWriter struct {
	nc      *nats.Conn
	subject string
}

func NewNatsWriter(natsURL string, subject string) *NatsWriter {
	nc, _ := nats.Connect(natsURL)
	return &NatsWriter{
		nc:      nc,
		subject: subject,
	}
}

func (n *NatsWriter) Write(p []byte) (int, error) {
	err := n.nc.Publish(n.subject, p)
	return len(p), err
}

func (n *NatsWriter) Sync() error {
	return nil
}

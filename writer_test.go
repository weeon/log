package log

import (
	"os"
	"testing"
)

func TestNats(t *testing.T) {
	natsURI := os.Getenv("NATS_URI")
	NatsWriter := NewNatsWriter(natsURI, "test")
	i, err := NatsWriter.Write([]byte("test"))
	if err != nil {
		t.Error(err)
	}
	t.Log(i)
}

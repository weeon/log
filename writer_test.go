package log

import "testing"

func TestNats(t *testing.T) {
	NatsWriter := NewNatsWriter("nats://localhost:4222", "test")
	i, err := NatsWriter.Write([]byte("test"))
	if err != nil {
		t.Error(err)
	}
	t.Log(i)
}

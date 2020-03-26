package pkg

import (
	"github.com/go-redis/redis/v7"
	"log"
)

type RedisPubSub struct {
	client *redis.Client
	pubsub *redis.PubSub
	name   string
}

func NewRedisPubSub(addr string, name string) *RedisPubSub {
	return &RedisPubSub{
		client: redis.NewClient(&redis.Options{Addr: addr}),
		pubsub: nil,
		name:   name,
	}
}

func (rps *RedisPubSub) Close() {
	if rps.client != nil {
		_ = rps.client.Close()
	}
	if rps.pubsub != nil {
		_ = rps.pubsub.Close()
	}
}

func (rps *RedisPubSub) Subscribe() chan []byte {
	rps.pubsub = rps.client.Subscribe(rps.name)

	out := make(chan []byte)

	c := rps.pubsub.Channel()
	go func() {
		for msg := range c {
			log.Printf("RedisPubSub.Subscribe, msg:%s", msg)
			out <- []byte(msg.Payload)
		}
		close(out)
	}()

	return out
}

func (rps *RedisPubSub) Publish(data []byte) error {
	return rps.client.Publish(rps.name, data).Err()
}

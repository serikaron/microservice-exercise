package pkg

type PubSub interface {
	Publish(data []byte) error
	Subscribe() chan []byte
	Close()
}

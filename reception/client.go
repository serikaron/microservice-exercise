package reception

import (
	"io"
	"io/ioutil"
	"log"
	"time"
)

type Client struct {
	id         uint32
	readWriter io.ReadWriter
	writeChan  chan []byte
	done       chan bool
}

func NewClient(clientId uint32, readWriter io.ReadWriter) *Client {
	return &Client{
		id:         clientId,
		readWriter: readWriter,
		writeChan:  make(chan []byte),
		done:       make(chan bool),
	}
}

func (client *Client) Close() {
	close(client.writeChan)
	close(client.done)
}

func (client *Client) Start(readChan chan []byte) {
	for {
		select {
		case <-client.done:
			return
		case buf := <-client.writeChan:
			n, err := client.readWriter.Write(buf)
			if err != nil {
				log.Println("client.Start() write failed: ", err)
				break
			}
			if n < len(buf) {
				log.Printf("client.Start() write not complete, buf len:%d, actually write:%d\n", len(buf), n)
				break
			}
		default:
			buf, err := ioutil.ReadAll(client.readWriter)
			if err != nil {
				log.Println("client.Start() read failed: ", err)
				return
			}
			if len(buf) > 0 {
				readChan <- buf
			} else {
				log.Println("client.Start() read empty buf")
			}
		}

		time.Sleep(time.Millisecond * 10)
	}
}

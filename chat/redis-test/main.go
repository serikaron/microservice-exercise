package main

import (
	"chat/pkg"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	rdsHost := flag.String("redis-host", "redis", "redis host")
	rdsPort := flag.Uint("redis-port", 3697, "redis port")
	flag.Parse()
	rdsAddr := fmt.Sprintf("%s:%d", *rdsHost, *rdsPort)
	log.Println("connect to redis:", rdsAddr)
	rdsPS := pkg.NewRedisPubSub(rdsAddr, "notify")

	c := rdsPS.Subscribe()

	wg := sync.WaitGroup{}
	done := make(chan bool)

	wg.Add(1)
	go func() {
		//defer wg.Done()
		for {
			select {
			case <-done:
				return
			case msg := <-c:
				log.Println("msg.byte:", msg)
				log.Println("msg:", string(msg))
				wg.Done()
			default:
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	//for {
	log.Println("publish:", []byte("test"))
	err := rdsPS.Publish([]byte("test"))
	if err != nil {
		log.Println("publish failed:", err)
	}
	//time.Sleep(time.Second)
	//}
	wg.Wait()
	close(done)
}

package main

import (
	"log"
	"mse/proto"
)

type listener struct {
	name   string
	stream proto.Chat_ListenServer
	done   chan error
}

type listenerHub struct {
	listenerMap map[string]*listener
	addChan     chan *listener
	notifyChan  chan *proto.ListenRsp
}

func newListenerHub() *listenerHub {
	return &listenerHub{
		listenerMap: make(map[string]*listener, 10),
		addChan:     make(chan *listener, 10),
		notifyChan:  make(chan *proto.ListenRsp, 10),
	}
}

func (lh *listenerHub) run(done chan bool) {
	defer log.Println("listenerHub.Run done")
	for {
		select {
		case <-done:
			for _, l := range lh.listenerMap {
				l.done <- nil
			}
			return
		case l := <-lh.addChan:
			log.Printf("listenerHub.Run add listener, name:%s", l.name)
			lh.listenerMap[l.name] = l
		case message := <-lh.notifyChan:
			log.Printf("listenerHub.Run notify message, message:%s", message)
			lh.send(message)
			//default:
			//	time.Sleep(10 * time.Millisecond)
		}
		log.Println("listenerHub Running...")
	}
}

func (lh *listenerHub) send(message *proto.ListenRsp) {
	for name, l := range lh.listenerMap {
		err := l.stream.Send(message)
		if err != nil {
			log.Printf("listenerHub.notify error:%v", err)
			if l.done != nil {
				l.done <- err
			}
			delete(lh.listenerMap, name)
		}
	}
}

func (lh *listenerHub) addListener(l *listener) {
	lh.addChan <- l
}

func (lh *listenerHub) notify(message *proto.ListenRsp) {
	lh.notifyChan <- message
}

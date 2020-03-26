package main

import (
	"log"
)

type Notifier interface {
	Notify(stream interface{}, message interface{}) error
}

type Listener struct {
	name   string
	stream interface{}
	done   chan error
}

type ListenerHub struct {
	listenerMap map[string]*Listener
	addChan     chan *Listener
	notifyChan  chan interface{}
	notifier    Notifier
}

func NewListenerHub(notifier Notifier) *ListenerHub {
	return &ListenerHub{
		listenerMap: make(map[string]*Listener, 10),
		addChan:     make(chan *Listener, 10),
		notifyChan:  make(chan interface{}, 10),
		notifier:    notifier,
	}
}

func (lh *ListenerHub) Run(done chan bool) {
	defer log.Println("ListenerHub.Run done")
	for {
		select {
		case <-done:
			for _, l := range lh.listenerMap {
				l.done <- nil
			}
			return
		case l := <-lh.addChan:
			log.Printf("ListenerHub.Run add listener, name:%s", l.name)
			lh.listenerMap[l.name] = l
		case message := <-lh.notifyChan:
			log.Printf("ListenerHub.Run notify message, message:%s", message)
			lh.notify(message)
			//default:
			//	time.Sleep(10 * time.Millisecond)
		}
		log.Println("ListenerHub Running...")
	}
}

func (lh *ListenerHub) notify(message interface{}) {
	for name, l := range lh.listenerMap {
		err := lh.notifier.Notify(l.stream, message)
		if err != nil {
			log.Printf("ListenerHub.notify error:%v", err)
			if l.done != nil {
				l.done <- err
			}
			delete(lh.listenerMap, name)
		}
	}
}

func (lh *ListenerHub) AddListener(l *Listener) {
	lh.addChan <- l
}

func (lh *ListenerHub) Notify(message interface{}) {
	lh.notifyChan <- message
}

package internal

import "sync"

type SafeHub struct {
	hub map[string]*Listener
	mux sync.Mutex
}

func NewSafeHub() *SafeHub {
	return &SafeHub{hub: make(map[string]*Listener)}
}

func (sh *SafeHub) Close() {
	sh.mux.Lock()
	defer sh.mux.Unlock()
	for _, l := range sh.hub {
		l.close()
	}
}

func (sh *SafeHub) Add(l *Listener) {
	sh.mux.Lock()
	defer sh.mux.Unlock()
	sh.hub[l.name] = l
}

func (sh *SafeHub) Remove(key string) {
	sh.mux.Lock()
	defer sh.mux.Unlock()
	if l, ok := sh.hub[key]; ok {
		delete(sh.hub, key)
		l.close()
	}
}

func (sh *SafeHub) Notify(msg string) {
	sh.mux.Lock()
	defer sh.mux.Unlock()
	for _, l := range sh.hub {
		l.c <- msg
	}
}

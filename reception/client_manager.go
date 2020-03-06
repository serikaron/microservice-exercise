package reception

import "io"

type ClientManager struct {
	clientId  uint32
	clientMap map[uint32]*Client
	readChan  chan []byte
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clientId:  0,
		clientMap: make(map[uint32]*Client),
		readChan:  make(chan []byte),
	}
}

func (cm *ClientManager) newClient(rw io.ReadWriter) {
	cm.clientId++

	client := NewClient(cm.clientId, rw)
	cm.clientMap[cm.clientId] = client
	go client.Start(cm.readChan)
}

func (cm *ClientManager) getClient(clientId uint32) *Client {
	return cm.clientMap[clientId]
}

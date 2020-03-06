package reception

import (
	"fmt"
	"testing"
)

func TestClientManager(t *testing.T) {
	mgr := NewClientManager()

	checkNewClient := func(t *testing.T, clientId uint32) {
		t.Helper()

		mockReadWriter := NewMockReaderWriter()
		mgr.newClient(mockReadWriter)
		client := mgr.getClient(clientId)
		if client == nil {
			t.Errorf("client not found, clientId:%d", clientId)
			return
		}
		if client.id != clientId {
			t.Errorf("client id not match, Got:%d, Want:%d", client.id, clientId)
		}
	}

	clientId := uint32(0)
	clientId++
	t.Run(fmt.Sprintf("new client %d", clientId), func(t *testing.T) {
		checkNewClient(t, clientId)
	})
	clientId++
	t.Run(fmt.Sprintf("new client %d", clientId), func(t *testing.T) {
		checkNewClient(t, clientId)
	})
	t.Fatal("not implement")
}

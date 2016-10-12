package sockets

import (
	"net"
	"sync"

	"github.com/gorilla/websocket"
)

var ActiveClients = make(map[string]map[ClientConn]int)
var ActiveClientsRWMutex sync.RWMutex

type ClientConn struct {
	RoomID    string
	Websocket *websocket.Conn
	ClientIP  net.Addr
}

func AddClient(cc ClientConn) {
	ActiveClientsRWMutex.Lock()
	if ActiveClients[cc.RoomID] == nil {
		ActiveClients[cc.RoomID] = make(map[ClientConn]int)
	}
	ActiveClients[cc.RoomID][cc] = 0
	ActiveClientsRWMutex.Unlock()
}

func DeleteClient(cc ClientConn) {
	ActiveClientsRWMutex.Lock()
	delete(ActiveClients[cc.RoomID], cc)
	ActiveClientsRWMutex.Unlock()
}

func BroadcastMessage(messageType int, message []byte, roomID string) {
	ActiveClientsRWMutex.RLock()
	defer ActiveClientsRWMutex.RUnlock()

	for client, _ := range ActiveClients[roomID] {
		if err := client.Websocket.WriteMessage(messageType, message); err != nil {
			return
		}
	}
}

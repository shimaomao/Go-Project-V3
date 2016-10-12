package app

import (
	"log"
	"net"
	"sync"

	"github.com/gorilla/websocket"
)

var ActiveClients = make(map[string]map[ClientConn]int)
var ActiveClientsRWMutex sync.RWMutex

type ClientConn struct {
	RoomID    string
	websocket *websocket.Conn
	clientIP  net.Addr
}

func addClient(cc ClientConn) {
	ActiveClientsRWMutex.Lock()
	if ActiveClients[cc.RoomID] == nil {
		ActiveClients[cc.RoomID] = make(map[ClientConn]int)
	}
	ActiveClients[cc.RoomID][cc] = 0
	ActiveClientsRWMutex.Unlock()
}

func deleteClient(cc ClientConn) {
	ActiveClientsRWMutex.Lock()
	delete(ActiveClients[cc.RoomID], cc)
	ActiveClientsRWMutex.Unlock()
}

func broadcastMessage(messageType int, message []byte, roomID string) {
	ActiveClientsRWMutex.RLock()
	defer ActiveClientsRWMutex.RUnlock()

	for client, _ := range ActiveClients[roomID] {
		if err := client.websocket.WriteMessage(messageType, message); err != nil {
			return
		}
	}
}

func broadcastJson(message interface{}, roomID string) {
	ActiveClientsRWMutex.RLock()
	defer ActiveClientsRWMutex.RUnlock()
	log.Println("sending message to room: %#s, %s", message, roomID)
	for client, _ := range ActiveClients[roomID] {
		if err := client.websocket.WriteJSON(message); err != nil {
			log.Println("err", err)
			return
		}
	}
}

func broadcastMessageToAll(messageType int, message []byte) {
	ActiveClientsRWMutex.RLock()
	defer ActiveClientsRWMutex.RUnlock()

	for _, clients := range ActiveClients {
		for client, _ := range clients {
			if err := client.websocket.WriteMessage(messageType, message); err != nil {
				log.Println("err", err)
				return
			}
		}
	}
}

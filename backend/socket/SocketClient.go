package socket

import (
	"github.com/gofiber/websocket/v2"
	"log"
)

type Client struct {
	conn *websocket.Conn
	key  string
}

var connToKeyMap = make(map[*websocket.Conn]string)
var keyToConnMap = make(map[string]*websocket.Conn)

var registerChan = make(chan Client)
var unregisterChan = make(chan *websocket.Conn)
var broadcastChan = make(chan Message)

func RunSocketHub() {
	for {
		select {
		case client := <-registerChan:
			connToKeyMap[client.conn] = client.key
			keyToConnMap[client.key] = client.conn

		case conn := <-unregisterChan:
			key := connToKeyMap[conn]
			delete(connToKeyMap, conn)
			delete(keyToConnMap, key)

		case message := <-broadcastChan:
			conn, found := keyToConnMap[message.Key]
			if found {
				if err := conn.WriteJSON(message.Content); err != nil {
					unregisterChan <- conn
					_ = conn.Close()
				}
			}
		}
	}
}

func CreateNewSocketConnection(conn *websocket.Conn) {
	// When the function returns, unregister the client and close the connection
	defer func() {
		log.Println("Socket client disconnected, closing the connection.")
		unregisterChan <- conn
		_ = conn.Close()
	}()

	// Register the client, registerChan <- conn
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("read error:", err)
			}
			return // Calls the deferred function, i.e. closes the connection on error
		}

		if messageType == websocket.TextMessage {
			log.Println("Socket client connected through the socket.")
			newClient := Client{
				conn: conn,
				key:  string(message),
			}
			registerChan <- newClient
		} else {
			log.Println("Socket message received of type", messageType)
		}
	}
}

func broadcastMessage(message Message) {
	broadcastChan <- message
}

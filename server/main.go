package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var connectionPool map[string]map[string][]*websocket.Conn
var connectionPoolMutex sync.Mutex

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}
	user := r.URL.Query().Get("username")
	room := r.URL.Query().Get("room")

	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	go reader(ws, room, user)
}

type Message struct {
	User    string `json:"user"`
	Message string `json:"message"`
	Room    string `json:"room"`
}

func reader(conn *websocket.Conn, room string, username string) {
	connectionPoolMutex.Lock()
	if _, ok := connectionPool[room]; !ok {
		connectionPool[room] = make(map[string][]*websocket.Conn)
	}
	if _, ok := connectionPool[room][username]; !ok {
		connectionPool[room][username] = []*websocket.Conn{}
	}
	connectionPool[room][username] = append(connectionPool[room][username], conn)
	connectionPoolMutex.Unlock()
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			// remove connection from pool
			for i, c := range connectionPool[room][username] {
				if c == conn {
					connectionPool[room][username] = append(connectionPool[room][username][:i], connectionPool[room][username][i+1:]...)
				}
			}
			return
		}
		log.Println(string(p), messageType)
		msg := Message{}
		err = json.Unmarshal(p, &msg)
		if err == nil {
			// fannout message to all connections in the room
			for k, c := range connectionPool[msg.Room] {
				if k != msg.User {
					for _, connie := range c {
						if err := connie.WriteMessage(messageType, p); err != nil {
							log.Println(err)
							return
						}
					}
				}
			}
		}
	}
}

func main() {

	connectionPool = make(map[string]map[string][]*websocket.Conn)
	http.HandleFunc("/ws", wsEndpoint)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

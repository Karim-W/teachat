package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type client struct {
	conn     *websocket.Conn
	callBack func([]byte)
}

func InitOrDie(callback func([]byte)) *client {
	url := _WEB_SOCKET_URL
	url += "?username=" + *USERNAME
	// url += "&password=" + *PASSWORD
	url += "&room=" + *ROOM_NAME
	// url += "&roompass=" + *ROOM_PASS
	c, resp, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		log.Printf("handshake failed with status %d", resp.StatusCode)
		log.Fatal("dial:", err)
	}

	//When the program closes close the connection
	// defer c.Close()
	cl := &client{
		conn:     c,
		callBack: callback,
	}
	go cl.ManageConnection()
	return cl
}

func (c *client) ManageConnection() {
	authPayload := map[string]string{
		"username": *USERNAME,
		"password": *PASSWORD,
		"room":     *ROOM_NAME,
		"roompass": *ROOM_PASS,
	}
	c.conn.WriteJSON(authPayload)
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		c.callBack(message)
	}
}

func (c *client) SendMessageJSON(msg interface{}) {
	c.conn.WriteJSON(msg)
}

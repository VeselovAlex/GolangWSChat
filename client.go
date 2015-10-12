// client
package main

import (
	"encoding/json"
	"log"
	"time"

	ws "github.com/VeselovAlex/GolangWSChat/Godeps/_workspace/src/github.com/gorilla/websocket"
)

type Client struct {
	Conn *ws.Conn
	Msg  chan *Message
	Room *Room

	Name string
}

func (c *Client) Read() {
	for {
		/*While connection is active*/
		if _, msg, err := c.Conn.ReadMessage(); err == nil {
			wrap := NewMessage(c.Name, msg)
			c.Room.Send <- wrap
		} else {
			break
		}
	}
	c.Conn.Close()
	log.Println("User", c.Name, "connection closed", "(reading)")
}

func (c *Client) Write() {
	ticker := time.NewTicker(50 * time.Second)
	var err error = nil
	for {
		select {
		case msg := <-c.Msg:
			/*Send incoming messages to remote receiver*/
			jsonMsg, err := json.Marshal(msg)
			if err != nil {
				break
			}
			err = c.Conn.WriteMessage(ws.TextMessage, jsonMsg)
		case <-ticker.C:
			//Keep connection established
			err = c.Conn.WriteMessage(ws.PongMessage, nil)
		}
		if err != nil {
			break
		}
	}
	c.Conn.Close()
	ticker.Stop()
	log.Println("User", c.Name, "connection closed", "(writing)")
}

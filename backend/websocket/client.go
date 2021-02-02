package websocket

import (
	"log"

	"github.com/gorilla/websocket"
)

// Client is the websocket client type
type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

// Message is the type which defines the messages sent back to the Websocket
type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

func (c *Client) Read() (Message, error) {
	/*defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()
	*/

	var message Message

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return message, err
		}
		message = Message{Type: messageType, Body: string(p)}
		// Below sends the message back to all clients
		//c.Pool.Broadcast <- message
		//fmt.Printf("Message Received: %+v\n", message)
		return message, nil
	}
}

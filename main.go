package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

type EventType string

var EventTypeJoin EventType = "join"
var EventTypeMessage EventType = "message"

type Message struct {
	Type      EventType `json:"type"`
	UserId    string    `json:"userId"`
	Content   string    `json:"content"`
	Timestamp string    `json:"timestamp"`
}

func (m Message) Serialize() []byte {
	jsonMessage, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return jsonMessage
}

func readClient(ws *websocket.Conn, messageChannel chan<- Message) {
	fmt.Printf("Start reading client: %v", ws.RemoteAddr())
	for {
		var message Message
		err := websocket.JSON.Receive(ws, &message)
		if err != nil && err != io.EOF {
			fmt.Printf("Failed to receive: %v\n", err)
			break
		}

		switch message.Type {
		case EventTypeMessage:
			fmt.Printf("Received: %v", message.Content)
		}
		messageChannel <- message
	}
}

func serveClient(clientsChannel <-chan *websocket.Conn, messageChannel chan<- Message) {
	for {
		select {
		case clientSock := <-clientsChannel:
			fmt.Println("Client connected")
			go readClient(clientSock, messageChannel)
		}
	}
}

type WebsocketServer struct {
	ClientChannel  chan *websocket.Conn
	MessageChannel chan Message
}

// Echo the data received on the WebSocket.
func (s WebsocketServer) EchoServer(ws *websocket.Conn) {
	s.ClientChannel <- ws
	go serveClient(s.ClientChannel, s.MessageChannel)
}

func main() {
	messageChannel := make(chan Message)
	clientChannels := make(chan *websocket.Conn)
	s := WebsocketServer{
		ClientChannel:  clientChannels,
		MessageChannel: messageChannel,
	}

	http.Handle("/", websocket.Handler(s.EchoServer))
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			panic("ListenAndServe: " + err.Error())
		}
	}()

	for {
		select {
		case message := <-messageChannel:
			fmt.Print("Message received:", message.Content, " At:", message.Timestamp)
		}
	}
}

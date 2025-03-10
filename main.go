package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/shakram02/nony-chat/adapters/http/handshaker"
	http_parser "github.com/shakram02/nony-chat/adapters/http/parser"
	"github.com/shakram02/nony-chat/adapters/websocket"
)

type EventType string

const BufferSize = 2048

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

func readMessage(socket net.Conn) ([]byte, error) {
	buffer := make([]byte, 2048)
	n, err := socket.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:n], nil
}

type ChatClient struct {
	id     string
	socket net.Conn
}

func initiateConnection(socket net.Conn, connections <-chan ChatClient) {
	// Read the handshake
	message, err := readMessage(socket)
	if err != nil {
		fmt.Printf("Failed to read handshake: %v", err)
		return
	}
	websocketHandshake, err := http_parser.Parse(message)
	handshaker := handshaker.New(socket, websocketHandshake)

	if err != nil {
		fmt.Printf("Failed to parse handshake: %v", err)
		err := handshaker.Reject()
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		// Reply to handshake
		err := handshaker.Handshake()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	for {
		// Read welcome message
		message, err = readMessage(socket)
		if err != nil {
			fmt.Printf("Failed to read: %v", err)
			return
		}

		packet := websocket.New(message)
		fmt.Printf("%s\n", packet)
	}

}

func main() {
	http.Handle("/", http.FileServer(http.Dir("public")))
	go http.ListenAndServe(":8000", nil)

	server, err := net.Listen("tcp", ":8080")
	chatClientChan := make(chan ChatClient)
	if err != nil {
		panic(fmt.Errorf("Failed to listen: %s", err))
	}

	for {
		conn, err := server.Accept()
		if err != nil {
			panic(fmt.Errorf("Failed to accept: %s", err))
		}

		go initiateConnection(conn, chatClientChan)
	}
}

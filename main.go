package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/shakram02/nony-chat/adapters/http/handshaker"
	http_parser "github.com/shakram02/nony-chat/adapters/http/parser"
	"github.com/shakram02/nony-chat/adapters/nony"
	"github.com/shakram02/nony-chat/adapters/websocket"
)

const BufferSize = 2048

func readSocketBytes(socket net.Conn) ([]byte, error) {
	buffer := make([]byte, BufferSize)
	n, err := socket.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:n], nil
}

var ErrInvalidFrame = errors.New("Invalid websocket packet")

func readWebsocketFrame(socket net.Conn) (*websocket.WebsocketFrame, error) {
	bytes, err := readSocketBytes(socket)
	if err != nil {
		return nil, err
	}

	frame := websocket.New(bytes)
	// TODO: refactor this function to websocket reader later, maybe handle fragmented packets.
	if frame == nil || frame.IsFragmented() {
		fmt.Printf("Failed to parse welcome packet: %v\n", err)
		return nil, ErrInvalidFrame
	}

	return frame, nil
}

type ChatClient struct {
	socket        net.Conn
	welcomePacket *nony.NonyPacket
}

func initiateConnection(socket net.Conn, connections chan<- ChatClient) {
	// Read the handshake
	welcomeMessage, err := readSocketBytes(socket)
	if err != nil {
		fmt.Printf("Failed to read handshake: %v", err)
		return
	}
	websocketHandshake, err := http_parser.Parse(welcomeMessage)
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

	// Read welcome message
	welcomePacket, err := readWebsocketFrame(socket)
	if err != nil {
		fmt.Printf("Failed to read: %v\n", err)
		return
	}

	packet := nony.New(*welcomePacket)
	if packet == nil {
		fmt.Printf("Failed to parse packet: %v\n", err)
		socket.Close()
		return
	}

	connections <- ChatClient{
		welcomePacket: packet,
		socket:        socket,
	}
}

func handleClient(socket net.Conn) {
	for {
		frame, err := readWebsocketFrame(socket)
		if err != nil {
			fmt.Printf("Failed to read frame: %v\n", err)
			socket.Close()
			return
		}

		packet := nony.New(*frame)
		if packet == nil {
			fmt.Printf("Failed to parse Nony packet: %v\n", err)
			socket.Close()
			return
		}

		fmt.Printf("[%s] // [%s]\n", packet.UserId, packet.Type)
	}
}

func handleConnections(connections chan ChatClient) {
	for client := range connections {
		// Register user
		fmt.Printf("[%s] joined [%s] at %s\n",
			client.welcomePacket.UserId,
			client.welcomePacket.RoomId,
			client.welcomePacket.Timestamp.Format(time.RFC3339),
		)

		go handleClient(client.socket)
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("public")))
	log.Println("HTTP ServerListening on port 8000")
	go func() {
		err := http.ListenAndServe(":8000", nil)
		if err != nil {
			log.Fatalf("Failed to start HTTP server: %s", err)
		}
	}()

	server, err := net.Listen("tcp", ":8080")
	log.Println("TCP Server Listening on port 8080")
	if err != nil {
		panic(fmt.Errorf("Failed to listen: %s", err))
	}

	connections := make(chan ChatClient)

	go handleConnections(connections)

	for {
		conn, err := server.Accept()
		if err != nil {
			panic(fmt.Errorf("Failed to accept: %s", err))
		}

		go initiateConnection(conn, connections)
	}
}

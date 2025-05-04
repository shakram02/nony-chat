package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/shakram02/nony-chat/adapters/nony"
	"github.com/shakram02/nony-chat/adapters/protocol/transport"
)

const BufferSize = 2048

var ErrInvalidFrame = errors.New("Invalid websocket packet")

type ChatClient struct {
	socket        net.Conn
	welcomePacket *nony.Packet
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

	for {
		conn, err := server.Accept()
		if err != nil {
			panic(fmt.Errorf("Failed to accept: %s", err))
		}

		go func() {
			tcpTransport := transport.NewTcp(conn, BufferSize)
			websocketsTransport := transport.NewWebsocket(tcpTransport)
			nonySocket := transport.NewNony(tcpTransport, websocketsTransport)

			err := nonySocket.Start()
			if err != nil {
				panic("failed to handshake client:" + err.Error())
			}

			for {
				packet, err := nonySocket.Read()
				if err != nil {
					panic("failed to read nony packet:" + err.Error())
				}

				if packet == nil {
					nonySocket.Close()
					break
				}

				fmt.Printf("[rx]: %v\n", packet.Content)
			}
		}()
	}
}

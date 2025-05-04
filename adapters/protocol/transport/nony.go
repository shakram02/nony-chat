package transport

import (
	"fmt"

	"github.com/shakram02/nony-chat/adapters/http/handshaker"
	http_parser "github.com/shakram02/nony-chat/adapters/http/parser"
	"github.com/shakram02/nony-chat/adapters/nony"
)

type NonySocket struct {
	tcpTransport       *Tcp
	websocketTransport *Websockets
}

func NewNony(
	tcpTransport *Tcp,
	websocketTransport *Websockets,
) *NonySocket {
	return &NonySocket{
		tcpTransport:       tcpTransport,
		websocketTransport: websocketTransport,
	}
}

func (n *NonySocket) Start() error {
	// Read HTTP upgrade request.
	// Handhshake client
	httpHandshake, err := n.tcpTransport.Read()
	if err != nil {
		n.tcpTransport.Write(handshaker.MakeRejectionResponse())
		n.tcpTransport.Close()
		return fmt.Errorf("Failed to read handshake: %v", err)
	}

	websocketHandshake, err := http_parser.ParseUpgradeRequest(httpHandshake)
	if err != nil {
		n.tcpTransport.Write(handshaker.MakeRejectionResponse())
		n.tcpTransport.Close()
		return fmt.Errorf("Failed to parse upgrade request")
	}

	handshakeResponse := handshaker.MakeAcceptanceResposne(websocketHandshake)
	err = n.tcpTransport.Write(handshakeResponse)
	if err != nil {
		n.tcpTransport.Close()
		return fmt.Errorf("Failed to send client handshake response")
	}

	return nil
}

func (n *NonySocket) Read() (*nony.Packet, error) {
	frame, err := n.websocketTransport.Read()
	if err != nil {
		n.websocketTransport.Close()
		return nil, fmt.Errorf("failed to read websocket packet: %w", err)
	}

	packet := nony.New(frame)
	return packet, nil
}

func (n *NonySocket) Close() {
	n.websocketTransport.Close()
}

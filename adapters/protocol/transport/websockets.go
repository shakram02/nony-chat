package transport

import (
	"fmt"

	"github.com/shakram02/nony-chat/adapters/websockets"
)

type Websockets struct {
	tcpTransport *Tcp
	isHandshaked bool
}

func NewWebsocket(tcpTransport *Tcp) *Websockets {
	return &Websockets{
		tcpTransport: tcpTransport,
		isHandshaked: false,
	}
}

func (w *Websockets) Read() (*websockets.Frame, error) {
	tcpMessage, err := w.tcpTransport.Read()
	if err != nil {
		return nil, err
	}

	// TODO: this is the adapter layer. Do we need that layer?
	frame := websockets.New(tcpMessage)
	if frame == nil {
		w.tcpTransport.Close()
		return nil, fmt.Errorf("corrupt frame: TCP connection closed")
	}

	if frame.IsFragmented() {
		// TODO: handle fragmented frames later.
		// Merge Websockets frame together
		panic("Received fragmented frame")
	}

	return frame, nil
}

func (w *Websockets) Write(frame *websockets.Frame) error {
	panic("not implemented")
}

func (w *Websockets) Close() error {
	return w.tcpTransport.Close()
}

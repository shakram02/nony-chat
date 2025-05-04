package transport

import (
	"fmt"
	"net"
)

type Tcp struct {
	socket     net.Conn
	bufferSize int

	isClosed bool
}

func NewTcp(socket net.Conn, bufferSize int) *Tcp {
	return &Tcp{
		socket:     socket,
		bufferSize: bufferSize,
		isClosed:   false,
	}
}

func (t *Tcp) Read() ([]byte, error) {
	if t.isClosed {
		return nil, fmt.Errorf("Connection closed")
	}

	buffer := make([]byte, t.bufferSize)
	n, err := t.socket.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:n], nil
}

func (t *Tcp) Write(data []byte) error {
	if t.isClosed {
		return fmt.Errorf("Connection closed")
	}

	n, err := t.socket.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write handshake: %w", err)
	}

	if n != len(data) {
		return fmt.Errorf("response not fully written, expected: %d, actual: %d", len(data), n)
	}

	return nil
}

func (t *Tcp) Close() error {
	if t.isClosed {
		return nil
	}

	err := t.socket.Close()
	t.isClosed = true

	return err
}

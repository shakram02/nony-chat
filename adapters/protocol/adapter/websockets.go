package adapter

import "github.com/shakram02/nony-chat/adapters/websockets"

type Websocket struct {
	canReceive bool
}

func (w *Websocket) CanReceive() bool {
	panic("not implemented")
}
func (w *Websocket) Receive([]byte) *websockets.Frame {
	panic("not implemented")
}
func (w *Websocket) Send(frame *websockets.Frame) []byte {
	panic("not implemented")
}

package websocket

import (
	"fmt"
)

type FrameOpCode uint8

const (
	OpContinuationFrame FrameOpCode = iota
	OpTextFrame
	OpBinaryFrame
	_Rx3
	_Rx4
	_Rx5
	_Rx6
	_Rx7
	OpConnectionClose
	OpPing
	OpPong
	_RxB
	_RxC
	_RxD
	_RxE
	_RxF
)

type websocketHeader struct {
	Fin           bool
	OpCode        FrameOpCode
	Mask          bool
	PayloadLength uint64
	MaskingKey    uint32
}

type WebsocketPacket struct {
	websocketHeader
	data []uint8
}

func New(raw []uint8) WebsocketPacket {
	parser := newParser(raw)
	header := parser.parseHeader()

	return WebsocketPacket{
		websocketHeader: header,
	}
}

func (p WebsocketPacket) String() string {
	out := ""
	out += "---------------------\n"
	out += fmt.Sprintf(" Is Fin: %v\n", p.Fin)
	out += fmt.Sprintf(" Op: %v\n", p.OpCode)
	out += fmt.Sprintf(" Is Masked: %v\n", p.Mask)
	out += fmt.Sprintf(" Length: %v\n", p.PayloadLength)
	out += "---------------------\n"
	return out
}

// func parseMask(container [2]uint8) [4]uint8 {

// }

package websocket

import (
	"encoding/binary"
	"fmt"
)

type FrameOpCode byte

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
	data []byte
}

func New(raw []byte) WebsocketPacket {
	return WebsocketPacket{
		websocketHeader: parseHeader(raw),
	}
}

func (p WebsocketPacket) String() string {
	out := ""

	out += fmt.Sprintf("Is Fin: %v\n", p.Fin)
	out += fmt.Sprintf("\tOp: %v\n", p.OpCode)
	out += fmt.Sprintf("\tIs Masked: %v\n", p.Mask)
	out += fmt.Sprintf("\tLength: %d\n", p.PayloadLength)
	out += "---------------------"
	return out
}

func parseHeader(headerBytes []byte) websocketHeader {
	fin := parseBit(headerBytes[0], 0)
	opCode := parseHeaderOpCode(headerBytes[0])
	length := uint64(0)
	// lengthBuffer := [9]byte{}
	// parseHeaderPayloadLength([9]byte(headerBytes[1:]))

	return websocketHeader{
		Fin:           fin,
		OpCode:        opCode,
		PayloadLength: length,
	}
}

func parseHeaderOpCode(container byte) FrameOpCode {
	return FrameOpCode(container & 0x0F)
}
func parseHeaderPayloadLength(container [9]byte) uint64 {
	payloadLen := container[0] & 0x7F
	// Most common case first
	if payloadLen < 126 {
		// if 0-125, that is the payload length
		return uint64(payloadLen)
	}
	if payloadLen == 126 {
		// Extended -> uint16
		// If 126, the following 2 bytes interpreted as a
		// 16-bit unsigned integer are the payload length
		lenContainer := []byte{
			container[1],
			container[2],
		}
		return uint64(binary.BigEndian.Uint16(lenContainer))

	}
	if payloadLen == 127 {
		// Extended -> uint64
		// If 127, the following 8 bytes interpreted as
		// a 64-bit unsigned integer (the most significant
		// bit MUST be 0) are the payload length
		return binary.BigEndian.Uint64(container[1:])
	}
	return 0
}
func parseBit(container byte, index int) bool {
	out := container & (1 << (7 - index))
	if out > 0 {
		return true
	}

	return false
}

package websockets

import (
	"encoding/binary"
	"fmt"
)

type payloadLengthMode uint8

const (
	Simple payloadLengthMode = iota
	Extended16Bits
	Extended64Bits
)

func (plm payloadLengthMode) String() string {
	switch plm {
	case Simple:
		return "Simple"
	case Extended16Bits:
		return "Extended16Bits"
	case Extended64Bits:
		return "Extended64Bits"
	}
	panic("Not implemented")
}

type PayloadLength struct {
	Length uint64
	Mode   payloadLengthMode
}

func parseHeaderPayloadLengthMode(input uint8) payloadLengthMode {
	payloadLen := input & 0x7F
	// Most common case first
	if payloadLen < 126 {
		return Simple
	}

	if payloadLen == 126 {
		return Extended16Bits
	}

	if payloadLen == 127 {
		return Extended64Bits
	}

	panic(fmt.Sprintf("Invalid payload length: %d", payloadLen))
}

func (p *FrameParser) parseHeaderPayloadLength(mode payloadLengthMode) uint64 {
	payloadLenPointer := p.pointer
	switch mode {
	case Simple:
		// if 0-125, that is the payload length
		return uint64(p.getCurrentByte())
	case Extended16Bits:
		// Extended -> uint16
		// If 126, the following 2 bytes interpreted as a
		// 16-bit unsigned integer are the payload length
		lenContainer := []uint8{
			p.raw[payloadLenPointer+1],
			p.raw[payloadLenPointer+2],
		}
		return uint64(binary.BigEndian.Uint16(lenContainer))
	case Extended64Bits:
		// Extended -> uint64
		// If 127, the following 8 bytes interpreted as
		// a 64-bit unsigned integer (the most significant
		// bit MUST be 0) are the payload length
		return binary.BigEndian.Uint64([]uint8{
			p.raw[payloadLenPointer+1],
			p.raw[payloadLenPointer+2],
			p.raw[payloadLenPointer+3],
			p.raw[payloadLenPointer+4],
			p.raw[payloadLenPointer+5],
			p.raw[payloadLenPointer+6],
			p.raw[payloadLenPointer+7],
			p.raw[payloadLenPointer+8],
		})
	}
	panic(fmt.Sprintf("Invalid payload length mode"))
}

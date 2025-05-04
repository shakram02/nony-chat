package websockets

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
	IsMasked      bool
	Mask          [4]byte
	PayloadLength uint64
	MaskingKey    uint32
}

type Frame struct {
	raw    []uint8
	header websocketHeader
	Data   []uint8
}

// TODO: rename this to parse or something, return result and error.
func New(raw []uint8) *Frame {
	parser := newParser(raw)
	frame := parser.parseFrame()
	if frame.header.IsMasked {
		unmask(frame.Data, frame.header.Mask)
	}

	return &frame
}

func unmask(data []uint8, mask [4]byte) {
	for i, b := range data {
		maskIndex := i % len(mask)
		data[i] = b ^ mask[maskIndex]
	}
}

func (f Frame) String() string {
	out := ""
	out += "---------------------\n"
	out += fmt.Sprintf(" Is Fin: %v\n", f.header.Fin)
	out += fmt.Sprintf(" Op: %v\n", f.header.OpCode)
	out += fmt.Sprintf(" Is Masked: %v\n", f.header.IsMasked)
	out += fmt.Sprintf(" Mask: %x\n", f.header.Mask)
	out += fmt.Sprintf(" Length: %v\n", f.header.PayloadLength)
	out += fmt.Sprintf(" Data: %s\n", string(f.Data))
	out += "---------------------\n"
	return out
}

// func parseMask(container [2]uint8) [4]uint8 {

// }

func (f Frame) IsFragmented() bool {
	return f.header.Fin == false || f.IsEndFragment()
}

func (f Frame) IsStartFragment() bool {
	return f.header.Fin == false && f.header.OpCode != OpContinuationFrame
}

func (f Frame) IsContinuationFragment() bool {
	return f.header.Fin == false && f.header.OpCode == OpContinuationFrame
}

func (f Frame) IsEndFragment() bool {
	return f.header.Fin == true && f.header.OpCode == OpContinuationFrame
}

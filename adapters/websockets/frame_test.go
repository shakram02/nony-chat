package websockets

import (
	"testing"
)

func TestParseBit(t *testing.T) {
	cases := []struct {
		input  byte
		index  int
		output bool
	}{
		{
			input:  0b_0000_0001,
			index:  7,
			output: true,
		},
		{
			input:  0b_0111_1110,
			index:  7,
			output: false,
		},
		{
			input:  0b_0000_0010,
			index:  6,
			output: true,
		},
		{
			input:  0b_0000_0101,
			index:  6,
			output: false,
		},
		{
			input:  0b_1000_0000,
			index:  0,
			output: true,
		},
		{
			input:  0b_0111_1111,
			index:  0,
			output: false,
		},
	}

	for _, c := range cases {
		out := parseBit(c.input, c.index)
		if out != c.output {
			t.Errorf("expected byte[%d] to be [%v] found [%v]", c.index, c.output, out)
		}
	}
}

func TestParseHeaderOpCode(t *testing.T) {
	cases := []struct {
		input  byte
		output FrameOpCode
	}{
		{
			input:  0x0,
			output: OpContinuationFrame,
		},
		{
			input:  0x1,
			output: OpTextFrame,
		},
		{
			input:  0x2,
			output: OpBinaryFrame,
		},
		{
			input:  0x8,
			output: OpConnectionClose,
		},
		{
			input:  0x9,
			output: OpPing,
		},
		{
			input:  0xA,
			output: OpPong,
		},
	}

	for _, c := range cases {
		out := parseHeaderOpCode(c.input)

		if out != c.output {
			t.Errorf("expected to find [%v] found [%v]", c.output, out)
		}
	}
}

func TestParsePayloadLengthMode(t *testing.T) {
	cases := []struct {
		inputPayloadLengthByte byte
		outputLengthMode       payloadLengthMode
	}{
		{
			inputPayloadLengthByte: 0,
			outputLengthMode:       Simple,
		},
		{
			inputPayloadLengthByte: 125,
			outputLengthMode:       Simple,
		},
		{
			inputPayloadLengthByte: 126,
			outputLengthMode:       Extended16Bits,
		},
		{
			inputPayloadLengthByte: 0b_1111_1111, // 127
			outputLengthMode:       Extended64Bits,
		},
	}

	for _, c := range cases {
		out := parseHeaderPayloadLengthMode(c.inputPayloadLengthByte)
		if out != c.outputLengthMode {
			t.Errorf("Expected length mode to be [%v] found [%v]", c.outputLengthMode, out)
		}
	}
}

func TestParsePayloadLength(t *testing.T) {
	cases := []struct {
		description     string
		input           []uint8
		inputLengthMode payloadLengthMode
		outputLength    uint64
	}{
		{
			description:     "if 0-125, that is the payload length",
			input:           []byte{0},
			outputLength:    0,
			inputLengthMode: Simple,
		},
		{
			description:     "if 0-125, that is the payload length",
			input:           []byte{125},
			outputLength:    125,
			inputLengthMode: Simple,
		},
		{
			description:     "If 126, the following 2 bytes interpreted as a 16-bit unsigned integer are the payload length",
			input:           []byte{126, 0, 0xFF},
			outputLength:    0xFF,
			inputLengthMode: Extended16Bits,
		},
		{
			description:     "If 126, the following 2 bytes interpreted as a 16-bit unsigned integer are the payload length",
			input:           []byte{126, 0xFF, 0xFF},
			outputLength:    0xFFFF,
			inputLengthMode: Extended16Bits,
		},
		{
			description:     "If 127, the following 8 bytes interpreted as a 64-bit unsigned integer (the most significant bit MUST be 0)",
			input:           []byte{127, 0xEF, 0xCD, 0xAB, 0x89, 0x67, 0x45, 0x23, 0x01},
			outputLength:    0xEFCDAB8967452301,
			inputLengthMode: Extended64Bits,
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			parser := FrameParser{
				raw: c.input,
			}
			out := parser.parseHeaderPayloadLength(c.inputLengthMode)

			if out != c.outputLength {
				t.Errorf("Input [%v] Expected: [%X] found [%X]", c.input, c.outputLength, out)
			}
		})
	}
}

func TestParsePayloadMask(t *testing.T) {
	cases := []struct {
		description string
		input       []byte
		wantMasked  bool
		wantMask    [4]byte
	}{
		{
			description: "unmasked frame",
			input:       []byte{0x81, 0x00}, // FIN=1, OpCode=1, MASK=0
			wantMasked:  false,
			wantMask:    [4]byte{0, 0, 0, 0},
		},
		{
			description: "masked frame",
			input:       []byte{0x81, 0x80, 0xAA, 0xBB, 0xCC, 0xDD}, // FIN=1, OpCode=1, MASK=1, with mask
			wantMasked:  true,
			wantMask:    [4]byte{0xAA, 0xBB, 0xCC, 0xDD},
		},
		{
			description: "masked frame with payload",
			// FIN=1, OpCode=1, MASK=1, Length=5, with mask, then payload
			input:      []byte{0x81, 0x85, 0x11, 0x22, 0x33, 0x44, 0x01, 0x02, 0x03, 0x04, 0x05},
			wantMasked: true,
			wantMask:   [4]byte{0x11, 0x22, 0x33, 0x44},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			frame := New(c.input)

			if frame.header.IsMasked != c.wantMasked {
				t.Errorf("IsMasked = %v, want %v", frame.header.IsMasked, c.wantMasked)
			}

			if frame.header.IsMasked && frame.header.Mask != c.wantMask {
				t.Errorf("Mask = %v, want %v", frame.header.Mask, c.wantMask)
			}
		})
	}
}

func TestUnmask(t *testing.T) {
	cases := []struct {
		description string
		rawPayload  string
		mask        [4]byte
	}{
		{
			description: "unmasking a frame with a mask",
			rawPayload:  "Hello, World!",
			mask:        [4]byte{0xAA, 0xBB, 0xCC, 0xDD},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			maskedInput := make([]byte, len(c.rawPayload))
			for i, b := range c.rawPayload {
				maskedInput[i] = byte(b) ^ c.mask[i%len(c.mask)]
			}

			frame := Frame{
				raw: []byte(c.rawPayload),
				header: websocketHeader{
					IsMasked: true,
					Mask:     c.mask,
				},
				Data: maskedInput,
			}
			unmask(frame.Data, frame.header.Mask)

			if string(frame.Data) != c.rawPayload {
				t.Errorf("Expected unmasked input to be [%v] found [%v]", c.rawPayload, frame.Data)
			}

		})
	}
}

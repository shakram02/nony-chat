package websocket

import "testing"

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

func TestParsePayloadLength(t *testing.T) {
	cases := []struct {
		description string
		input       [9]byte
		output      uint64
	}{
		{
			description: "if 0-125, that is the payload length",
			input:       [9]byte{0},
			output:      0,
		},
		{
			description: "if 0-125, that is the payload length",
			input:       [9]byte{125},
			output:      125,
		},
		{
			description: "If 126, the following 2 bytes interpreted as a 16-bit unsigned integer are the payload length",
			input:       [9]byte{126, 0, 0xFF},
			output:      0xFF,
		},
		{
			description: "If 126, the following 2 bytes interpreted as a 16-bit unsigned integer are the payload length",
			input:       [9]byte{126, 0xFF, 0xFF},
			output:      0xFFFF,
		},
		{
			description: "If 127, the following 8 bytes interpreted as a 64-bit unsigned integer (the most significant bit MUST be 0)",
			input:       [9]byte{127, 0xEF, 0xCD, 0xAB, 0x89, 0x67, 0x45, 0x23, 0x01},
			output:      0xEFCDAB8967452301,
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			out := parseHeaderPayloadLength(c.input)

			if out != c.output {
				t.Errorf("Input [%v] Expected: [%X] found [%X]", c.input, c.output, out)
			}

		})
	}
}

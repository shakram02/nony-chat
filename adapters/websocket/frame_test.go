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
	// _ := []struct {
	// 	hasMask   bool
	// 	inputMask [4]byte
	// 	output    [4]byte
	// }{
	// 	{
	// 		hasMask:   false,
	// 		inputMask: [4]byte{},
	// 		output:    [4]byte{},
	// 	},
	// }
}

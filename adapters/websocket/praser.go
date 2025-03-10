package websocket

type FrameParser struct {
	raw     []byte
	pointer uint8
}

func newParser(raw []byte) *FrameParser {
	return &FrameParser{
		raw:     raw,
		pointer: 0,
	}
}

func (p *FrameParser) parseFrame() WebsocketFrame {
	header := p.parseHeader()

	frame := WebsocketFrame{
		header: header,
		Data:   p.raw[p.pointer:],
	}

	return frame
}

func (p *FrameParser) parseHeader() websocketHeader {
	fin := parseBit(p.getCurrentByte(), 0)
	opCode := parseHeaderOpCode(p.getCurrentByte())

	p.Advance(1)

	isMasked := parseBit(p.getCurrentByte(), 0)
	mode := parseHeaderPayloadLengthMode(p.getCurrentByte())
	length := p.parseHeaderPayloadLength(mode)
	// Payload length:  7 bits, 7+16 bits, or 7+64 bits
	switch mode {
	case Simple:
		p.Advance(1)
	case Extended16Bits:
		p.Advance(1 + 2)
	case Extended64Bits:
		p.Advance(1 + 8)
	}

	header := websocketHeader{
		Fin:           fin,
		OpCode:        opCode,
		IsMasked:      isMasked,
		PayloadLength: length,
	}

	// Will the mask bytes exist as 0000 if
	// the mask bit is unset? yes ->
	// frame-masking-key = 4( %x00-FF )
	//                     ; present only if frame-masked is 1
	//                     ; 32 bits in length
	if isMasked {
		mask := p.raw[p.pointer : p.pointer+4]
		header.Mask = [4]byte(mask)
		p.Advance(4)
	}

	return header
}

func parseHeaderOpCode(input uint8) FrameOpCode {
	return FrameOpCode(input & 0x0F)
}

func parseBit(input uint8, index int) bool {
	out := input & (1 << (7 - index))
	if out > 0 {
		return true
	}

	return false
}

func (p *FrameParser) getCurrentByte() uint8 {
	return p.raw[p.pointer]
}

func (p *FrameParser) Advance(by byte) {
	p.pointer += by
}

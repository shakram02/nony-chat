package nony

import (
	"encoding/json"
	"time"

	"github.com/shakram02/nony-chat/adapters/websockets"
)

type NonyPacketType string

const (
	NonyPacketTypeJoin    NonyPacketType = "join"
	NonyPacketTypeMessage NonyPacketType = "message"
)

type PacketContent struct {
	Text string `json:"text"`
}

type Packet struct {
	Type      NonyPacketType `json:"type"`
	UserId    string         `json:"userId"`
	RoomId    string         `json:"roomId"`
	Content   *PacketContent `json:"content"`
	Timestamp time.Time      `json:"timestamp"`
}

func New(websocketFrame *websockets.Frame) *Packet {
	data := websocketFrame.Data

	packet, err := parse(data)
	if err != nil {
		return nil
	}

	return packet
}

func parse(data []byte) (*Packet, error) {
	packet := &Packet{}

	err := json.Unmarshal(data, packet)
	if err != nil {
		return nil, err
	}

	return packet, nil
}

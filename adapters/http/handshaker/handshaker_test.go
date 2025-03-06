package handshaker

import (
	"testing"
)

func TestHandshakeResponse(t *testing.T) {

	connectionKey := "dGhlIHNhbXBsZSBub25jZQ=="
	response := makeHandshakeResponse(connectionKey)

	if response.WebsocketAccept != "s3pPLMBiTxaQ9kYGzzhZRbK+xOo=" {
		t.Errorf("Expected websocket accept to be s3pPLMBiTxaQ9kYGzzhZRbK+xOo=, got %s", response.WebsocketAccept)
	}
}

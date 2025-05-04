package handshaker

import (
	"crypto/sha1"
	"encoding/base64"
	"strings"

	http_parser "github.com/shakram02/nony-chat/adapters/http/parser"
)

// https://datatracker.ietf.org/doc/html/rfc6455#section-4.1
// If the response lacks a |Sec-WebSocket-Accept| header field or
// the |Sec-WebSocket-Accept| contains a value other than the
// base64-encoded SHA-1 of the concatenation of the |Sec-WebSocket-
// Key| (as a string, not base64-decoded) with the string "258EAFA5-
// E914-47DA-95CA-C5AB0DC85B11" but ignoring any leading and
// trailing whitespace, the client MUST _Fail the WebSocket
// Connection_.
const rfc6455ServerResponseGuid = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
const lineSep = "\r\n"

type handshakeResponse struct {
	ClientUuid      string
	WebsocketAccept string
}

type HandshakedClient struct {
	RemoteAddr       string
	SocketIdentifier string
}

func MakeAcceptanceResposne(clientHandshake http_parser.WebsocketHandshake) []byte {
	websocketAccept := makeHandshakeAcceptHeaderValue(clientHandshake.Headers.SecWebSocketKey)
	return makeResponse(websocketAccept)
}

func makeResponse(resp handshakeResponse) []byte {
	responseString := ""
	responseString += "HTTP/1.1 101 Switching Protocols" + lineSep
	responseString += "Upgrade: websocket" + lineSep
	responseString += "Connection: Upgrade" + lineSep
	responseString += "Sec-WebSocket-Accept: " + resp.WebsocketAccept + lineSep
	responseString += lineSep

	return []byte(responseString)
}

func makeHandshakeAcceptHeaderValue(websocketKey string) handshakeResponse {
	trimmed := strings.TrimSpace(websocketKey)
	handshakeAccept := trimmed + rfc6455ServerResponseGuid
	hasher := sha1.New()
	hasher.Write([]byte(handshakeAccept))
	data := hasher.Sum(nil)

	dst := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(dst, data)

	return handshakeResponse{
		ClientUuid:      websocketKey,
		WebsocketAccept: string(dst),
	}
}

func MakeRejectionResponse() []byte {
	responseString := "HTTP/1.1 400 Bad Request\r\n\r\n"
	return []byte(responseString)
}

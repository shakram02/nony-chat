package http_parser

import (
	"fmt"
	"net/http"
	"strings"
)

type HandshakeRequestLine struct {
	Uri string
}

type HandshakeHeaders struct {
	//The handshake MUST be a valid HTTP request as specified by [RFC2616].
	// https://datatracker.ietf.org/doc/html/rfc2616#page-128
	// The Host field value MUST represent
	// the naming authority of the origin server or gateway given by the
	// original URL. This allows the origin server or gateway to
	// differentiate between internally-ambiguous URLs, such as the root "/"
	// URL of a server for multiple host names on a single IP address.
	Host string
	// An |Upgrade| header field containing the value "websocket",
	// treated as an ASCII case-insensitive value.
	Upgrade string
	//  The request MUST contain a |Connection| header field whose value
	// MUST include the "Upgrade" token.
	Connection string
	// The request MUST include a header field with the name
	// |Sec-WebSocket-Key|.  The value of this header field MUST be a
	// nonce consisting of a randomly selected 16-byte value that has
	// been base64-encoded (see Section 4 of [RFC4648]).  The nonce
	// MUST be selected randomly for each connection.
	SecWebSocketKey string
}

type WebsocketHandshake struct {
	RequestLine HandshakeRequestLine
	Headers     HandshakeHeaders
}

// https://datatracker.ietf.org/doc/html/rfc6455#section-4.1
var requiredHeaders = map[string]string{
	"Host":              "",
	"Upgrade":           "websocket",
	"Connection":        "Upgrade",
	"Sec-WebSocket-Key": "",
	// The request MUST include a header field with the name
	// |Sec-WebSocket-Version|.  The value of this header field MUST be
	// 13.
	"Sec-WebSocket-Version": "13",
}

func ParseUpgradeRequest(request []byte) (WebsocketHandshake, error) {
	requestString := string(request)

	httpRequestParts := strings.Split(requestString, "\r\n")
	if len(httpRequestParts) == 0 {
		return WebsocketHandshake{}, fmt.Errorf("Invalid HTTP request: %s", requestString)
	}

	requestLine, err := parseHandshakeRequestLine(httpRequestParts[0])
	if err != nil {
		return WebsocketHandshake{}, fmt.Errorf("Invalid request line: %s", httpRequestParts[0])
	}

	headerParts := httpRequestParts[1:]
	headers, err := parseHandshakeHeaders(headerParts)
	if err != nil {
		return WebsocketHandshake{}, fmt.Errorf("Invalid headers: %s", err)
	}

	return WebsocketHandshake{
		RequestLine: requestLine,
		Headers:     headers,
	}, nil
}

func parseHandshakeRequestLine(requestLine string) (HandshakeRequestLine, error) {
	parts := strings.Split(strings.TrimSpace(requestLine), " ")
	if len(parts) < 3 {
		return HandshakeRequestLine{}, fmt.Errorf("Invalid Handshake Request-Line")
	}

	parts = parts[:3] // The handshake has just 3 parts e.g. GET /chat HTTP/1.1
	if parts[0] != "GET" {
		return HandshakeRequestLine{}, fmt.Errorf("Invalid method")
	}

	if parts[1][0] != '/' {
		return HandshakeRequestLine{}, fmt.Errorf("Invalid URI")
	}

	major, minor, ok := http.ParseHTTPVersion(parts[2])
	if !ok {
		return HandshakeRequestLine{}, fmt.Errorf("Invalid Protocol version")
	}

	acceptedVersion := (major == 1 && minor == 1) || (major > 1)
	if !acceptedVersion {
		return HandshakeRequestLine{}, fmt.Errorf("Invalid Protocol version")
	}

	return HandshakeRequestLine{Uri: parts[1]}, nil
}

func parseHandshakeHeaders(headerLines []string) (HandshakeHeaders, error) {
	headers := parseHttpHeaders(headerLines)
	if !hasRequiredHanshakeHeaders(headers) {
		return HandshakeHeaders{}, fmt.Errorf("Invalid headers")
	}

	return HandshakeHeaders{
		Host:            headers["Host"],
		Upgrade:         headers["Upgrade"],
		Connection:      headers["Connection"],
		SecWebSocketKey: headers["Sec-WebSocket-Key"],
	}, nil

}

func hasRequiredHanshakeHeaders(headers map[string]string) bool {
	for k, v := range requiredHeaders {
		value, ok := headers[k]

		if !ok {
			return false
		}

		if v == "" {
			// Header value isn't required for validation
			continue
		}

		if k == "Connection" && strings.Contains(value, v) {
			continue
		}

		if v != value {
			return false
		}
	}

	return true
}

func parseHttpHeaders(headerLines []string) map[string]string {
	headers := make(map[string]string)
	for _, line := range headerLines {
		if strings.TrimSpace(line) == "" {
			// Body separator
			break
		}
		splits := strings.Split(line, ": ")

		if len(splits) != 2 {
			// Skip invalid headers
			continue
		}

		key := splits[0]
		value := splits[1]
		headers[key] = value
	}

	return headers
}

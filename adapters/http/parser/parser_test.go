package http_parser

import "testing"

func TestParseRequestLine(t *testing.T) {
	cases := []struct {
		description string
		input       string
		fails       bool
		output      HandshakeRequestLine
	}{
		{
			description: "Valid GET request line",
			input:       "GET / HTTP/1.1\r\n",
			output:      HandshakeRequestLine{Uri: "/"},
		},
		{
			description: "Valid GET request line with URI",
			input:       "GET /chat HTTP/1.1\r\n",
			output:      HandshakeRequestLine{Uri: "/chat"},
		},
		{
			description: "Invalid POST request line",
			input:       "POST /chat HTTP/1.1\r\n",
			fails:       true,
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			actual, err := parseHandshakeRequestLine(c.input)
			if c.fails {
				if err == nil {
					t.Errorf("Expected error for input: %s", c.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for input: %s", c.input)
				}
				if actual.Uri != c.output.Uri {
					t.Errorf("Expected URI: %s, got: %s", c.output.Uri, actual.Uri)
				}
			}
		})
	}
}

func TestParseHasValidHandshakeHeaders(t *testing.T) {
	cases := []struct {
		description string
		input       map[string]string
		valid       bool
	}{
		{
			description: "Valid handshake headers",
			input: map[string]string{
				"Host":                  "astro",
				"Upgrade":               "websocket",
				"Sec-WebSocket-Key":     "kBQW2M+CkClJ1bvTT8O4LA==",
				"Connection":            "Upgrade",
				"Sec-WebSocket-Version": "13",
			},
			valid: true,
		},
		{
			description: "Valid handshake headers with keep-alive",
			input: map[string]string{
				"Host":                  "astro",
				"Upgrade":               "websocket",
				"Sec-WebSocket-Key":     "kBQW2M+CkClJ1bvTT8O4LA==",
				"Connection":            "keep-alive, Upgrade", // Upgrade, must be present. Not necessarily the full value
				"Sec-WebSocket-Version": "13",
			},
			valid: true,
		},
		{
			description: "Invalid handshake headers missing Host",
			input: map[string]string{
				"Upgrade":               "websocket",
				"Sec-WebSocket-Key":     "kBQW2M+CkClJ1bvTT8O4LA==",
				"Connection":            "Upgrade",
				"Sec-WebSocket-Version": "13",
			},
			valid: false,
		},

		{
			description: "Invalid handshake headers missing Upgrade",
			input: map[string]string{
				"Host":                  "astro",
				"Sec-WebSocket-Key":     "kBQW2M+CkClJ1bvTT8O4LA==",
				"Connection":            "Upgrade",
				"Sec-WebSocket-Version": "13",
			},
			valid: false,
		},
		{
			description: "Invalid handshake headers missing Sec-WebSocket-Key",
			input: map[string]string{
				"Host":                  "astro",
				"Upgrade":               "websocket",
				"Connection":            "Upgrade",
				"Sec-WebSocket-Version": "13",
			},
			valid: false,
		},
		{
			description: "Invalid handshake headers missing Sec-WebSocket-Version",
			input: map[string]string{
				"Host":                  "astro",
				"Upgrade":               "websocket",
				"Sec-WebSocket-Key":     "kBQW2M+CkClJ1bvTT8O4LA==",
				"Sec-WebSocket-Version": "13",
			},
			valid: false,
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			actual := hasRequiredHanshakeHeaders(c.input)
			if actual != c.valid {
				t.Errorf("Expected valid: %t, got: %t", c.valid, actual)
			}
		})
	}
}

func TestParseHandshakeHeaders(t *testing.T) {
	cases := []struct {
		input  []string
		fails  bool
		output HandshakeHeaders
	}{
		{
			input: []string{
				"Host: astro",
				"Upgrade: websocket",
				"Sec-WebSocket-Key: kBQW2M+CkClJ1bvTT8O4LA==",
				"Connection: Upgrade",
				"Sec-WebSocket-Version: 13",
			},
			output: HandshakeHeaders{
				Host:            "astro",
				Upgrade:         "websocket",
				SecWebSocketKey: "kBQW2M+CkClJ1bvTT8O4LA==",
				Connection:      "Upgrade",
			},
		},
		{
			input: []string{
				"Host: astro",
				"Upgrade: websocket",
				"Sec-WebSocket-Key: kBQW2M+CkClJ1bvTT8O4LA==",
				"Connection: keep-alive, Upgrade", // Upgrade, must be present. Not necessarily the full value
				"Sec-WebSocket-Version: 13",
			},
			output: HandshakeHeaders{
				Host:            "astro",
				Upgrade:         "websocket",
				SecWebSocketKey: "kBQW2M+CkClJ1bvTT8O4LA==",
				Connection:      "keep-alive, Upgrade",
			},
		},
		{
			input: []string{
				"Host: astro",
				"Upgrade: websocket",
			},
			fails: true,
		},
	}

	for _, c := range cases {
		actual, err := parseHandshakeHeaders(c.input)
		if c.fails {
			if err == nil {
				t.Errorf("Expected error for input: %s", c.input)
			}
		} else {

			if err != nil {
				t.Errorf("Unexpected error for input: %s", c.input)
			}

			if actual.Host != c.output.Host {
				t.Errorf("Expected Host: %s, got: %s", c.output.Host, actual.Host)
			}

			if actual.Upgrade != c.output.Upgrade {
				t.Errorf("Expected Upgrade: %s, got: %s", c.output.Upgrade, actual.Upgrade)
			}

			if actual.SecWebSocketKey != c.output.SecWebSocketKey {
				t.Errorf("Expected Sec-WebSocket-Key: %s, got: %s", c.output.SecWebSocketKey, actual.SecWebSocketKey)
			}

			if actual.Connection != c.output.Connection {
				t.Errorf("Expected Connection: %s, got: %s", c.output.Connection, actual.Connection)
			}
		}
	}
}

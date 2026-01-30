package signature

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// ToPrefixedMessage prepends the Ethereum Signed Message prefix to a message.
// This follows the standard format: "\x19Ethereum Signed Message:\n" + len(message) + message
//
// Example:
//
//	result := ToPrefixedMessage(NewSignableMessage("hello world"))
//	// Returns the prefixed message as hex
func ToPrefixedMessage(message SignableMessage) string {
	var messageHex string

	// Determine the message content
	if message.Raw != nil {
		switch v := message.Raw.(type) {
		case string:
			// Raw hex string
			if isHex(v) {
				messageHex = v
			} else {
				messageHex = stringToHex(v)
			}
		case []byte:
			// Raw bytes
			messageHex = bytesToHex(v)
		default:
			messageHex = "0x"
		}
	} else {
		// Plain string message
		messageHex = stringToHex(message.Message)
	}

	// Calculate the size in bytes
	size := sizeHex(messageHex)

	// Create the prefix
	prefix := stringToHex(fmt.Sprintf("%s%d", PresignMessagePrefix, size))

	// Concatenate prefix and message
	return concatHex(prefix, messageHex)
}

// ToPrefixedMessageBytes returns the prefixed message as bytes.
func ToPrefixedMessageBytes(message SignableMessage) []byte {
	hexStr := ToPrefixedMessage(message)
	return hexToBytes(hexStr)
}

// Helper functions

func isHex(s string) bool {
	return strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X")
}

func stringToHex(s string) string {
	return "0x" + hex.EncodeToString([]byte(s))
}

func bytesToHex(b []byte) string {
	return "0x" + hex.EncodeToString(b)
}

func hexToBytes(s string) []byte {
	s = strings.TrimPrefix(s, "0x")
	s = strings.TrimPrefix(s, "0X")
	if len(s)%2 != 0 {
		s = "0" + s
	}
	b, _ := hex.DecodeString(s)
	return b
}

func sizeHex(h string) int {
	h = strings.TrimPrefix(h, "0x")
	h = strings.TrimPrefix(h, "0X")
	return (len(h) + 1) / 2
}

func concatHex(values ...string) string {
	var builder strings.Builder
	builder.WriteString("0x")
	for _, h := range values {
		h = strings.TrimPrefix(h, "0x")
		h = strings.TrimPrefix(h, "0X")
		builder.WriteString(h)
	}
	return builder.String()
}

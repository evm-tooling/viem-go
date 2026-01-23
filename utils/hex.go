package utils

import "encoding/hex"

// BytesToHex converts a byte slice to a hex string prefixed with "0x".
func BytesToHex(b []byte) string {
	return "0x" + hex.EncodeToString(b)
}

// BytesToHexUnprefixed converts a byte slice to a hex string without the "0x" prefix.
func BytesToHexUnprefixed(b []byte) string {
	return hex.EncodeToString(b)
}

// HexToBytes converts a hex string (with or without "0x" prefix) to a byte slice.
func HexToBytes(s string) ([]byte, error) {
	if len(s) >= 2 && s[0:2] == "0x" {
		s = s[2:]
	}
	return hex.DecodeString(s)
}

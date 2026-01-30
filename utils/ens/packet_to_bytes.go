package ens

import (
	"strings"
)

// PacketToBytes encodes an ENS name into a DNS packet ByteArray.
// This is used for the off-chain resolution protocol (CCIP-Read).
//
// The format follows DNS encoding:
// - Each label is prefixed with its length byte
// - The packet ends with a zero byte
//
// Example:
//
//	bytes := PacketToBytes("vitalik.eth")
//	// []byte{7, 'v', 'i', 't', 'a', 'l', 'i', 'k', 3, 'e', 't', 'h', 0}
//
//	hex := PacketToBytesHex("vitalik.eth")
//	// "0x077669746c696b03657468000"
//
// @see https://docs.ens.domains/resolution/names#dns
func PacketToBytes(packet string) []byte {
	// Strip leading and trailing dots
	value := strings.Trim(packet, ".")

	if value == "" {
		return []byte{0}
	}

	// Calculate total size needed
	labels := strings.Split(value, ".")
	totalLen := 1 // For final zero byte
	for _, label := range labels {
		encoded := []byte(label)
		// If the length is > 255, we'll encode as labelhash
		if len(encoded) > 255 {
			encoded = []byte(EncodeLabelhash(Labelhash(label)))
		}
		totalLen += 1 + len(encoded) // length byte + content
	}

	// Build the packet
	bytes := make([]byte, totalLen)
	offset := 0

	for _, label := range labels {
		encoded := []byte(label)
		// If the length is > 255, make the encoded label value a labelhash
		// This is compatible with the universal resolver
		if len(encoded) > 255 {
			encoded = []byte(EncodeLabelhash(Labelhash(label)))
		}

		bytes[offset] = byte(len(encoded))
		copy(bytes[offset+1:], encoded)
		offset += len(encoded) + 1
	}

	// Final zero byte is already there (from make)
	return bytes
}

// PacketToBytesHex returns the DNS packet as a hex string.
func PacketToBytesHex(packet string) string {
	return bytesToHex(PacketToBytes(packet))
}

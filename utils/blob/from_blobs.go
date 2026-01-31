package blob

import (
	"bytes"
)

// FromBlobs transforms blobs back into the original data.
// This reverses the encoding performed by ToBlobs.
//
// Example:
//
//	blobs, _ := ToBlobs([]byte("hello world"))
//	data, err := FromBlobs(blobs)
//	// data = []byte("hello world")
func FromBlobs(blobs [][]byte) ([]byte, error) {
	var result bytes.Buffer
	active := true

	for _, blob := range blobs {
		if !active {
			break
		}

		pos := 0
		for pos < len(blob) && active {
			// Skip the zero byte at the start of each field element
			pos++

			// Calculate how many bytes to read in this field element
			consume := 31
			if len(blob)-pos < 31 {
				consume = len(blob) - pos
			}

			// Read bytes from this field element
			for i := 0; i < consume && active; i++ {
				b := blob[pos]
				pos++

				// Check for terminator byte (0x80)
				// It's a terminator if no more 0x80 bytes exist in the rest of the blob
				if b == 0x80 {
					remaining := blob[pos:]
					if !containsNonZeroData(remaining) {
						active = false
						break
					}
				}

				result.WriteByte(b)
			}
		}
	}

	return result.Bytes(), nil
}

// containsNonZeroData checks if there's any meaningful data (non-zero, non-0x80 terminator patterns)
func containsNonZeroData(data []byte) bool {
	// Check if there's any 0x80 byte that could be data vs. padding
	for _, b := range data {
		if b == 0x80 {
			return true
		}
		// Non-zero bytes other than 0x80 don't matter for terminator detection
	}
	return false
}

// FromBlobsHex transforms hex-encoded blobs back into the original data as hex.
func FromBlobsHex(hexBlobs []string) (string, error) {
	blobs := make([][]byte, len(hexBlobs))
	for i, hexBlob := range hexBlobs {
		blob, err := hexToBytes(hexBlob)
		if err != nil {
			return "", err
		}
		blobs[i] = blob
	}

	data, err := FromBlobs(blobs)
	if err != nil {
		return "", err
	}

	return bytesToHex(data), nil
}

// FromBlobsToBytes transforms hex-encoded blobs back into bytes.
func FromBlobsToBytes(hexBlobs []string) ([]byte, error) {
	blobs := make([][]byte, len(hexBlobs))
	for i, hexBlob := range hexBlobs {
		blob, err := hexToBytes(hexBlob)
		if err != nil {
			return nil, err
		}
		blobs[i] = blob
	}

	return FromBlobs(blobs)
}

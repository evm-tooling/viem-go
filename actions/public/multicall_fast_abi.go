package public

import (
	"encoding/binary"
	"fmt"
)

// Hand-rolled ABI encoder/decoder for the multicall3 aggregate3 function.
//
// This bypasses go-ethereum's reflect-heavy accounts/abi package entirely,
// encoding and decoding directly to/from bytes. For a 100-call multicall,
// this eliminates ~1000 reflect operations and ~400 big.Int allocations
// that the generic go-ethereum packer performs.
//
// The ABI layouts are fixed and known at compile time:
//   Encode: aggregate3(tuple(address,bool,bytes)[])
//   Decode: returns tuple(bool,bytes)[]

// pad32 rounds n up to the nearest multiple of 32.
func pad32(n int) int {
	return (n + 31) &^ 31
}

// writeUint256 writes a uint64 as a big-endian uint256 at buf[off:off+32].
// Caller must ensure buf is zero-initialized (from make) for bytes [off:off+24].
func writeUint256(buf []byte, off int, v uint64) {
	binary.BigEndian.PutUint64(buf[off+24:off+32], v)
}

// readUint256AsInt reads a big-endian uint256 at buf[off:off+32] as an int.
// Only reads the low 8 bytes since aggregate3 array sizes/offsets fit in uint64.
func readUint256AsInt(buf []byte, off int) int {
	return int(binary.BigEndian.Uint64(buf[off+24 : off+32]))
}

// encodeAggregate3Fast encodes Call3 structs directly to ABI-encoded bytes.
// Zero reflection, zero big.Int allocations, single buffer allocation.
//
// ABI layout for tuple(address target, bool allowFailure, bytes callData)[]:
//
//   [offset to array = 32]                          (32 bytes)
//   [array length = N]                              (32 bytes)
//   [offset to tuple[0], tuple[1], ... tuple[N-1]]  (N * 32 bytes)
//   [tuple[0] data]                                 (variable)
//   [tuple[1] data]                                 (variable)
//   ...
//
// Each tuple (address, bool, bytes):
//   [address left-padded to 32]   (32 bytes)
//   [allowFailure as uint256]     (32 bytes)
//   [offset to bytes = 96]        (32 bytes)  -- always 3*32, points past head
//   [callData length]             (32 bytes)
//   [callData right-padded to 32] (ceil32 bytes)
func encodeAggregate3Fast(calls []Call3) []byte {
	n := len(calls)
	if n == 0 {
		buf := make([]byte, 64)
		writeUint256(buf, 0, 32) // offset
		return buf               // length = 0 is zero-initialized
	}

	// Pre-calculate sizes for a single allocation
	tupleSizes := make([]int, n)
	totalTupleData := 0
	for i, c := range calls {
		// head(3*32) + callData length(32) + padded callData
		sz := 128 + pad32(len(c.CallData))
		tupleSizes[i] = sz
		totalTupleData += sz
	}

	// Total: offset(32) + length(32) + N offsets(N*32) + tuple data
	total := 64 + n*32 + totalTupleData
	buf := make([]byte, total)

	// Outer offset to array data
	writeUint256(buf, 0, 32)
	// Array length
	writeUint256(buf, 32, uint64(n))

	// Write tuple offsets (relative to start of offsets area at byte 64)
	tupleOffset := n * 32 // first tuple starts after N offset words
	for i := range calls {
		writeUint256(buf, 64+i*32, uint64(tupleOffset))
		tupleOffset += tupleSizes[i]
	}

	// Write tuple data
	pos := 64 + n*32
	for _, c := range calls {
		// address (left-padded: 12 zero bytes + 20 address bytes)
		copy(buf[pos+12:pos+32], c.Target[:])
		pos += 32

		// allowFailure (bool as uint256)
		if c.AllowFailure {
			buf[pos+31] = 1
		}
		pos += 32

		// offset to callData = 96 (always 3 head words * 32)
		writeUint256(buf, pos, 96)
		pos += 32

		// callData length
		writeUint256(buf, pos, uint64(len(c.CallData)))
		pos += 32

		// callData (right-padded, padding is zero from make)
		if len(c.CallData) > 0 {
			copy(buf[pos:], c.CallData)
			pos += pad32(len(c.CallData))
		}
	}

	return buf
}

// decodeAggregate3Fast decodes aggregate3 return data directly from ABI bytes.
// Zero reflection, zero big.Int allocations, direct byte slicing.
//
// ABI layout for tuple(bool success, bytes returnData)[]:
//
//   [offset to array = 32]                          (32 bytes)
//   [array length = N]                              (32 bytes)
//   [offset to tuple[0], tuple[1], ... tuple[N-1]]  (N * 32 bytes)
//   [tuple[0] data]                                 (variable)
//   ...
//
// Each result tuple (bool, bytes):
//   [success as uint256]            (32 bytes)
//   [offset to bytes = 64]          (32 bytes)  -- always 2*32
//   [returnData length]             (32 bytes)
//   [returnData right-padded to 32] (ceil32 bytes)
func decodeAggregate3Fast(data []byte) ([]aggregate3Result, error) {
	if len(data) < 64 {
		return nil, fmt.Errorf("aggregate3 result too short: %d bytes", len(data))
	}

	// Read outer offset (should be 32)
	offset := readUint256AsInt(data, 0)
	if offset < 0 || offset+32 > len(data) {
		return nil, fmt.Errorf("aggregate3: invalid array offset %d (data len %d)", offset, len(data))
	}

	// Read array length
	n := readUint256AsInt(data, offset)
	if n < 0 || n > 1000000 {
		return nil, fmt.Errorf("aggregate3: invalid array length %d", n)
	}
	if n == 0 {
		return []aggregate3Result{}, nil
	}

	// Offsets area starts right after the length word
	offsetsStart := offset + 32
	if offsetsStart+n*32 > len(data) {
		return nil, fmt.Errorf("aggregate3: data too short for %d tuple offsets", n)
	}

	results := make([]aggregate3Result, n)
	for i := 0; i < n; i++ {
		// Read offset to this tuple (relative to offsetsStart)
		tupleRel := readUint256AsInt(data, offsetsStart+i*32)
		tupleStart := offsetsStart + tupleRel

		if tupleStart+64 > len(data) {
			return nil, fmt.Errorf("aggregate3: tuple %d start out of bounds (offset %d, data len %d)", i, tupleStart, len(data))
		}

		// success = bool at tupleStart (check byte 31 of the uint256 word)
		results[i].Success = data[tupleStart+31] == 1

		// Read offset to returnData (relative to tupleStart, should be 64 for (bool,bytes))
		rdOffset := readUint256AsInt(data, tupleStart+32)
		rdStart := tupleStart + rdOffset

		if rdStart+32 > len(data) {
			return nil, fmt.Errorf("aggregate3: returnData offset out of bounds for tuple %d", i)
		}

		// returnData length
		rdLen := readUint256AsInt(data, rdStart)
		if rdLen < 0 || rdStart+32+rdLen > len(data) {
			return nil, fmt.Errorf("aggregate3: returnData out of bounds for tuple %d (len %d, available %d)", i, rdLen, len(data)-rdStart-32)
		}

		// Extract returnData (copy to own slice to avoid holding the entire response buffer)
		if rdLen > 0 {
			results[i].ReturnData = make([]byte, rdLen)
			copy(results[i].ReturnData, data[rdStart+32:rdStart+32+rdLen])
		}
	}

	return results, nil
}

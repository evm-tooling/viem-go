package data

// IsBytes checks if a value is a byte slice.
// In Go, this is primarily useful for type assertions in generic contexts.
//
// Example:
//
//	IsBytes([]byte{0x01, 0x02}) // true
//	IsBytes("0x0102")          // false
//	IsBytes(nil)               // false
func IsBytes(value any) bool {
	if value == nil {
		return false
	}
	_, ok := value.([]byte)
	return ok
}

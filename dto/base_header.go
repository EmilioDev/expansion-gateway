package dto

// the base header of all the headers
type BaseHeader struct {
	// Example of possible future fields:
	// Flags byte
	// Timestamp int64
	// Nonce []byte
}

func (h BaseHeader) IsValid() bool {
	// Later you might implement validation rules
	return true
}

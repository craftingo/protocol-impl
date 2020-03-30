package protocol_impl

import "math"

type (
	Boolean       bool
	Byte          int8
	UnsignedByte  uint8
	Short         int16
	UnsignedShort uint16
	Int           int32
	Long          int64
	Float         float32
	Double        float64
	String        string
	Chat          = String
	Identifier    = String
	VarInt        int32
	VarLong       int64
	Position      struct {
		X, Y, Z int
	}
	Angle int8
	// TODO: Add UUID data type
	// TODO: Add NBT data type
	ByteArray []byte
)

// Encode a wrapped boolean
func (value Boolean) Encode() []byte {
	if value {
		return []byte{0x01}
	}
	return []byte{0x00}
}

// TODO: Decode wrapped booleans

// Encode a wrapped byte
func (value Byte) Encode() []byte {
	return []byte{byte(value)}
}

// TODO: Decode wrapped bytes

// Encode a wrapped unsigned byte
func (value UnsignedByte) Encode() []byte {
	return []byte{byte(value)}
}

// TODO: Decode wrapped unsigned bytes

// Encode a wrapped short
func (value Short) Encode() []byte {
	number := uint16(value)
	return []byte{
		byte(number >> 8),
		byte(number),
	}
}

// TODO: Decode wrapped shorts

// Encode a wrapped unsigned short
func (value UnsignedShort) Encode() []byte {
	number := uint16(value)
	return []byte{
		byte(number >> 8),
		byte(number),
	}
}

// TODO: Decode wrapped unsigned shorts

// Encode a wrapped integer
func (value Int) Encode() []byte {
	number := uint32(value)
	return []byte{
		byte(number >> 24),
		byte(number >> 16),
		byte(number >> 8),
		byte(number),
	}
}

// TODO: Decode wrapped integers

// Encode a wrapped long
func (value Long) Encode() []byte {
	number := uint64(value)
	return []byte{
		byte(number >> 56),
		byte(number >> 48),
		byte(number >> 40),
		byte(number >> 32),
		byte(number >> 24),
		byte(number >> 16),
		byte(number >> 8),
		byte(number),
	}
}

// TODO: Decode wrapped longs

// Encode a wrapped float
func (value Float) Encode() []byte {
	return Int(math.Float32bits(float32(value))).Encode()
}

// TODO: Decode wrapped floats

// Encode a wrapped double
func (value Double) Encode() []byte {
	return Long(math.Float64bits(float64(value))).Encode()
}

// TODO: Decode wrapped doubles

// Encode a wrapped string
func (value String) Encode() (encoded []byte) {
	byteArray := []byte(value)
	encoded = append(encoded, VarInt(len(byteArray)).Encode()...)
	encoded = append(encoded, byteArray...)
	return
}

// TODO: Decode wrapped strings

// Encode a VarInt
func (value VarInt) Encode() (encoded []byte) {
	number := uint32(value)
	for {
		temp := number & 0x7F
		number >>= 7
		if number != 0 {
			temp |= 0x80
		}
		encoded = append(encoded, byte(temp))
		if number == 0 {
			break
		}
	}
	return
}

// TODO: Decode VarInts

// Encode a VarLong
func (value VarLong) Encode() (encoded []byte) {
	number := uint64(value)
	for {
		temp := number & 0x7F
		number >>= 7
		if number != 0 {
			temp |= 0x80
		}
		encoded = append(encoded, byte(temp))
		if number == 0 {
			break
		}
	}
	return
}

// TODO: Decode VarLongs

// Encode a position
func (value Position) Encode() []byte {
	encoded := make([]byte, 8)
	position := uint64(value.X&0x3FFFFFF)<<38 | uint64((value.Z&0x3FFFFFF)>>12) | uint64(value.Y&0xFFF)
	for i := 7; i >= 0; i-- {
		encoded[i] = byte(position)
		position >>= 8
	}
	return encoded
}

// TODO: Decode positions

// Encode a wrapped byte array
func (value ByteArray) Encode() []byte {
	return append(VarInt(len(value)).Encode(), value...)
}

// TODO: Decode wrapped byte arrays

package protocol_impl

import (
	"errors"
	"io"
	"math"
)

// FieldReader represents an object that can read packet contents
type FieldReader interface {
	io.ByteReader
	io.Reader
}

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

// ReadNBytes reads n bytes from the given reader
func ReadNBytes(reader FieldReader, n int) (readBytes []byte, err error) {
	readBytes = make([]byte, n)
	for i := 0; i < n; i++ {
		readBytes[i], err = reader.ReadByte()
		if err != nil {
			return
		}
	}
	return
}

// Encode a wrapped boolean
func (value Boolean) Encode() []byte {
	if value {
		return []byte{0x01}
	}
	return []byte{0x00}
}

// Decode a wrapped boolean
func (value *Boolean) Decode(reader FieldReader) error {
	received, err := reader.ReadByte()
	if err != nil {
		return err
	}

	*value = received != 0
	return nil
}

// Encode a wrapped byte
func (value Byte) Encode() []byte {
	return []byte{byte(value)}
}

// Decode a wrapped byte
func (value *Byte) Decode(reader FieldReader) error {
	received, err := reader.ReadByte()
	if err != nil {
		return err
	}

	*value = Byte(received)
	return nil
}

// Encode a wrapped unsigned byte
func (value UnsignedByte) Encode() []byte {
	return []byte{byte(value)}
}

// Decode a wrapped unsigned byte
func (value *UnsignedByte) Decode(reader FieldReader) error {
	received, err := reader.ReadByte()
	if err != nil {
		return err
	}

	*value = UnsignedByte(received)
	return nil
}

// Encode a wrapped short
func (value Short) Encode() []byte {
	number := uint16(value)
	return []byte{
		byte(number >> 8),
		byte(number),
	}
}

// Decode a wrapped short
func (value *Short) Decode(reader FieldReader) error {
	readBytes, err := ReadNBytes(reader, 2)
	if err != nil {
		return err
	}

	*value = Short(int16(readBytes[0])<<8 | int16(readBytes[1]))
	return nil
}

// Encode a wrapped unsigned short
func (value UnsignedShort) Encode() []byte {
	number := uint16(value)
	return []byte{
		byte(number >> 8),
		byte(number),
	}
}

// Decode a wrapped unsigned short
func (value *UnsignedShort) Decode(reader FieldReader) error {
	readBytes, err := ReadNBytes(reader, 2)
	if err != nil {
		return err
	}

	*value = UnsignedShort(int16(readBytes[0])<<8 | int16(readBytes[1]))
	return nil
}

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

// Decode a wrapped integer
func (value *Int) Decode(reader FieldReader) error {
	readBytes, err := ReadNBytes(reader, 4)
	if err != nil {
		return err
	}

	*value = Int(int32(readBytes[0])<<24 | int32(readBytes[1])<<16 | int32(readBytes[2])<<8 | int32(readBytes[3]))
	return nil
}

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

// Decode a wrapped long
func (value *Long) Decode(reader FieldReader) error {
	readBytes, err := ReadNBytes(reader, 8)
	if err != nil {
		return err
	}

	*value = Long(int64(readBytes[0])<<56 | int64(readBytes[1])<<48 | int64(readBytes[2])<<40 |
		int64(readBytes[3])<<32 | int64(readBytes[4])<<24 | int64(readBytes[5])<<16 |
		int64(readBytes[6])<<8 | int64(readBytes[7]))
	return nil
}

// Encode a wrapped float
func (value Float) Encode() []byte {
	return Int(math.Float32bits(float32(value))).Encode()
}

// Decode a wrapped float
func (value *Float) Decode(reader FieldReader) error {
	var intValue Int
	err := intValue.Decode(reader)
	if err != nil {
		return err
	}

	*value = Float(math.Float32frombits(uint32(intValue)))
	return nil
}

// Encode a wrapped double
func (value Double) Encode() []byte {
	return Long(math.Float64bits(float64(value))).Encode()
}

// Decode a wrapped double
func (value *Double) Decode(reader FieldReader) error {
	var longValue Long
	err := longValue.Decode(reader)
	if err != nil {
		return err
	}

	*value = Double(math.Float64frombits(uint64(longValue)))
	return nil
}

// Encode a wrapped string
func (value String) Encode() (encoded []byte) {
	byteArray := []byte(value)
	encoded = append(encoded, VarInt(len(byteArray)).Encode()...)
	encoded = append(encoded, byteArray...)
	return
}

// Decode a wrapped string
func (value *String) Decode(reader FieldReader) error {
	var stringLength VarInt
	err := stringLength.Decode(reader)
	if err != nil {
		return err
	}

	byteArray, err := ReadNBytes(reader, int(stringLength))
	if err != nil {
		return err
	}

	*value = String(byteArray)
	return nil
}

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

// Decode a VarInt
func (value *VarInt) Decode(reader FieldReader) error {
	var temp uint32
	for i := 0; ; i++ {
		currentByte, err := reader.ReadByte()
		if err != nil {
			return err
		}

		temp |= uint32(currentByte&0x7F) << uint32(7*i)

		if currentByte&0x80 == 0 {
			break
		} else if i > 5 {
			return errors.New("too big VarInt")
		}
	}

	*value = VarInt(temp)
	return nil
}

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

// Decode a VarLong
func (value *VarLong) Decode(reader FieldReader) error {
	var temp uint64
	for i := 0; ; i++ {
		currentByte, err := reader.ReadByte()
		if err != nil {
			return err
		}

		temp |= uint64(currentByte&0x7F) << uint64(7*i)

		if currentByte&0x80 == 0 {
			break
		} else if i > 10 {
			return errors.New("too big VarLong")
		}
	}

	*value = VarLong(temp)
	return nil
}

// Encode a Position
func (value Position) Encode() []byte {
	encoded := make([]byte, 8)
	position := uint64(value.X&0x3FFFFFF)<<38 | uint64((value.Z&0x3FFFFFF)>>12) | uint64(value.Y&0xFFF)
	for i := 7; i >= 0; i-- {
		encoded[i] = byte(position)
		position >>= 8
	}
	return encoded
}

// Decode a Position
func (value *Position) Decode(reader FieldReader) error {
	var read Long
	if err := read.Decode(reader); err != nil {
		return err
	}

	x := int(read >> 38)
	y := int(read & 0xFFF)
	z := int(read << 26 >> 38)

	if x >= 1<<25 {
		x -= 1 << 26
	}
	if y >= 1<<11 {
		y -= 1 << 12
	}
	if z >= 1<<25 {
		z -= 1 << 26
	}

	value.X, value.Y, value.Z = x, y, z
	return nil
}

// Encode a wrapped byte array
func (value ByteArray) Encode() []byte {
	return append(VarInt(len(value)).Encode(), value...)
}

// Decode a wrapped byte array
func (value *ByteArray) Decode(reader FieldReader) error {
	var arrayLength VarInt
	err := arrayLength.Decode(reader)
	if err != nil {
		return err
	}

	*value = make([]byte, arrayLength)
	_, err = reader.Read(*value)
	return err
}

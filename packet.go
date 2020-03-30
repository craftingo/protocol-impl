package protocol_impl

import (
	"bytes"
	"compress/zlib"
	"errors"
	"io"
)

// Packet represents a protocol packet
type Packet struct {
	ID   byte
	Data []byte
}

// GeneratePacket generates a new packet with the given fields
func GeneratePacket(id byte, data ...Encodable) (packet Packet) {
	packet.ID = id
	for _, field := range data {
		packet.Data = append(packet.Data, field.Encode()...)
	}
	return
}

// Unwrap decodes a packet and splits it up into multiple fields
func (packet Packet) Unwrap(fields ...Decodable) error {
	reader := bytes.NewReader(packet.Data)
	for _, field := range fields {
		err := field.Decode(reader)
		if err != nil {
			return err
		}
	}
	return nil
}

// Pack packs a packet
func (packet *Packet) Pack(threshold int) (packed []byte) {
	data := []byte{packet.ID}
	data = append(data, packet.Data...)

	if threshold > 0 {
		if len(data) > threshold {
			length := VarInt(len(data)).Encode()
			data = Compress(data)

			packed = append(packed, VarInt(len(length)+len(data)).Encode()...)
			packed = append(packed, length...)
			packed = append(packed, data...)
		} else {
			packed = append(packed, VarInt(int32(len(data)+1)).Encode()...)
			packed = append(packed, 0x00)
			packed = append(packed, data...)
		}
	} else {
		packed = append(packed, VarInt(int32(len(data))).Encode()...)
		packed = append(packed, data...)
	}

	return
}

// ReadPacket reads a packet from a byte reader
func ReadPacket(reader io.ByteReader, compressed bool) (*Packet, error) {
	var length int
	for i := 0; i < 5; i++ {
		byte, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}

		length |= int(byte&0x7F) << uint(7*i)
		if byte&0x80 == 0 {
			break
		}
	}

	if length < 1 {
		return nil, errors.New("too short packet")
	}

	data := make([]byte, length)
	var err error
	for i := 0; i < length; i++ {
		data[i], err = reader.ReadByte()
		if err != nil {
			return nil, err
		}
	}

	if compressed {
		return Decompress(data)
	}

	return &Packet{
		ID:   data[0],
		Data: data[1:],
	}, nil
}

// Compress compresses a byte array using zlib
func Compress(data []byte) []byte {
	var buffer bytes.Buffer

	writer := zlib.NewWriter(&buffer)
	_, _ = writer.Write(data)
	_ = writer.Close()

	return buffer.Bytes()
}

// Decompress decompresses a byte array into a packet
func Decompress(data []byte) (*Packet, error) {
	reader := bytes.NewReader(data)

	var size VarInt
	err := size.Decode(reader)
	if err != nil {
		return nil, err
	}

	decompressed := make([]byte, size)
	if size != 0 {
		zlibReader, err := zlib.NewReader(reader)
		if err != nil {
			return nil, err
		}
		defer zlibReader.Close()

		_, err = io.ReadFull(zlibReader, decompressed)
		if err != nil {
			return nil, err
		}
	} else {
		decompressed = data[1:]
	}

	return &Packet{
		ID:   decompressed[0],
		Data: decompressed[1:],
	}, nil
}

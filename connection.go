package protocol_impl

import (
	"bufio"
	"crypto/cipher"
	"io"
	"net"
)

// Connection represents a Minecraft connection
type Connection struct {
	Socket net.Conn
	io.ByteReader
	io.Writer
	threshold int
}

// Listener represents a connection listener
type Listener struct {
	net.Listener
}

// Listen listens for new Minecraft connections
func Listen(address string) (*Listener, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	return &Listener{listener}, nil
}

// Accept accepts a new Minecraft connection
func (listener Listener) Accept() (Connection, error) {
	connection, err := listener.Listener.Accept()
	return Connection{
		Socket:     connection,
		ByteReader: bufio.NewReader(connection),
		Writer:     connection,
	}, err
}

// SetThreshold sets the threshold for a Minecraft connection
func (connection *Connection) SetThreshold(threshold int) {
	connection.threshold = threshold
}

// SetCipher defines the cipher streams for a Minecraft connection
func (connection *Connection) SetCipher(encodeStream, decodeStream cipher.Stream) {
	connection.ByteReader = bufio.NewReader(cipher.StreamReader{
		S: decodeStream,
		R: connection.Socket,
	})

	connection.Writer = cipher.StreamWriter{
		S: encodeStream,
		W: connection.Socket,
	}
}

// ReadPacket reads a packet from a Minecraft connection
func (connection *Connection) ReadPacket() (Packet, error) {
	packet, err := ReadPacket(connection.ByteReader, connection.threshold > 0)
	if err != nil {
		return Packet{}, err
	}
	return *packet, nil
}

// WritePacket writes a packet to to a Minecraft connection
func (connection *Connection) WritePacket(packet Packet) error {
	_, err := connection.Write(packet.Pack(connection.threshold))
	return err
}

// Close closes a Minecraft connection
func (connection Connection) Close() error {
	return connection.Socket.Close()
}

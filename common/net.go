package common

import (
	"encoding/binary"
	"io"

	"github.com/golang/protobuf/proto"
)

// The first byte of a network message indicates its type
const (
	MsgWelcome byte = iota
	MsgAuthResponse
	MsgAuthResult
)

// SendMessage sends a protobuf message along with a header containing its type and length
func SendMessage(dst io.Writer, msgtype byte, msg proto.Message) {
	out, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}

	header := make([]byte, 5)
	header[0] = msgtype
	binary.BigEndian.PutUint32(header[1:], uint32(len(out)))

	dst.Write(header)
	dst.Write(out)
}

// ReadMessage reads a protobuf message along with a header containing its type and length
func ReadMessage(src io.Reader) (byte, []byte, error) {
	header := make([]byte, 5)
	src.Read(header)

	msgtype := header[0]
	length := binary.BigEndian.Uint32(header[1:5])

	buf := make([]byte, length)
	_, err := src.Read(buf)
	if err != nil {
		return 0, nil, err
	}

	return msgtype, buf, nil
}

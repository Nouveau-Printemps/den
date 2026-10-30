package proto

import (
	"encoding/binary"
	"encoding/json/v2"
	"errors"
	"io"
)

type Message struct {
	Code uint8
	Body io.Reader
}

func (msg *Message) As[T any]() (T, error) {
	var v T
	err := json.UnmarshalRead(msg.Body, &v)
	return v, err
}

var (
	ErrInvalidMessage = errors.New("invalid message")
)

func ReadMessage(r io.Reader) (*Message, error) {
	var header [4]byte
	_, err := io.ReadFull(r, header[:])
	if err != nil {
		return nil, err
	}
	ln := binary.BigEndian.Uint16(header[:])
	//flags := header[2]
	code := header[3]
	return &Message{Code: code, Body: io.LimitReader(r, int64(ln))}, nil
}

func SendMessage[C ResponseCode | CommandCode](w io.Writer, code C, data any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	var header [4]byte
	binary.BigEndian.PutUint16(header[:], uint16(len(b)))
	header[3] = uint8(code)
	_, err = w.Write(header[:])
	if err != nil {
		return err
	}
	return w.Write(b)
}

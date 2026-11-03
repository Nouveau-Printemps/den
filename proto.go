package den

import (
	"encoding/binary"
	. "encoding/json/jsontext"
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

// ReadMessage is a low-level function reading a [Message] from an [io.Reader].
//
// For a higher-level usage, see [Client].
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

// SendMessage is a low-level function sending a [Message] to an [io.Writer].
//
// For a higher-level usage, see [Client].
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
	_, err = w.Write(b)
	return err
}

type CommandCode uint8

const (
	CodeHello CommandCode = iota
	CodeAuthn
	CodeExists
)

type CommandHello struct {
	Name string `json:"name"`
}

type CommandAuthn struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CommandExists struct {
	Username string `json:"username"`
}

type ResponseCode uint8

// Available [ResponseCode]
const (
	CodeOk ResponseCode = iota
	// CodeInvalid is returned when the message is malformed.
	CodeInvalid
	// CodeBad is returned when the data inside the message is invalid.
	CodeBad
	CodeInternalError
)

type ResponseError struct {
	Details string `json:"details"`
}

func (e *ResponseError) Error() string {
	return e.Details
}

type ResponseAuthn struct {
	Key    *string  `json:"pgp_key,omitempty"`
	Groups []string `json:"groups,omitempty"`
	// RawData is populated when unmarshalled
	RawData Value `json:"data,omitempty"`
	// CustomData is used when marshalled.
	CustomData any `json:"-"`
}

func (r *ResponseAuthn) MarshalJSONTo(e *Encoder) error {
	e.WriteToken(BeginObject)
	if r.Key != nil {
		e.WriteToken(String("pgp_key"))
		e.WriteToken(String(*r.Key))
	}
	if len(r.Groups) != 0 {
		e.WriteToken(String("groups"))
		e.WriteToken(BeginArray)
		for _, g := range r.Groups {
			e.WriteToken(String(g))
		}
		e.WriteToken(EndArray)
	}
	if r.CustomData != nil {
		out, err := json.Marshal(r.CustomData)
		if err != nil {
			return err
		}
		err = e.WriteValue(Value(out))
		if err != nil {
			return err
		}
	}
	return nil
}

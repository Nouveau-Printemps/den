package den

import (
	. "encoding/json/jsontext"
	"encoding/json/v2"
)

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

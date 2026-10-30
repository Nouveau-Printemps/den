package den

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
	Key *string `json:"pgp_key,omitempty"`
}

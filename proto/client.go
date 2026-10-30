package proto

type CommandCode uint8

const (
	CodeHello CommandCode = iota
	CodeAuthn
)

type CommandHello struct {
	Name string `json:"name"`
}

type CommandAuthn struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

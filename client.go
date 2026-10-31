package den

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

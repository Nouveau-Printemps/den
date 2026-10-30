package den

import (
	"io"
	"net"

	"github.com/nouveau-printemps/den/proto"
)

type Client struct {
	conn net.Conn
}

func Dial(servername, network, address string) (*Client, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	err = proto.SendMessage(conn, proto.CodeHello, proto.CommandHello{Name: servername})
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Authn(username, password string) (bool, error) {
	err := proto.SendMessage(c.conn, proto.CodeAuthn, proto.CommandAuthn{
		Username: username,
		Password: password,
	})
	if err != nil {
		return false, err
	}
	msg, err := proto.ReadMessage(c.conn)
	if err != nil {
		return false, err
	}
	switch proto.ResponseCode(msg.Code) {
	case proto.CodeOk:
		_, _ = io.ReadAll(msg.Body)
		return true, nil
	case proto.CodeBad:
		_, _ = io.ReadAll(msg.Body)
		return false, nil
	case proto.CodeInternalError:
		e, err := msg.As[*proto.ResponseError]()
		if err != nil {
			return false, err
		}
		return false, e
	case proto.CodeInvalid:
		_, _ = io.ReadAll(msg.Body)
		return false, proto.ErrInvalidMessage
	default:
		panic("invalid response code")
	}
}

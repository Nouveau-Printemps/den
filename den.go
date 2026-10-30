package den

import (
	"crypto/tls"
	"errors"
	"io"
	"net"

	"github.com/nouveau-printemps/den/proto"
)

type Client struct {
	conn net.Conn
}

// Dial creates a new connection to the server.
//
// See [net.Dial] for more information.
func Dial(clientname, network, address string) (*Client, error) {
	var c Client
	var err error
	c.conn, err = net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return &c, c.hello(clientname)
}

// Dial creates a new connection to the server using TLS.
//
// See [tls.Dial] for more information.
func DialTLS(clientname, network, address string, cfg *tls.Config) (*Client, error) {
	var c Client
	var err error
	c.conn, err = tls.Dial(network, address, cfg)
	if err != nil {
		return nil, err
	}
	return &c, c.hello(clientname)
}

func (c *Client) hello(servername string) error {
	err := proto.SendMessage(c.conn, proto.CodeHello, proto.CommandHello{Name: servername})
	if err != nil {
		return err
	}
	msg, err := proto.ReadMessage(c.conn)
	if err != nil {
		return err
	}
	_, _ = io.ReadAll(msg.Body)
	if msg.Code != uint8(proto.CodeOk) {
		return errors.New("den: cannot connect to server")
	}
	return nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

// Authn verifies the login credentials.
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

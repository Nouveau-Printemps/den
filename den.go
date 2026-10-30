package den

import (
	"crypto/tls"
	"errors"
	"io"
	"net"
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
	err := SendMessage(c.conn, CodeHello, CommandHello{Name: servername})
	if err != nil {
		return err
	}
	msg, err := ReadMessage(c.conn)
	if err != nil {
		return err
	}
	_, _ = io.ReadAll(msg.Body)
	if msg.Code != uint8(CodeOk) {
		return errors.New("den: cannot connect to server")
	}
	return nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

// Authn verifies the login credentials.
func (c *Client) Authn(username, password string) (*ResponseAuthn, error) {
	err := SendMessage(c.conn, CodeAuthn, CommandAuthn{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	msg, err := ReadMessage(c.conn)
	if err != nil {
		return nil, err
	}
	switch ResponseCode(msg.Code) {
	case CodeOk:
		return msg.As[*ResponseAuthn]()
	case CodeBad:
		return nil, nil
	case CodeInternalError:
		e, err := msg.As[*ResponseError]()
		if err != nil {
			return nil, err
		}
		return nil, e
	case CodeInvalid:
		_, _ = io.ReadAll(msg.Body)
		return nil, ErrInvalidMessage
	default:
		panic("invalid response code")
	}
}

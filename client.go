package den

import (
	"crypto/tls"
	"errors"
	"io"
	"net"
)

// Client is connected to a den server.
//
// Use [Dial] or [DialTLS] to connect a new one.
//
// A Client is safe for concurrent use.
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

// Close a [Client].
// It must not be used after.
func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) response[T any]() (v T, err error) {
	msg, err := ReadMessage(c.conn)
	if err != nil {
		return
	}
	switch ResponseCode(msg.Code) {
	case CodeOk:
		return msg.As[T]()
	case CodeBad:
		_, _ = io.ReadAll(msg.Body)
		return
	case CodeInternalError:
		var e error
		e, err = msg.As[*ResponseError]()
		if err != nil {
			return
		}
		return v, e
	case CodeInvalid:
		_, _ = io.ReadAll(msg.Body)
		return v, ErrInvalidMessage
	default:
		panic("invalid response code")
	}
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
	return c.response[*ResponseAuthn]()
}

// Exists verifies if a user exists.
func (c *Client) Exists(username string) (*ResponseAuthn, error) {
	err := SendMessage(c.conn, CodeExists, CommandExists{
		Username: username,
	})
	if err != nil {
		return nil, err
	}
	return c.response[*ResponseAuthn]()
}

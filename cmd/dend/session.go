package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"

	"github.com/Nouveau-Printemps/den"
)

type Session struct {
	name string
	conn net.Conn
}

func NewSession(conn net.Conn) (*Session, error) {
	msg, err := den.ReadMessage(conn)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			err = den.SendMessage(conn, den.CodeInvalid, &den.ResponseError{
				Details: "expected hello command",
			})
			if err != nil {
				slog.Error("sending invalid command", "error", err)
			}
		}
	}()
	if msg.Code != uint8(den.CodeHello) {
		return nil, fmt.Errorf("invalid code: %d", msg.Code)
	}
	hello, err := msg.As[*den.CommandHello]()
	if err != nil {
		return nil, err
	}
	return &Session{hello.Name, conn}, nil
}

func (s *Session) Handle(cfg *Config, conn net.Conn) {
	for {
		msg, err := den.ReadMessage(conn)
		if errors.Is(err, io.EOF) {
			return
		}
		if err != nil {
			den.SendMessage(conn, den.CodeInvalid, &den.ResponseError{
				Details: "invalid command",
			})
			continue
		}
		switch den.CommandCode(msg.Code) {
		case den.CodeAuthn:
			authn, err := msg.As[*den.CommandAuthn]()
			if err != nil {
				den.SendMessage(conn, den.CodeInvalid, &den.ResponseError{
					Details: "invalid body",
				})
				continue
			}
			user := cfg.Authn(authn.Username, authn.Password)
			if user == nil {
				err = den.SendMessage(conn, den.CodeBad, &den.ResponseError{
					Details: "bad credentials",
				})
			} else {
				err = den.SendMessage(conn, den.CodeOk, &den.ResponseAuthn{
					Key: user.Key(),
				})
			}
			if err != nil {
				slog.Error("responding", "error", err)
			}
		case den.CodeExists:
			exists, err := msg.As[*den.CommandExists]()
			if err != nil {
				den.SendMessage(conn, den.CodeInvalid, &den.ResponseError{
					Details: "invalid body",
				})
				continue
			}
			user := cfg.Exists(exists.Username)
			if user == nil {
				err = den.SendMessage(conn, den.CodeBad, &den.ResponseError{
					Details: "not found",
				})
			} else {
				err = den.SendMessage(conn, den.CodeOk, &den.ResponseAuthn{
					Key: user.Key(),
				})
			}
			if err != nil {
				slog.Error("responding", "error", err)
			}
		default:
			den.SendMessage(conn, den.CodeInvalid, &den.ResponseError{
				Details: "invalid code",
			})
		}
	}
}

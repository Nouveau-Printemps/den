package main

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	Users map[string]*User `toml:"user"`
}

type User struct {
	Password      string  `toml:"password"`
	PGPPubKey     *string `toml:"pgp_pub_key"`
	PGPPubKeyFile *string `toml:"pgp_pub_key_file"`
}

func (cfg *Config) Exists(username string) bool {
	_, ok := cfg.Users[username]
	return ok
}

func (cfg *Config) Authn(username, password string) *User {
	u, ok := cfg.Users[username]
	if !ok {
		time.Sleep(1 * time.Second)
		return nil
	}
	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) != nil {
		time.Sleep(1 * time.Second)
		return nil
	}
	return u
}

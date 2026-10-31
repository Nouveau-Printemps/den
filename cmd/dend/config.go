package main

import (
	_ "embed"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"golang.org/x/crypto/bcrypt"
)

//go:embed default.toml
var defaultConfig []byte

type Config struct {
	Users map[string]*User `toml:"user"`
}

type User struct {
	Password      string   `toml:"password"`
	PGPPubKey     *string  `toml:"pgp_pub_key"`
	PGPPubKeyFile *string  `toml:"pgp_pub_key_file"`
	Groups        []string `toml:"groups"`
}

func (cfg *Config) Exists(username string) *User {
	return cfg.Users[username]
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

func (u *User) Key() *string {
	if u.PGPPubKey != nil {
		return u.PGPPubKey
	}
	if u.PGPPubKeyFile != nil {
		b, err := os.ReadFile(*u.PGPPubKeyFile)
		if err != nil {
			slog.Warn("reading pgp key file", "error", err, "file", *u.PGPPubKeyFile)
			return nil
		}
		return new(string(b))
	}
	return nil
}

func Parse() (cfg *Config) {
	data, err := toml.DecodeFile(config, &cfg)
	if errors.Is(err, os.ErrNotExist) {
		err = os.WriteFile(config, defaultConfig, 0o600)
		if err == nil {
			err = toml.Unmarshal(defaultConfig, &cfg)
		}
	}
	if err != nil {
		panic(err)
	}
	for _, k := range data.Undecoded() {
		slog.Warn("key not decoded", "key", k)
	}
	return
}

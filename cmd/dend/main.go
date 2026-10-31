package main

import (
	"context"
	"flag"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
)

var (
	verbose bool
	config  string = "/etc/dend/config.toml"
	address string = "/run/dend.sock"
	uid     uint   = 0
	gid     uint   = 0
)

var Version = "dev"

func init() {
	flag.BoolVar(&verbose, "v", verbose, "increase verbosity")
	flag.StringVar(&config, "config", config, "path to the config")
	flag.StringVar(&address, "address", address, "address to listen to")
	flag.UintVar(&uid, "user", uid, "uid to run as")
	flag.UintVar(&gid, "group", gid, "gid to run as")
}

func main() {
	flag.Parse()
	ctx, cancel := signal.NotifyContext(
		context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM,
	)
	slog.Info("starting...", "version", Version)
	defer cancel()
	var cfg Config
	data, err := toml.DecodeFile(config, &cfg)
	if err != nil {
		panic(err)
	}
	for _, k := range data.Undecoded() {
		slog.Warn("key not decoded", "key", k)
	}
	network := "tcp"
	unix := address[0] == '/'
	if unix {
		network = "unix"
	}
	l, err := new(net.ListenConfig).Listen(ctx, network, address)
	if err != nil {
		panic(err)
	}
	if unix {
		err = os.Chmod(address, 0o660)
		if err != nil {
			panic(err)
		}
		err = os.Chown(address, int(uid), int(gid))
		if err != nil {
			panic(err)
		}
	}
	slog.Debug("changing uid and gid", "uid", uid)
	err = syscall.Setuid(int(uid))
	if err != nil {
		panic(err)
	}
	err = syscall.Setgid(int(gid))
	if err != nil {
		panic(err)
	}
	errc := make(chan error)
	go func() {
		errc <- Serve(&cfg, l)
	}()
	slog.Info("started")
	select {
	case err = <-errc:
		panic(err)
	case <-ctx.Done():
	}
	slog.Info("stopping")
}

func Serve(cfg *Config, l net.Listener) error {
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go func() {
			s, err := NewSession(conn)
			defer conn.Close()
			if err != nil {
				slog.Error("handling connection", "error", err)
				return
			}
			s.Handle(cfg, conn)
		}()
	}
}

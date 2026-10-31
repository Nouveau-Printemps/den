package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"log/syslog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/nyttikord/logos"
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
	flag.UintVar(&uid, "uid", uid, "uid to run as")
	flag.UintVar(&gid, "gid", gid, "gid to run as")
}

func main() {
	flag.Parse()
	lvl := slog.LevelInfo
	if verbose {
		lvl = slog.LevelDebug
	}
	logos, err := logos.NewSyslog("dend", syslog.LOG_AUTH, &logos.Options{
		Level: lvl,
	})
	if err != nil {
		panic(err)
	}
	slog.SetDefault(slog.New(logos))
	ctx, cancel := signal.NotifyContext(
		context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM,
	)
	slog.Info("starting...", "version", Version)
	defer cancel()
	var cfg Config
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
	ctx = withConfig(withLogger(ctx, slog.Default()), &cfg)
	network := "tcp"
	unix := address[0] == '/'
	if unix {
		network = "unix"
	}
	l, err := new(net.ListenConfig).Listen(ctx, network, address)
	if err != nil {
		panic(err)
	}
	defer l.Close()
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
	slog.Debug("changing uid and gid", "uid", uid, "gid", gid)
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
		errc <- Serve(ctx, l)
	}()
	slog.Info("started")
	select {
	case err = <-errc:
		panic(err)
	case <-ctx.Done():
	}
	slog.Info("stopping")
}

func Serve(ctx context.Context, l net.Listener) error {
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go func() {
			s, err := NewSession(conn)
			defer conn.Close()
			if err != nil {
				logger(ctx).Error("handling connection", "error", err)
				return
			}
			s.Handle(
				withLogger(ctx, logger(ctx).With(
					"module", "session",
					"remote", conn.RemoteAddr(),
					"name", s.name,
				)),
				conn,
			)
		}()
	}
}

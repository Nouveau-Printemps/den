package main

import (
	"context"
	"flag"
	"log/slog"
	"log/syslog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/nyttikord/logos"
)

var (
	verbose bool
	config  string = "/etc/dend/config.toml"
	address string = "/run/den/dend.sock"
)

var Version = "dev"

func init() {
	flag.BoolVar(&verbose, "v", verbose, "increase verbosity")
	flag.StringVar(&config, "config", config, "path to the config")
	flag.StringVar(&address, "address", address, "address to listen to")
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
	cfg := Parse()
	ctx = withConfig(withLogger(ctx, slog.Default()), cfg)
	network := "tcp"
	unix := address[0] == '/'
	if unix {
		network = "unix"
	}
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP)
		for {
			select {
			case <-ctx.Done():
				return
			case <-sig:
				*cfg = *Parse()
			}
		}
	}()
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

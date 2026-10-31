version := `git describe --tags`

dev:
    cd cmd/dend && go run . -config dev.toml -address :8000 -v -uid $(id -u) -gid $(id -g)

dev-log:
    journalctl --user -S today --facility auth

build:
    cd cmd/dend && go build -ldflags "-s -w -X main.Version={{ version }}" -o ../../dend .

dev:
    cd cmd/dend && go run . -config dev.toml -address :8000 -v -uid $(id -u) -gid $(id -g)

dev-log:
    journalctl --user -S today --facility auth

build: build-server
    mv pkg/server/server dend

[working-directory: 'pkg/server']
build-server:
    go build -ldflags "-s -w -X main.Version={{ `git describe --tags` }}" .

install:
    install -m 755 ./dend /usr/local/sbin/

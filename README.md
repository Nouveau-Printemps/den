# den

Simple authentification protocol without authorization.

This repo contains a server, a client library and a formal definition of the protocol.

## Why

At Nouveau Printemps, we have multiple applications that share the same users.
Already existing protocols are too complex or not adapted to our needs: our applications manage theirselves
authorization, but we wanted to share the authentification.
Den is designed to fix this: centralize authentification without touching authorization.

## Building

You can build the server with
```shell
$ just build
```

The server is designed to only work with UNIX and UNIX-like OSes.
Only Linux is tested.

## Organization

The root folder contains the formal definition of the protocol and a Go client library.

[`cmd/dend/`](./cmd/dend/) is a separate Go module: it implements a complete server with the root library.

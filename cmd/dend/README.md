# dend

`dend` is the official server implementation of the den protocol.

## Usage

User information is stored in a config file: `/etc/dend/config.toml` by default.
It contains its username, its password hashed with bcrypt, its groups and its public PGP key.

```toml
[user]
foo.password = "$2y$10$wWJy64n5j4pU2CscdWRq.er5Z2.V31U5Ntz5uTopLy1fMX87GQEDG"
foo.groups = ["bar"]
```

The server listens to `/run/dend.sock` by default.
You can customize who executes the daemon with `-uid` and `gid` flags.
Other uids and gids can't communicate to the server.

The daemon logs to the `AUTH` `syslog(3)` facility.

## Building

See the root [`README.md`](/README.md)

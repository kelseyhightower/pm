# pm

pm is a simple package manager for self-contained binaries.

## Usage

### Build

```
$ pm build -b hello -m metadata.json -p private.key
```
```
hello-1.0.0-darwin-amd64.tar.gz
hello-1.0.0-darwin-amd64.tar.gz.sha256
hello-1.0.0-darwin-amd64.tar.gz.asc
hello-1.0.0-darwin-amd64.tar.gz.sha256.asc
```

### Download 

The get subcommand downloads packages into the pm cache directory `/opt/pm/cache`. The remote package will only be downloaded if the remote copy is newer than the local copy.

```
$ pm get https://dl.example.com/hello-1.0.0-darwin-amd64.tar.gz
```

### Install

```
$ pm install https://dl.example.com/hello-1.0.0-darwin-amd64.tar.gz
```

```
$ pm install hello-1.0.0-darwin-amd64.tar.gz
```

### List

```
$ pm list
```
```
hello-1.0.0
```

### Remove

```
$ pm remove hello-1.0.0
```

### Verify

Trusted public keys must be added to `/opt/pm/.keyring`

```
$ pm verify hello-1.0.0
```

```
$ pm verify hello-1.0.0-darwin-amd64.tar.gz
```

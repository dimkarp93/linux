## Overview
Example with simple http-server and http-client with customization of local and remote port.
Server wait some fixed time and return client its remote address.

### Getting started
Run server than run client.
Before doing command go to dir with this manual

#### Server
```shell
go run server/server.go
```

#### Client
```shell
go run client/server.go
```

### Options

#### Server
- `--host` - address for bind
- `--port` - port for bind

#### Client
- `--remoteHost` - address of server
- `--remotePort` - port of server
- `--host` - local address to send request
- `--port` - local port to send request
### go-http-cache
A simple http proxy that writes data to filesystem and serves from cache.

Overrides
|Name|Purpose|
|---|---|
|CONFIG_FILE|Config File to be used|
|LISTEN_ADDRESS|Address this server will listen on|
|TARGET|Proxy Target to call when cache does not have data|
|IGNORE_HEADERS|Whether to consider header keys and values when looking up cache|

### Running Locally
```
LISTEN_ADDRESS=:8070 TARGET=http://localhost:8080 go run app.go
```
Doing this will listen on 8070 and proxy calls to 8080. Filesystem path /tmp/cache is used as cache location from which subsequent responses will be served. 

### Docker
Same thing can be run in docker like this.

```
docker run -e TARGET=http://someendpoint:8080 LISTEN_ADDRESS=:8070 -p 8070:8070 -it skhatri/go-http-cache:latest
```


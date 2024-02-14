### go-http-cache
A simple http proxy that writes data to filesystem and serves from cache.

Overrides
|Name|Purpose|Default|
|---|---|---|
|CONFIG_FILE|Config File to be used|config.yaml|
|LISTEN_ADDRESS|Address this server will listen on|8070|
|TARGET|Proxy Target to call when cache does not have data|http://localhost:8080|
|IGNORE_HEADERS|Whether to consider header keys and values when looking up cache. Each header is computed when determining cache key|false|
|LOG_REQUEST_HEADERS|Whether to log request headers into cached data file. Be careful if you are using Auth Tokens|false|
|SKIP_VERIFY_TLS|Whether to verify certs for https calls|false|

### Running Locally
```
TARGET=https://jsonplaceholder.typicode.com SKIP_VERIFY_TLS=true go run app.go
```
Doing this will listen on 8070 and proxy calls to 8080. Filesystem path /tmp/cache is used as cache location from which subsequent responses will be served. 

### Docker
Same thing can be run in docker like this.

```
docker run -e TARGET=https://jsonplaceholder.typicode.com \
-e LISTEN_ADDRESS=0.0.0.0:8070 \
-e LOG_REQUEST_HEADERS=true \
-e SKIP_VERIFY_TLS=true -p 8070:8070 \
-it skhatri/go-http-cache:latest
```

Test it by doing a curl
```curl http://localhost:8070/todos```



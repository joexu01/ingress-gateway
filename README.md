To start reverse proxy:

```shell
$ go run main.go
```

To start reverse proxy without security verifications*:

```shell
$ SEC_VER=disabled go run main.go
```

*`proxy_http_router/router.go`
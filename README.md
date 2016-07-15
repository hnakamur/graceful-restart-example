graceful-restart-example
========================

An example of graceful restart of go webapp using [zenazn/goji: Goji is a minimalistic web framework for Golang that's high in antioxidants.](https://github.com/zenazn/goji) and [lestrrat/go-server-starter: Go port of start_server utility (Server::Starter)](https://github.com/lestrrat/go-server-starter).

## How to set up on CentOS 7

Install Go: [Getting Started - The Go Programming Language](https://golang.org/doc/install).

Install the `start_server` command in [lestrrat/go-server-starter](https://github.com/lestrrat/go-server-starter).

```
go get github.com/lestrrat/go-server-starter/cmd/start_server
```

Install the `graceful-restart-example` command from this repository.

```
go get github.com/hnakamur/graceful-restart-example
```

```
sed "|/usr/local/gocode|$GOPATH|" graceful-restart-example.service > /etc/systemd/system/graceful-restart-example.service
systemctl daemon-reaload
```

## Example session

```
$ sudo systemctl start graceful-restart-example
$ : > curl.log && for i in `seq 1 100`; do curl -s 127.0.0.1:7777/sleep/$i >> curl.log & done; sudo systemctl reload graceful-restart-example & for i in `seq 101 200`; do curl -s 127.0.0.1:7777/sleep/$i >> curl.log & done & tail -f curl.log
```

```
$ less /tmp/graceful-restart-example.log
```

## TODO
* Switch to [goji/goji](https://github.com/goji/goji)

## License
MIT

# httpseek [![GoDoc](https://godoc.org/github.com/hnz/httpseek?status.svg)](https://godoc.org/github.com/hnz/httpseek) [![Go Report Card](https://goreportcard.com/badge/github.com/hnz/httpseek)](https://goreportcard.com/report/github.com/hnz/httpseek) [![coverage](https://img.shields.io/codacy/coverage/c44df2d9c89a4809896914fd1a40bedd.svg)](https://gocover.io/github.com/hnz/httpseek)

httpseek is a go package that implements io.Seeker and io.SectionReader interface on HTTP Response bodies.
This allows the client to read parts of a file without downloading it entirely.

To achieve this the HTTP Range request is used. If the HTTP server does not support this,
the request will fail.

Buffering is not supported, so reading many small chunks of data will result in many requests to
the http server, which could generate more load than downloading the whole file at once.


Install library
---------------

    go get github.com/hnz/httpseek


License
-------

MIT Licensed. See [LICENSE.md](LICENSE.md).

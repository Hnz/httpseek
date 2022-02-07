# httpseek [![GoDoc](https://godoc.org/github.com/hnz/httpseek?status.svg)](https://godoc.org/github.com/hnz/httpseek) [![Go Report Card](https://goreportcard.com/badge/github.com/hnz/httpseek)](https://goreportcard.com/report/github.com/hnz/httpseek)

httpseek is a go package that implements io.Seeker and io.ReaderAt interface on HTTP Response bodies.
This allows the client to read parts of a file without downloading it entirely.

To achieve this the HTTP Range request is used. If the HTTP server does not support this,
the request will fail.

To prevent lots of tiny http request a buffer is implemented.


Install library
---------------

    go get github.com/hnz/httpseek


License
-------

MIT Licensed. See [LICENSE.md](LICENSE.md).

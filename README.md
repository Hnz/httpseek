# httpseek

[![Version](https://img.shields.io/github/v/tag/hnz/httpseek?label=version&sort=semver&style=for-the-badge)](https://github.com/hnz/httpseek/tags)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white&style=for-the-badge)](https://pkg.go.dev/github.com/hnz/httpseek#section-readme)
[![Go Report Card](https://goreportcard.com/badge/github.com/hnz/httpseek?style=for-the-badge)](https://goreportcard.com/report/github.com/hnz/httpseek)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/hnz/httpseek/go.yml?style=for-the-badge)](https://github.com/hnz/httpseek/actions/workflows/go.yml)

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

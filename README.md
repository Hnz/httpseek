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

    Copyright (c) 2022 Hans van Leeuwen

    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:

    The above copyright notice and this permission notice shall be included in
    all copies or substantial portions of the Software.

    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
    THE SOFTWARE.
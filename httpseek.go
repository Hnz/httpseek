// Copyright (c) 2022 Hans van Leeuwen. MIT Licensed. See README.md for full license.

/*
httpseek is a go package that implements io.Seeker and io.SectionReader interface on HTTP Response bodies.
This allows the client to read parts of a file without downloading it entirely.

To achieve this the HTTP Range request is used. If the HTTP server does not support this,
the request will fail.

Buffering is not supported, so reading many small chunks of data will result in many requests to
the http server, which could generate more load than downloading the whole file at once.
*/
package httpseek

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// See https://pkg.go.dev/net/http#Client
type Client struct {
	*http.Client
}

// See https://pkg.go.dev/net/http#Client.Do
func (c *Client) Do(req *http.Request) (*Response, error) {
	if req.Method == "GET" || req.Method == "" {
		// Retrieve the content-lenght of the request using head, but don't make the GET request yet.

		if c.Client == nil {
			c.Client = new(http.Client)
		}

		head, err := c.Client.Head(req.URL.String())
		r := Response{Response: head}
		if err != nil {
			return &r, err
		}

		if head.Header.Get("Accept-Ranges") != "bytes" {
			log.Fatalln("HTTP Server does not accept ranges")
		}

		r.Request = req
		r.Body.request = req
		r.Body.client = c
		r.Body.contentlength = head.ContentLength
		r.ContentLength = head.ContentLength
		return &r, err
	}

	res, err := c.Client.Do(req)
	return &Response{Response: res}, err
}

// See https://pkg.go.dev/net/http#Client.Get
func (c *Client) Get(url string) (resp *Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

type Response struct {
	*http.Response

	Body ResponseBody
}

type ResponseBody struct {
	Body io.ReadCloser

	client        *Client
	contentlength int64
	offset        int64
	request       *http.Request
}

// See https://pkg.go.dev/io#Closer
func (o *ResponseBody) Close() error {
	return o.Body.Close()
}

// See https://pkg.go.dev/io#Reader
func (o *ResponseBody) Read(p []byte) (int, error) {
	return o.ReadAt(p, o.offset)
}

// See https://pkg.go.dev/io#ReaderAt
func (o *ResponseBody) ReadAt(p []byte, off int64) (n int, err error) {
	req, err := http.NewRequest("GET", o.request.URL.String(), nil)
	if err != nil {
		return n, err
	}

	// Add range to the request
	r := fmt.Sprintf("bytes=%d-%d", o.offset, o.offset+int64(len(p)))
	req.Header.Add("Range", r)
	res, err := o.client.Client.Do(req)
	if err != nil {
		return n, err
	}

	return res.Body.Read(p)
}

// See https://pkg.go.dev/io#Seeker
func (o *ResponseBody) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		o.offset = offset
	case io.SeekCurrent:
		o.offset += offset
	case io.SeekEnd:
		o.offset = o.contentlength - offset

	default:
		return 0, os.ErrInvalid
	}
	return o.offset, nil
}

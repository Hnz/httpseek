// Copyright (c) 2022 Hans van Leeuwen. MIT Licensed. See LICENSE.md for full license.

/*
httpseek is a go package that implements io.Seeker and io.ReaderAt interface on HTTP Response bodies.
This allows the client to read parts of a file without downloading it entirely.

To achieve this the HTTP Range request is used. If the HTTP server does not support this,
the request will fail.

To prevent lots of tiny http request a buffer is implemented.
*/
package httpseek

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
)

// Client is an implementation of http.Client that allows for the response body to be seeked
type Client struct {
	*http.Client // If empty use http.DefaultClient

	Blocksize int64
}

// See https://pkg.go.dev/net/http#Client.Do
func (c *Client) Do(req *http.Request) (*Response, error) {
	if req.Method == "GET" || req.Method == "" {
		// Retrieve the content-length of the request using head, but don't make the GET request yet.
		if c.Client == nil {
			c.Client = http.DefaultClient
		}
		url := req.URL.String()
		head, err := c.Head(url)
		if err != nil {
			return nil, err
		}

		// Check if server supports Range
		if head.Header.Get("Accept-Ranges") != "bytes" {
			return nil, &RangeNotSupported{Url: url}
		}

		r := Response{Response: head}
		r.Request = req
		r.ContentLength = head.ContentLength
		r.Body.client = c
		r.Body.contentlength = head.ContentLength
		r.Body.request = req
		r.Body.ReaderAtBuffer.Blocksize = c.Blocksize

		// Store the version of the file so if it changes during the read, we still use the same version
		if lastmod := head.Header.Get("ETag"); lastmod != "" {
			r.Body.lastmod = lastmod
		} else if lastmod := head.Header.Get("Last-Modified"); lastmod != "" {
			r.Body.lastmod = lastmod
		}

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

// See https://pkg.go.dev/net/http#Response
type Response struct {
	*http.Response

	Body ResponseBody
}

// ResponseBody implents the io.ReaderAt and io.Seeker interfaces
type ResponseBody struct {
	ReaderAtBuffer

	Body io.ReadCloser

	client        *Client
	contentlength int64
	lastmod       string // Can be an Etag or a Last-Modified datetime. Will be passed as If-Range header
	offset        int64
	request       *http.Request
}

// See https://pkg.go.dev/io#Closer
func (o *ResponseBody) Close() error {
	return o.Body.Close()
}

// See https://pkg.go.dev/io#Reader
func (o *ResponseBody) Read(p []byte) (n int, err error) {
	print("READ", len(p))
	n, err = o.ReadAt(p, o.offset)
	if err == nil {
		o.offset += int64(n)
		print(n, o.offset)
	}
	return n, err
}

// See https://pkg.go.dev/io#ReaderAt
func (o *ResponseBody) ReadAt(p []byte, off int64) (n int, err error) {
	print("READAT", len(p), off, o.offset)
	req, err := http.NewRequest("GET", o.request.URL.String(), nil)
	if err != nil {
		return n, err
	}

	// Add Range header to the request
	r := fmt.Sprintf("bytes=%d-%d", off, off+int64(len(p)))
	req.Header.Add("Range", r)

	// Add If-Range header to the request
	if o.lastmod != "" {
		req.Header.Add("If-Range", o.lastmod)
	}

	// Do the request with http.Client.
	res, err := o.client.Client.Do(req)
	if err != nil {
		return n, err
	}

	return res.Body.Read(p)
}

// See https://pkg.go.dev/io#Seeker
func (o *ResponseBody) Seek(offset int64, whence int) (int64, error) {
	print("SEEK", offset, whence)
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

type ReaderAtBuffer struct {
	Blocksize     int64
	ContentLength int64
	ReaderAt      io.ReaderAt

	buffer [][]byte
}

func (o *ReaderAtBuffer) ReadAt(p []byte, off int64) (n int, err error) {

	// If Blocksize == 0 buffering is disabled
	if o.Blocksize == 0 {
		return o.ReaderAt.ReadAt(p, off)
	}

	if o.ContentLength == 0 {
		panic("ContentLength not set")
	}
	if o.ReaderAt == nil {
		panic("ReaderAt not set")
	}

	var out []byte
	bytesrequested := int64(len(p))

	blocktotal := int(math.Ceil(float64(o.ContentLength) / float64(o.Blocksize)))

	// Create the buffer
	if o.buffer == nil {
		o.buffer = make([][]byte, blocktotal)
	}

	// Calculate which blocks to read from
	firstblock, firstblockstart := divmod(off, o.Blocksize)
	lastblock, lastblockend := divmod(off+bytesrequested-1, o.Blocksize)

	print("First", firstblock, "firstblockstart", firstblockstart, "Last", lastblock, "lastblockend", lastblockend)

	for i := int64(0); i <= lastblock-firstblock; i++ {
		blocknr := firstblock + i
		blocksize := o.Blocksize
		if blocknr == int64(blocktotal) {
			blocksize = lastblockend
		}

		start := int64(0)
		end := blocksize
		if blocknr == firstblock {
			start = off - o.Blocksize*blocknr
		}

		if o.buffer[blocknr] == nil {
			// Create block
			o.buffer[blocknr] = make([]byte, blocksize)

			print("Filled block", blocknr, "with data starting from", blocknr*o.Blocksize, "blocksize", blocksize)

			// Fill block
			// TODO: If block next to each other are empty, retrieve them in one go
			_, err := o.ReaderAt.ReadAt(o.buffer[blocknr], blocknr*o.Blocksize)
			if err != nil && err != io.EOF {
				return n, err
			}
		}

		print("Read block", blocknr, "from", start, "to", end, "block", string(o.buffer[blocknr]))

		// Append the data to out
		out = append(out, o.buffer[blocknr][start:end]...)
	}

	for i, x := range o.buffer {
		print("Block", i, ":", x)
	}
	print("")

	n = copy(p, out)

	return n, err
}

// RangeNotSupported is returned when the remote http server does not support
// the http Range header
type RangeNotSupported struct {
	Url string
}

func (e *RangeNotSupported) Error() string {
	return "Range header not supported by " + e.Url
}

// Logger logs debug messages if set
var Logger *log.Logger

func print(x ...interface{}) {
	if Logger != nil {
		log.Println(x...)
	}
}

func divmod(numerator, denominator int64) (quotient, remainder int64) {
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}

// Compile-time check to see if interfaces are implemented correctly
var _ io.ReaderAt = (*ResponseBody)(nil)
var _ io.Reader = (*ResponseBody)(nil)
var _ io.Closer = (*ResponseBody)(nil)
var _ io.Seeker = (*ResponseBody)(nil)

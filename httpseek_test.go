// Copyright (c) 2022 Hans van Leeuwen. MIT Licensed. See README.md for full license.

package httpseek

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"testing"
)

func Example() {
	client := new(Client)
	resp, err := client.Get("http://textfiles.com/100/phrack.01.phk")
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 33)
	x, err := resp.Body.Seek(555, io.SeekStart)
	y, err := resp.Body.Read(buf)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Read %d bytes starting from position %d. Output: %s", y, x, buf)
	// Output:
	// Read 33 bytes starting from position 555. Output: Welcome to the Phrack Inc. Philes
}

func TestReaderAtBuffer(t *testing.T) {

	b := bytes.NewReader([]byte("0123456789"))
	r := ReaderAtBuffer{Blocksize: 3, ContentLength: 10, ReaderAt: b}
	var p []byte
	var err error

	// Enable logger
	//Logger = log.New(os.Stdout, "", 0)

	p = make([]byte, 2)
	_, err = r.ReadAt(p, 1)
	handleError(err)
	assert(string(p), "12")

	p = make([]byte, 2)
	_, err = r.ReadAt(p, 3)
	handleError(err)
	assert(string(p), "34")

	p = make([]byte, 5)
	_, err = r.ReadAt(p, 5)
	handleError(err)
	assert(string(p), "56789")

}

func assert(a, b interface{}) {
	if a != b {
		log.Fatalln(a, "does not equal ", b)
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

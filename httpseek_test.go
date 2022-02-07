// Copyright (c) 2022 Hans van Leeuwen. MIT Licensed. See README.md for full license.

package httpseek_test

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hnz/httpseek"
)

func Example() {
	client := new(httpseek.Client)
	client.Blocksize = 64 * 1024
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

func ExampleClient() {

	client := new(httpseek.Client)
	client.Blocksize = 64 * 1024
	resp, err := client.Get("https://climate.onebuilding.org/WMO_Region_6_Europe/LUX_Luxembourg/LUX_LU_Luxembourg.AP.065900_TMYx.2004-2018.zip")
	if err != nil {
		panic(err)
	}

	archive, err := zip.NewReader(&resp.Body, resp.ContentLength)
	if err != nil {
		panic(err)
	}

	for _, f := range archive.File {
		fmt.Println("Found file", f.Name)
	}
	// Output:
	// Found file LUX_LU_Luxembourg.AP.065900_TMYx.2004-2018.clm
	// Found file LUX_LU_Luxembourg.AP.065900_TMYx.2004-2018.ddy
	// Found file LUX_LU_Luxembourg.AP.065900_TMYx.2004-2018.epw
	// Found file LUX_LU_Luxembourg.AP.065900_TMYx.2004-2018.rain
	// Found file LUX_LU_Luxembourg.AP.065900_TMYx.2004-2018.stat
	// Found file LUX_LU_Luxembourg.AP.065900_TMYx.2004-2018.wea
}

func TestReaderAtBuffer(t *testing.T) {

	b := bytes.NewReader([]byte("0123456789"))
	r := httpseek.ReaderAtBuffer{Blocksize: 3, ContentLength: 10, ReaderAt: b}
	var p []byte
	var err error

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
	if err != io.EOF {
		handleError(err)
	}
	assert(string(p), "56789")
}

func TestClient(t *testing.T) {

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b := bytes.NewReader([]byte("0123456789"))
		http.ServeContent(w, r, "", time.Now(), b)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client := new(httpseek.Client)
	res, err := client.Get(server.URL)
	handleError(err)
	r := res.Body

	var p []byte

	// Enable logger
	//httpseek.Logger = log.New(os.Stdout, "", 0)

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
	if err != io.EOF {
		handleError(err)
	}
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

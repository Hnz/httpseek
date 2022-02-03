// Copyright (c) 2022 Hans van Leeuwen. MIT Licensed. See README.md for full license.

package httpseek_test

import (
	"fmt"
	"io"

	"github.com/hnz/httpseek"
)

// Compile-time check of interface implementations.
var _ io.Reader = (*httpseek.ResponseBody)(nil)
var _ io.Closer = (*httpseek.ResponseBody)(nil)
var _ io.Seeker = (*httpseek.ResponseBody)(nil)

func Example() {
	client := &httpseek.Client{}
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

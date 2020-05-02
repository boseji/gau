// ####
// ##
// ## Boseji's Safe IO Utilities
// ##
// ## SPDX-License-Identifier: GPL-2.0-only
// ##
// ## Copyright (C) 2020 Abhijit Bose <boseji@users.noreply.github.com>
// ##
// ####

// Package sioutil implements safer alternatives to standard `io.ioutils`
package sioutil

import (
	"bytes"
	"fmt"
	"io"
)

const (
	// Version of the Package
	Version = 0.1
)

// ReadAll safely reads from 'ir' until an error or EOF and returns the data it read.
// A successful call returns err == nil, not err == EOF. Because ReadAll is
// defined to read from src until EOF, it does not treat an EOF from Read
// as an error to be reported. It returns a 'string' type containing
// the data read back and the length of data.
func ReadAll(ir io.Reader) (int64, string, error) {
	// Avoid Panic due to nil Reader
	if ir == nil {
		return 0, "", fmt.Errorf("Error (sioutil.ReadAll): Reader was nil")
	}
	// Create a Buffer for reading data
	// b := bytes.NewBufferString("")
	var b bytes.Buffer // Reduce Allocation Time

	// Perform IO Copy - bytes.Buffer implements Writer interface
	// len, err := io.Copy(b, ir)
	len, err := io.Copy(&b, ir)
	if err != nil {
		return len, "", err
	}
	// Perform Conversion from bytes.Buffer to String
	return len, b.String(), err
}

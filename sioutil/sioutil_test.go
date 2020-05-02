// ####
// ##
// ## Boseji's Safe IO Utilities
// ##
// ## SPDX-License-Identifier: GPL-2.0-only
// ##
// ## Copyright (C) 2020 Abhijit Bose <boseji@users.noreply.github.com>
// ##
// ####

package sioutil_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

// zeroErrReader implements the Reader interface
// to produce error when a read attempt is made.
type zeroErrReader struct {
	err error
}

// Read is the method for Reader Interface
func (r zeroErrReader) Read(p []byte) (int, error) {
	return copy(p, []byte{0}), r.err
}

func TestReadAll(t *testing.T) {
	type args struct {
		ir io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		want1   string
		wantErr bool
	}{
		{
			"Positive Test Case",
			args{strings.NewReader("Hari Aum")},
			8,
			"Hari Aum",
			false,
		},
		{
			"Negative Test Case 1",
			args{nil},
			0,
			"",
			true,
		},
		{
			"Negative Test Case 2",
			args{zeroErrReader{fmt.Errorf("Error No Data")}},
			1,
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ReadAll(tt.args.ir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadAll() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ReadAll() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	fmt.Println("Package Version =", Version)
}

func BenchmarkReadAll512(b *testing.B) {
	bs := bytes.Repeat([]byte{0}, 512+1)
	rd := bytes.NewReader(bs)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ReadAll(rd)
		rd.Reset(bs)
	}
}
func BenchmarkIoutilReadAll512(b *testing.B) {
	bs := bytes.Repeat([]byte{0}, 512+1)
	rd := bytes.NewReader(bs)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ioutil.ReadAll(rd)
		rd.Reset(bs)
	}
}

func BenchmarkReadAllRand(b *testing.B) {
	bs := bytes.Repeat([]byte{0}, b.N+1)
	rd := bytes.NewReader(bs)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ReadAll(rd)
		rd.Reset(bs)
	}
}

func BenchmarkIoutilReadAllRand(b *testing.B) {
	bs := bytes.Repeat([]byte{0}, b.N+1)
	rd := bytes.NewReader(bs)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ioutil.ReadAll(rd)
		rd.Reset(bs)
	}
}

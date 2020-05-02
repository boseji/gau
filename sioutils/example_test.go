// ####
// ##
// ## Boseji's Safe IO Utilities
// ##
// ## SPDX-License-Identifier: GPL-2.0-only
// ##
// ## Copyright (C) 2020 Abhijit Bose <boseji@users.noreply.github.com>
// ##
// ####

package sioutils_test

import (
	"fmt"
	"log"
	"strings"
)

func ExampleReadAll() {

	r := strings.NewReader("Go is a general-purpose language designed with systems programming in mind.")

	len, s, err := sioutil.ReadAll(r)

	if err != nil {

		log.Fatal(err)

	}

	fmt.Printf("%d:%s", len, s)

	// Output:
	// 75:Go is a general-purpose language designed with systems programming in mind.

}

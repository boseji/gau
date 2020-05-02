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
	"fmt"
	"log"
	"strings"

	"github.com/boseji/gau/sioutil"
)

func ExampleReadAll() {

	r := strings.NewReader("Go is a general-purpose language designed with systems programming in mind.")

	len, s, err := sioutil.ReadAll(r)

	if err != nil {

		log.Fatal(err)

	}

	fmt.Printf("%d:%s", len, s)

	// Output:

	//44:Go is a general-purpose language designed with systems programming in mind.

}

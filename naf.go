// Copyright (c) 2018 Ihor Horun

// N.A.F.
// New Age Fidonet tosser & reader

package main

import (
	"fmt"
	"naf/ftnlib"
)

type ftnMessage struct {
}

func main() {
	fmt.Println("New Age Fidonet Tosser & Editor, v.0.1.1")

	/*
		Check address parsing
	*/
	addr, _ := ftnlib.ParseFtnAddress("2:5020/0@fidonet.Org")
	fmt.Println(addr.Dump3D())
	fmt.Println(addr.Dump5D())
	fmt.Println(addr)

	/*
		Check nodelist parsing
	*/
	ftnlib.ParseNodelist("/Users/snake/Projects/ftn-components/src/snakemkua/FTNNodelistBundle/Tests/Resources/nodelist.082")
	//fmt.Println(nodelist)
}

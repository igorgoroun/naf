// Copyright (c) 2018 Ihor Horun

// N.A.F.
// New Age Fidonet tosser & reader

package main

import (
	"fmt"
	"naf/ftnlib"
	"os"
)

func main() {

	app, err := ftnlib.InitApp()

	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	/*
	   Check address parsing
	*/
	//addr, _ := ftnlib.ParseFtnAddress("2:46/12@fidonet.Org")
	//fmt.Println(addr.Dump3D())
	//fmt.Println(addr.Dump5D())
	//fmt.Println(addr)

	/*
	   Check nodelist parsing
	*/
	/*
	   ndl, _ := ftnlib.ParseNodelist("/Users/snake/Projects/ftn-components/src/snakemkua/FTNNodelistBundle/Tests/Resources/nodelist.082")
	   res, err := ndl.UpdateIndex(app)
	   if !res {
	       if err != nil {
	           log.Fatal("Nodelist error", err)
	       } else {
	           log.Fatal("Undefined nodelist error")
	       }
	   }
	*/
	os.Exit(app.Run(os.Args, os.Stdout))
}

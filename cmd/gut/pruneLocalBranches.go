package main

import (
	"fmt"
	"os"

	"github.com/codecare/gut/internal/gut"
)

func main() {
	if len(os.Args) < 2 { //&& os.Args[1] == "prune" {
		fmt.Println("usage: \ngut /path/to/gitdir/")
		os.Exit(1)
	}

	var baseDir = os.Args[1]
	fmt.Println(baseDir)

	branches, err := gut.ReadBranches(baseDir)
	if err != nil {
		panic(err)
	}
	gut.PrintBranches(branches)
}

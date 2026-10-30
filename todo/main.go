package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(os.Args[1:], "data", os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "錯誤：", err)
		os.Exit(1)
	}
}

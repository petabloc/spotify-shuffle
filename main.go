package main

import (
	"fmt"
	"os"

	"github.com/user/spotify-shuffle/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
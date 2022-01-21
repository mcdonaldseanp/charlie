package main

import (
	"github.com/McdonaldSeanp/charlie/auth"
	"os"
	"fmt"
)

func main() {
	switch os.Args[1] {
	case "load":
		switch os.Args[2] {
			case "yubikey":
				err := auth.ConnectYubikey()
				if err != nil {
					fmt.Printf("Did not load yubikey!\n%s", err)
					os.Exit(1)
				}
			default:
				fmt.Printf("Unknown noun!\n")
				os.Exit(1)
			}
	default:
		fmt.Printf("Unknown command!\n")
		os.Exit(1)
	}
	os.Exit(0)
}
package main

import (
	"github.com/McdonaldSeanp/charlie/auth"
	"github.com/McdonaldSeanp/charlie/githelpers"
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
		case "set":
			switch os.Args[2] {
				case "branch":
					err := githelpers.Setgitbranch(os.Args[3])
					if err != nil {
						fmt.Printf("Did not set branch!\n%s", err)
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
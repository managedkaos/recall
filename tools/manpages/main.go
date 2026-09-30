// Command manpages generates recall's manual in ./man.
package main

import (
	"log"

	"github.com/managedkaos/recall/cmd"
)

func main() {
	if err := cmd.GenerateMan("man"); err != nil {
		log.Fatal(err)
	}
}

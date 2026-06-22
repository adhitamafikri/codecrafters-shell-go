package main

import (
	"github.com/codecrafters-io/shell-starter-go/services/shell"
)

func main() {
	sh := shell.NewShell()
	sh.Bootstrap()
}

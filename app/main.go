package main

import (
	"fmt"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func printError(command string) {
	fmt.Printf("%s: command not found", command)
}

func main() {
	// TODO: Uncomment the code below to pass the first stage
	fmt.Print("$ ")

	var command = ""
	fmt.Scan(&command)
	printError(command)
}

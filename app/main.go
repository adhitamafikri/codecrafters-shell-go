package main

import (
	"fmt"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

var commandRegistry = map[string]string{
	"xyz":  "xyz",
	"echo": "echo",
	"cd":   "cd",
	"exit": "exit",
}

func GetCommand(input string) (string, error) {
	val, ok := commandRegistry[input]
	if !ok {
		return "", fmt.Errorf("%s: command not found", input)
	}

	return val, nil
}

func main() {
	// REPL
	for {
		fmt.Print("$ ")

		var command = ""
		fmt.Scan(&command)

		cmd, err := GetCommand(command)
		if err != nil {
			fmt.Print(err, "\n")

		} else {
			fmt.Print(cmd, "\n")
		}
	}
}

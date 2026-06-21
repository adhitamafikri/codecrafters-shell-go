package main

import (
	"fmt"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

var commands = map[string]string{
	"xyz":  "xyz",
	"echo": "echo",
	"cd":   "cd",
}

func GetCommand(inputCommand string) (string, error) {
	val, ok := commands[inputCommand]
	if !ok {
		return "", fmt.Errorf("%s: command not found", inputCommand)
	}

	return val, nil
}

func main() {
	fmt.Print("$ ")

	var command = ""
	fmt.Scan(&command)

	cmd, err := GetCommand(command)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Print(cmd)
}

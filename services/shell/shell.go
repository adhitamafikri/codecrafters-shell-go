package shell

import (
	"fmt"
	"os"
	"syscall"
)

var commandRegistry = map[string]string{
	"xyz":  "xyz",
	"echo": "echo",
	"cd":   "cd",
	"exit": "exit",
}

type shell struct{}

type Shell interface {
	Bootstrap()
	getCommand() (string, error)
	handleCommand(command string)
	handleGracefulShutdown()
}

func NewShell() *shell {
	return &shell{}
}

func (s *shell) Bootstrap() {
	for {
		fmt.Print("$ ")

		var command = ""
		fmt.Scan(&command)

		cmd, err := s.getCommand(command)
		if err != nil {
			fmt.Print(err, "\n")
		} else {
			s.commandHandler(cmd)
		}
	}
}

func (s *shell) getCommand(input string) (string, error) {
	val, ok := commandRegistry[input]
	if !ok {
		return "", fmt.Errorf("%s: command not found", input)
	}

	return val, nil
}

func (s *shell) commandHandler(command string) {
	switch command {
	case commandRegistry["exit"]:
		s.handleGracefulShutdown()
	default:
		fmt.Println("Unimplemented...")
	}
}

func (s *shell) handleGracefulShutdown() {
	// Find current running process
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		panic(err)
	}

	// fmt.Println("Bye...")

	// Emit Termination Signal
	err = p.Signal(syscall.SIGTERM)
	if err != nil {
		fmt.Println("Failed to send termination signal.", err)
	}
}

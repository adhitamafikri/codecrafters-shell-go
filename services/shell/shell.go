package shell

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var commandRegistry = map[string]string{
	"xyz":  "xyz",
	"echo": "echo",
	"cd":   "cd",
	"exit": "exit",
}

type shell struct {
	reader *bufio.Reader
}

type Shell interface {
	Bootstrap()
	getInput() string
	getCommand(input string) (string, error)
	getArgs(input string) []string
	handleInput(input string)
	handleEcho(args []string)
	handleGracefulShutdown()
}

func NewShell() *shell {
	return &shell{
		reader: bufio.NewReader(os.Stdin),
	}
}

func (s *shell) Bootstrap() {
	for {
		input, err := s.getInput()
		if err != nil {
			fmt.Print(err, "\n")
		} else {
			err = s.handleInput(input)

			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (s *shell) getInput() (string, error) {
	fmt.Print("$ ")
	input, err := s.reader.ReadString('\n')

	if err != nil {
		return "", fmt.Errorf("Something went wrong with the reader\n")
	}

	return input, nil
}

func (s *shell) getCommand(input string) (string, error) {
	splitResult := strings.Split(strings.TrimSpace(input), " ")
	command := strings.TrimSpace(splitResult[0])

	val, ok := commandRegistry[command]
	if !ok {
		return "", fmt.Errorf("%s: command not found", command)
	}

	return val, nil
}

func (s *shell) getArgs(input string) []string {
	splitResult := strings.Split(strings.TrimSpace(input), " ")

	var args []string
	if len(splitResult) > 1 {
		args = splitResult[1:]
	}

	return args
}

// Note: Orchestrator function, responsible for: extracting command from input, extracting values from input, routing to appropriate handler
func (s *shell) handleInput(input string) error {
	command, err := s.getCommand(input)

	if err != nil {
		return err
	}

	args := s.getArgs(input)

	switch command {
	case commandRegistry["echo"]:
		s.handleEcho(args)
	case commandRegistry["exit"]:
		s.handleGracefulShutdown()
	default:
		fmt.Println("Unimplemented...")
	}

	return nil
}

func (s *shell) handleEcho(args []string) {
	length := len(args)

	for idx, arg := range args {
		if idx != length-1 {
			fmt.Printf("%s ", arg)
		} else {
			fmt.Print(arg)
		}
	}

	fmt.Print("\n")
}

func (s *shell) handleGracefulShutdown() {
	// Find current running process
	// p, err := os.FindProcess(os.Getpid())
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Bye...")

	// Emit Termination Signal
	// err = p.Signal(syscall.SIGTERM)
	// if err != nil {
	// 	fmt.Println("Failed to send termination signal.", err)
	// }

	// Straightforward way to exit the process
	os.Exit(0)
}

package shell

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var commandRegistry = map[string]map[string]string{
	"builtin": {
		"xyz":  "xyz",
		"type": "type",
		"echo": "echo",
		"cd":   "cd",
		"exit": "exit",
	},
}

type shell struct {
	reader *bufio.Reader
}

type Shell interface {
	Bootstrap()
	getInput() string
	getCommandAndArgs(input string) (string, []string, error)
	handleInputBuiltin(input string)
	handleTypeBuiltin(args []string)
	handleNonBuiltin(args []string)
	handleEchoBuiltin(args []string)
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
			trimmedInput := strings.TrimSpace(input)

			if len(trimmedInput) > 0 {
				err = s.handleInputBuiltin(trimmedInput)

				if err != nil {
					fmt.Println(err)
				}
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

func (s *shell) getCommandAndArgs(input string) (string, []string, error) {
	splitResult := strings.Split(input, " ")
	command := splitResult[0]

	cmd, ok := commandRegistry["builtin"][command]
	if !ok {
		return "", []string{}, fmt.Errorf("%s: command not found", command)
	}

	var args []string
	if len(splitResult) > 1 {
		args = splitResult[1:]
	}

	return cmd, args, nil
}

// Note: Orchestrator function, responsible for: extracting command from input, extracting values from input, routing to appropriate handler
func (s *shell) handleInputBuiltin(input string) error {
	cmd, args, err := s.getCommandAndArgs(input)

	if err != nil {
		return err
	}

	switch cmd {
	case commandRegistry["builtin"]["type"]:
		s.handleTypeBuiltin(args)
	case commandRegistry["builtin"]["echo"]:
		s.handleEchoBuiltin(args)
	case commandRegistry["builtin"]["exit"]:
		s.handleGracefulShutdown()
	default:
		// fmt.Println("Unimplemented...")
		s.handleNonBuiltin(args)
	}

	return nil
}

func (s *shell) handleTypeBuiltin(args []string) {
	cmd, _, err := s.getCommandAndArgs(args[0])
	if err != nil {
		s.handleNonBuiltin(args)
		return
	}

	fmt.Printf("%s is a shell builtin\n", cmd)
}

func (s *shell) handleNonBuiltin(args []string) {
	arg := args[0]
	path, err := exec.LookPath(arg)

	if err != nil {
		fmt.Printf("%s: not found\n", arg)
		return
	}

	fmt.Printf("%s is %s\n", arg, path)

}

func (s *shell) handleEchoBuiltin(args []string) {
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

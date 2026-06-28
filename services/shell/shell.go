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
	getCommandAndArgs(input string) (cmd string, args []string, isBuiltin bool)
	getExecutablePath(arg string) (string, error)
	handleInput(input string)
	handleTypeBuiltin(args []string)
	handleNonBuiltin(cmd string, args []string)
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
				err = s.handleInput(trimmedInput)

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

func (s *shell) getCommandAndArgs(input string) (cmd string, args []string, isBuiltin bool) {
	splitResult := strings.Split(input, " ")
	command := splitResult[0]

	cmd, ok := commandRegistry["builtin"][command]
	if ok {
		isBuiltin = true
	} else {
		cmd = command
		isBuiltin = false
	}

	if len(splitResult) > 1 {
		args = splitResult[1:]
	}

	return cmd, args, isBuiltin
}

func (s *shell) getExecutablePath(arg string) (string, error) {
	path, err := exec.LookPath(arg)

	if err != nil {
		return "", fmt.Errorf("%s: not found", arg)
	}

	return path, nil
}

// Note: Orchestrator function, responsible for: extracting command from input, extracting values from input, routing to appropriate handler
func (s *shell) handleInput(input string) error {
	cmd, args, isBuiltin := s.getCommandAndArgs(input)

	if isBuiltin {
		switch cmd {
		case commandRegistry["builtin"]["type"]:
			s.handleTypeBuiltin(args)
		case commandRegistry["builtin"]["echo"]:
			s.handleEchoBuiltin(args)
		case commandRegistry["builtin"]["exit"]:
			s.handleGracefulShutdown()
		default:
			fmt.Println("Unimplemented...")
		}
	} else {
		s.handleNonBuiltin(cmd, args)
	}

	return nil
}

func (s *shell) handleTypeBuiltin(args []string) {
	if len(args) == 0 {
		fmt.Println("Not enough arguments provided for type builtin")
		return
	}

	cmd, _, isBuiltin := s.getCommandAndArgs(args[0])
	if !isBuiltin {
		path, err := s.getExecutablePath(args[0])
		if err != nil {
			fmt.Printf("%s: not found\n", args[0])
			return
		}

		fmt.Printf("%s is %s\n", args[0], path)
		return
	}

	fmt.Printf("%s is a shell builtin\n", cmd)
}

// This method is responsible for executing a program, just like what we usually do on our terminal
// Example: 'hx ~/.claude/config', 'cowsay "hey" | lolcat', 'ls -lah', .etc
func (s *shell) handleNonBuiltin(cmd string, args []string) {
	_, err := s.getExecutablePath(cmd)

	if err != nil {
		fmt.Printf("%s: not found\n", cmd)
		return
	}

	// Define the command, execute, and capture the output
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		fmt.Printf("Failed to execute %s command\n", cmd)
	} else {
		fmt.Printf("$ %s\n", out)
	}

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

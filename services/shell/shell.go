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
		"pwd":  "pwd",
		"exit": "exit",
	},
}

type shell struct {
	reader *bufio.Reader
}

type Shell interface {
	Bootstrap()
	getInput() string
	parseShellInput(input string) (cmd string, args []string, isBuiltin bool)
	getExecutablePath(arg string) (string, error)
	handleInput(input string)
	handleTypeBuiltin(args []string)
	handleNonBuiltin(cmd string, args []string)
	handleEchoBuiltin(args []string)
	handlePwdBuiltin() string
	handleCdBuiltin(args []string)
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

func (s *shell) parseShellInput(input string) (cmd string, args []string, isBuiltin bool) {
	spaceIndex := strings.Index(input, " ")
	if spaceIndex == -1 {
		_, ok := commandRegistry["builtin"][input]
		return input, []string{}, ok
	}

	cmd = input[0:spaceIndex]
	args = strings.Split(input[spaceIndex+1:], " ")

	// Arg rules:
	// 1. Any args having ONLY single quote must error
	// 2. Args without single quotes must treated as trimmed values, no trailing, leading, or empty whitespaces spaces, except whitespace given from the end-printing logic
	// 3. Any args having single quotes without any value between them must be removed, the rest of args must be concatenated as is
	// 4. Any args having single quotes with any value between them must be retained as is. No missing spaces, signs, characters, .etc

	fmt.Println("the raw args:", args)
	length := len(args)

	foundPrefixSq, foundSuffixSq := false, false
	tempQuotedArg := ""
	var bucket []string

	// Example args: 'hello     this world' belongs''to me     'tanja' i love'tanja'much
	// Expected result: hello      this world belongsto me tanja i lovetanjamuch

	for i := 0; i < length; i++ {
		// Check has prefix and suffix sq
		foundPrefixSq = strings.HasPrefix(args[i], "'")
		foundSuffixSq = strings.HasSuffix(args[i], "'")

		if len(tempQuotedArg) == 0 {
			if foundPrefixSq && !foundSuffixSq {
				tempQuotedArg = args[i][1:]
				bucket = append(bucket, tempQuotedArg)
			} else if !foundPrefixSq && foundSuffixSq {
				lArg := len(args[i])
				bucket = append(bucket, args[i][0:lArg-1])
				tempQuotedArg = ""
			} else if foundPrefixSq && foundSuffixSq {
				lArg := len(args[i])
				bucket = append(bucket, args[i][1:lArg-1])
			} else {
				// push to bucket as long as it's not a whitespace character
				if args[i] != "" {
					bucket = append(bucket, args[i])
				}
			}
		} else {
			// Things to do when we have an opening value of quoted args
			if foundPrefixSq && !foundSuffixSq {
				tempQuotedArg = args[i][1:]
				tempQuotedArg = ""
				bucket = append(bucket, tempQuotedArg)
			} else if !foundPrefixSq && foundSuffixSq {
				lArg := len(args[i])
				tempQuotedArg = ""
				bucket = append(bucket, args[i][0:lArg-1])
			} else if foundPrefixSq && foundSuffixSq {
				lArg := len(args[i])
				bucket = append(bucket, args[i][1:lArg-1])
			} else {
				bucket = append(bucket, args[i])
			}
		}
	}

	var result []string
	for _, item := range bucket {
		result = append(result, strings.ReplaceAll(item, "'", ""))
	}

	_, ok := commandRegistry["builtin"][cmd]
	isBuiltin = ok

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
	cmd, args, isBuiltin := s.parseShellInput(input)

	if isBuiltin {
		switch cmd {
		case commandRegistry["builtin"]["type"]:
			s.handleTypeBuiltin(args)
		case commandRegistry["builtin"]["echo"]:
			s.handleEchoBuiltin(args)
		case commandRegistry["builtin"]["pwd"]:
			s.handlePwdBuiltin()
		case commandRegistry["builtin"]["cd"]:
			s.handleCdBuiltin(args)
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

	cmd, _, isBuiltin := s.parseShellInput(args[0])
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
		fmt.Printf("%s", out)
	}

}

func (s *shell) handleEchoBuiltin(args []string) {
	result := ""
	isOpening := false

	// sanitizing the args first
	var bucket []string
	for _, item := range args {
		hasPrefixSq := strings.HasPrefix(item, "'")
		hasSuffixSq := strings.HasSuffix(item, "'")

		if hasPrefixSq && hasSuffixSq {
			// result += strings.ReplaceAll(item, "'", "")
			bucket = append(bucket, strings.ReplaceAll(item, "'", ""))
		} else if hasPrefixSq && !hasSuffixSq {
			isOpening = true
			// result += strings.ReplaceAll(item, "'", "")
			bucket = append(bucket, strings.ReplaceAll(item, "'", ""))
		} else if !hasPrefixSq && hasSuffixSq {
			isOpening = false
			bucket = append(bucket, strings.ReplaceAll(item, "'", ""))
		} else {
			if item == "" {
				if isOpening {
					bucket = append(bucket, item)
				}
			} else {
				bucket = append(bucket, strings.ReplaceAll(item, "'", ""))
			}
		}
	}

	// Concatenating
	for _, item := range bucket {
		result += item + " "
	}
	fmt.Printf("%s\n", strings.TrimSuffix(result, " "))
}

func (s *shell) handlePwdBuiltin() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get current working directory")
		return
	}

	fmt.Println(dir)
}

func (s *shell) handleCdBuiltin(args []string) {
	if len(args) == 0 {
		fmt.Println("Not enough arguments provided")
		return
	}

	dir := args[0]

	if dir == "~" {
		d, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Failed to get home directory")
			return
		}
		dir = d
	}

	err := os.Chdir(dir)
	if err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", dir)
	}
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

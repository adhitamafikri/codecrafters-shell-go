package shell

import (
	"reflect"
	"testing"
)

func TestParseShellInput(t *testing.T) {
	s := NewShell()
	tests := []struct {
		input string
		cmd   string
		args  []string
	}{
		{
			input: "cat '/tmp/cow/f   37' '/tmp/cow/f   87' '/tmp/cow/f   71'",
			cmd:   "cat",
			args:  []string{"/tmp/cow/f   37", "/tmp/cow/f   87", "/tmp/cow/f   71"},
		},
		{input: "echo Playboi''Carti", cmd: "echo", args: []string{"PlayboiCarti"}},
		{input: "echo Playboi'Car'ti", cmd: "echo", args: []string{"PlayboiCarti"}},
		{input: "echo ''", cmd: "echo", args: []string{""}},
		{input: "echo 'hello     world'", cmd: "echo", args: []string{"hello     world"}},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			cmd, args, _, err := s.parseShellInput(test.input)
			if err != nil {
				t.Fatalf("parseShellInput() error = %v", err)
			}
			if cmd != test.cmd || !reflect.DeepEqual(args, test.args) {
				t.Errorf("parseShellInput() = (%q, %#v), want (%q, %#v)", cmd, args, test.cmd, test.args)
			}
		})
	}
}

func TestParseShellInputRejectsUnclosedSingleQuote(t *testing.T) {
	_, _, _, err := NewShell().parseShellInput("echo 'hello")
	if err == nil {
		t.Fatal("parseShellInput() error = nil, want an error")
	}
}

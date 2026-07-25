package processors

import (
	"strings"
)

// Returns command and args from the given input
func ParseShellInput(input string) (string, []string) {
	spaceIndex := strings.Index(input, " ")
	command := input[0:spaceIndex]
	rawArgs := input[spaceIndex:]
	args := strings.Split(rawArgs, " ")

	// Arg rules:
	// 1. Any args having ONLY single quote must error
	// 2. Args without single quotes must treated as trimmed values, no trailing, leading, or empty whitespaces spaces, except whitespace given from the end-printing logic
	// 3. Any args having single quotes without any value between them must be removed, the rest of args must be concatenated as is
	// 4. Any args having single quotes with any value between them must be retained as is. No missing spaces, signs, characters, .etc

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

	return command, result
}

func ProcessEchoArgs(args []string) string {
	// Arg rules:
	// 1. Any args having ONLY single quote must error
	// 2. Args without single quotes must treated as trimmed values, no trailing, leading, or empty whitespaces spaces, except whitespace given from the end-printing logic
	// 3. Any args having single quotes without any value between them must be removed, the rest of args must be concatenated as is
	// 4. Any args having single quotes with any value between them must be retained as is. No missing spaces, signs, characters, .etc

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

	result := ""
	lBucket := len(bucket)
	for idx, item := range bucket {
		if idx < lBucket {
			result += strings.ReplaceAll(item, "'", "") + " "
		} else {
			result += strings.ReplaceAll(item, "'", "")
		}
	}

	return result
}

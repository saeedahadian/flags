package flags

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var flagKeyRegex = regexp.MustCompile(`-{1,2}(.+)`)

type Flags []*Flag

func (flags *Flags) Append(flag *Flag) {
	*flags = append(*flags, flag)
}

func Parse() (*Flags, error) {
	args := getFlagArgs()
	var flags *Flags
	for i := 0; i < len(args); i++ {
		arg := args[i]

		matches := flagKeyRegex.FindStringSubmatch(arg)

		if matches == nil {
			continue
		}

		flagKey := matches[1]

		if !hasNextIndex(i, args) {
			flags.Append(MustNewBoolFlag(flagKey, "true"))
			break
		}

		nextArg := args[i+1]
		nextMatches := flagKeyRegex.FindStringSubmatch(nextArg)

		if nextMatches != nil {
			flags.Append(MustNewBoolFlag(flagKey, "true"))
			continue
		}

		switch getStringType(nextArg) {
		case "int":
			flags.Append(MustNewIntFlag(flagKey, nextArg))
		case "bool":
			flags.Append(MustNewBoolFlag(flagKey, nextArg))
		case "string":
			flags.Append(MustNewStringFlag(flagKey, nextArg))
		}
		i++
	}
	return flags, nil
}


func getFlagArgs() []string {
	startIdx, endIdx := -1, -1
	for i, arg := range os.Args[1:] {
		if startIdx == -1 && flagKeyRegex.MatchString(arg) {
			startIdx = i
			continue
		}

		if endIdx == -1 && arg == "--" {
			endIdx = i
			break
		}
	}
	
	if startIdx == -1 {
		return nil
	}

	if endIdx == -1 {
		return os.Args[startIdx:]
	}

	return os.Args[startIdx:endIdx]
}

func hasNextIndex[T any](currIdx int, slice []T) bool {
	return currIdx < (len(slice) - 1)
}

func getStringType(s string) string {
	if _, err := strconv.Atoi(s); err == nil {
		return "int"
	}

	if _, err := strconv.ParseBool(s); err == nil {
		return "bool"
	}

	return "string"
}

func ParseFlags(args []string) ([]*Flag, error) {
	var flags []*Flag

	if len(args) == 0 {
		return flags, nil
	}

	for i := 1; i < len(args); i++ {
		arg := args[i]
		prevArg := args[i-1]
		_, isFlagKey := parseKey(arg)
		prevArgKey, isPrevFlagKey := parseKey(prevArg)
		if isPrevFlagKey {
			if isFlagKey {
				flag, err := NewBoolFlag(prevArgKey, "true")
				if err != nil {
					return nil, fmt.Errorf("key %q: %w", prevArgKey, err)
				}
				flags = append(flags, flag)
				continue
			}

			if _, err := strconv.Atoi(arg); err == nil {
				flag, err := NewIntFlag(prevArgKey, arg)
				if err != nil {
					return nil, fmt.Errorf("key %q: %w", prevArgKey, err)
				}
				flags = append(flags, flag)
				continue
			}

			if _, err := strconv.ParseBool(arg); err == nil {
				flag, err := NewBoolFlag(prevArgKey, arg)
				if err != nil {
					return nil, fmt.Errorf("key %q: %w", prevArgKey, err)
				}
				flags = append(flags, flag)
				continue
			}

			flag, err := NewStringFlag(prevArgKey, arg)
			if err != nil {
				return nil, fmt.Errorf("key %q: %w", prevArgKey, err)
			}
			flags = append(flags, flag)
		}
	}

	lastArg := args[len(args)-1]
	lastArgKey, isLastArgKey := parseKey(lastArg)
	if !isLastArgKey {
		return flags, nil
	}
	
	flag, err := NewBoolFlag(lastArgKey, "true")
	if err != nil {
		return nil, fmt.Errorf("key %q: %w", lastArgKey, err)
	}
	flags = append(flags, flag)

	return flags, nil
}

func parseKey(arg string) (string, bool) {
	matches := flagKeyRegex.FindStringSubmatch(arg)
	if len(matches) < 2 {
		return "", false
	}
	return matches[1], true
}

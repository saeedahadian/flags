package flags

import (
	"fmt"
	"regexp"
	"strconv"
)

var flagKeyRegex = regexp.MustCompile(`-{1,2}(.+)`)

type FlagValue interface {
	String() string
	Set(string) error
}

type intValue int

func (i *intValue) String() string {
	return strconv.Itoa(int(*i))
}

func (i *intValue) Set(value string) error {
	val, err := strconv.Atoi(value)
	if err != nil {
		return err
	}

	*i = intValue(val)
	return nil
}

type stringValue string

func (s *stringValue) Set(value string) error {
	*s = stringValue(value)
	return nil
}

func (s *stringValue) String() string {
	return string(*s)
}

type boolValue bool

func (b *boolValue) Set(value string) error {
	val, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}

	*b = boolValue(val)
	return nil
}

func (b *boolValue) String() string {
	return strconv.FormatBool(bool(*b))
}

func ParseFlags(args []string) ([]*Flag, error) {
	var flags []*Flag
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
	return flags, nil
}

func parseKey(arg string) (string, bool) {
	matches := flagKeyRegex.FindStringSubmatch(arg)
	if len(matches) < 2 {
		return "", false
	}
	return matches[1], true
}
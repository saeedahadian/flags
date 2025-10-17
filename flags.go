package flags

import (
	"fmt"
	"os"
	"strconv"
)

type Flag struct {
	Key   string
	Value FlagValue
}

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


func NewIntFlag(key string, value string) (*Flag, error) {
	var val intValue
	if err := val.Set(value); err != nil {
		return nil, fmt.Errorf("flag %q: %w", key, err)
	}
	return &Flag{Key: key, Value: &val}, nil
}

func MustNewIntFlag(key string, value string) *Flag {
	flag, err := NewIntFlag(key, value)
	if err != nil {
		printKeyValueErrAndFail(key, value, err)
	}
	return flag
}

func NewBoolFlag(key string, value string) (*Flag, error) {
	var val boolValue
	if err := val.Set(value); err != nil {
		return nil, fmt.Errorf("flag %q: %w", key, err)
	}
	return &Flag{Key: key, Value: &val}, nil
}

func MustNewBoolFlag(key string, value string) *Flag {
	flag, err := NewBoolFlag(key, value)
	if err != nil {
		printKeyValueErrAndFail(key, value, err)
	}
	return flag
}

func NewStringFlag(key string, value string) (*Flag, error) {
	var val stringValue
	if err := val.Set(value); err != nil {
		return nil, fmt.Errorf("flag %q: %w", key, err)
	}
	return &Flag{Key: key, Value: &val}, nil
}

func MustNewStringFlag(key string, value string) *Flag {
	flag, err := NewStringFlag(key, value)
	if err != nil {
		printKeyValueErrAndFail(key, value, err)
	}
	return flag
}

func printKeyValueErrAndFail(key string, value string, err error) {
	fmt.Printf("key %q, value %q: %s", key, value, err.Error())
	os.Exit(1)
}
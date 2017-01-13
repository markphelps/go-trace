package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

type boundIntValue struct {
	val      *int
	min, max int
}

func newBoundIntValue(val int, p *int, min, max int) *boundIntValue {
	*p = val
	return &boundIntValue{p, min, max}
}

func (i *boundIntValue) String() string {
	if i.val == nil {
		return ""
	}
	return fmt.Sprintf("%d", *i.val)
}

func (i *boundIntValue) Set(value string) error {
	v, err := strconv.Atoi(value)
	if err != nil {
		return err
	}

	if v > i.max {
		v = i.max
	} else if v < i.min {
		v = i.min
	}
	*i.val = v
	return nil
}

// BoundIntVar defines an int flag with specified name, default value, minimum value, maximum value and usage string.
// The argument p points to an int variable in which to store the value of the flag.
func BoundIntVar(p *int, name string, value, min, max int, usage string) {
	flag.Var(newBoundIntValue(value, p, min, max), name, usage)
}

// BoundInt defines an int flag with specified name, default value, minimum value, maximum value and usage string.
// The return value is the address of an int variable that stores the value of the flag.
func BoundInt(name string, value, min, max int, usage string) *int {
	p := &value
	BoundIntVar(p, name, value, min, max, usage)
	return p
}

type boundFloat64Value struct {
	val      *float64
	min, max float64
}

func newBoundFloat64Value(val float64, p *float64, min, max float64) *boundFloat64Value {
	*p = val
	return &boundFloat64Value{p, min, max}
}

func (f *boundFloat64Value) String() string {
	if f.val == nil {
		return ""
	}
	return fmt.Sprintf("%f", *f.val)
}

func (f *boundFloat64Value) Set(value string) error {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}

	if v > f.max {
		v = f.max
	} else if v < f.min {
		v = f.min
	}

	*f.val = v
	return nil
}

// BoundFloat64Var defines a float64 flag with specified name, default value, minimum value, maximum value and usage string.
// The argument p points to a float64 variable in which to store the value of the flag.
func BoundFloat64Var(p *float64, name string, value, min, max float64, usage string) {
	flag.Var(newBoundFloat64Value(value, p, min, max), name, usage)
}

// BoundFloat64 defines a float64 flag with specified name, default value, minimum value, maximum value and usage string.
// The return value is the address of an float64 variable that stores the value of the flag.
func BoundFloat64(name string, value, min, max float64, usage string) *float64 {
	p := &value
	BoundFloat64Var(p, name, value, min, max, usage)
	return p
}

type filenameValue struct {
	val        *string
	extensions []string
}

func newFilenameValue(val string, p *string, extensions []string) *filenameValue {
	*p = val
	return &filenameValue{p, extensions}
}

func (f *filenameValue) String() string {
	if f.val == nil {
		return ""
	}
	return *f.val
}

func (f *filenameValue) Set(value string) error {
	var err error

	allowed := strings.Join(f.extensions, ",")
	switch ext := strings.ToLower(filepath.Ext(value)); ext {
	case allowed:
		*f.val = value
	default:
		err = fmt.Errorf("Invalid extension: %s, allowed: %s", ext, allowed)
	}
	return err
}

// FilenameVar defines a string flag with specified name, default value, allowed values and usage string.
// The argument p points to a string variable in which to store the value of the flag.
func FilenameVar(p *string, name, value string, extensions []string, usage string) {
	flag.Var(newFilenameValue(value, p, extensions), name, usage)
}

// Filename defines a string flag with specified name, default value, allowed values and usage string.
// The return value is the address of an string variable that stores the value of the flag.
func Filename(name, value string, extensions []string, usage string) *string {
	p := &value
	FilenameVar(p, name, value, extensions, usage)
	return p
}

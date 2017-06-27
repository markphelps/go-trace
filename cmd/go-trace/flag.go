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

type filenameValue struct {
	val        *string
	extensions map[string]interface{}
}

func newFilenameValue(val string, p *string, extensions map[string]interface{}) *filenameValue {
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
	ext := strings.ToLower(filepath.Ext(value))

	if _, ok := f.extensions[ext]; !ok {
		return fmt.Errorf("Invalid extension: %s", ext)
	}

	*f.val = value
	return nil
}

// FilenameVar defines a string flag with specified name, default value, allowed extension values and usage string.
// The argument p points to a string variable in which to store the value of the flag.
func FilenameVar(p *string, name, value string, extensions map[string]interface{}, usage string) {
	flag.Var(newFilenameValue(value, p, extensions), name, usage)
}

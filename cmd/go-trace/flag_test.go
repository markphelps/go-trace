package main

import "testing"

func TestNewBoundIntValue(t *testing.T) {
	p := new(int)

	expected := 1
	min := 0
	max := 10

	i := newBoundIntValue(expected, p, min, max)

	if *p != expected {
		t.Errorf("Invalid reference, expected: %v, got: %v", expected, *p)
	}
	if *i.val != expected {
		t.Errorf("Invalid val, expected: %v, got: %v", expected, *i.val)
	}
	if i.min != min {
		t.Errorf("Invalid min, expected: %v, got: %v", min, i.min)
	}
	if i.max != max {
		t.Errorf("Invalid max, expected: %v, got: %v", max, i.max)
	}
}

func TestBoundIntValueString(t *testing.T) {
	testCases := []struct {
		val      *int
		expected string
	}{
		{val: new(int), expected: "0"},
		{val: nil, expected: ""},
	}

	for _, tc := range testCases {
		i := boundIntValue{val: tc.val}
		if i.String() != tc.expected {
			t.Errorf("Invalid String() value, expected: %v, got: %v", tc.expected, i.String())
		}
	}
}

func TestBoundIntValueSet(t *testing.T) {
	testCases := []struct {
		value    string
		expected int
		min, max int
		err      bool
	}{
		{value: "0", expected: 0, min: 0, max: 10},
		{value: "5", expected: 5, min: 0, max: 10},
		{value: "10", expected: 10, min: 0, max: 10},
		{value: "100", expected: 10, min: 0, max: 10},
		{value: "-1", expected: 0, min: 0, max: 10},
		{value: "foo", expected: 0, min: 0, max: 10, err: true},
	}

	for _, tc := range testCases {
		i := boundIntValue{val: new(int), min: tc.min, max: tc.max}
		err := i.Set(tc.value)

		if err != nil && tc.err != true {
			t.Error("Expected parse error, got none")
		}

		if *i.val != tc.expected {
			t.Errorf("Invalid value, expected: %v, got: %v", tc.expected, *i.val)
		}
	}
}

func TestNewBoundFloat64Value(t *testing.T) {
	p := new(float64)

	expected := 1.0
	min := 0.0
	max := 10.0

	i := newBoundFloat64Value(expected, p, min, max)

	if *p != expected {
		t.Errorf("Invalid reference, expected: %v, got: %v", expected, *p)
	}
	if *i.val != expected {
		t.Errorf("Invalid val, expected: %v, got: %v", expected, *i.val)
	}
	if i.min != min {
		t.Errorf("Invalid min, expected: %v, got: %v", min, i.min)
	}
	if i.max != max {
		t.Errorf("Invalid max, expected: %v, got: %v", max, i.max)
	}
}

func TestBoundFloat64ValueString(t *testing.T) {
	testCases := []struct {
		val      *float64
		expected string
	}{
		{val: new(float64), expected: "0.000000"},
		{val: nil, expected: ""},
	}

	for _, tc := range testCases {
		i := boundFloat64Value{val: tc.val}
		if i.String() != tc.expected {
			t.Errorf("Invalid String() value, expected: %v, got: %v", tc.expected, i.String())
		}
	}
}

func TestBoundFloat64ValueSet(t *testing.T) {
	testCases := []struct {
		value    string
		expected float64
		min, max float64
		err      bool
	}{
		{value: "0.0", expected: 0.0, min: 0, max: 10},
		{value: "5.1", expected: 5.1, min: 0, max: 10},
		{value: "10.0", expected: 10, min: 0, max: 10},
		{value: "100.0", expected: 10, min: 0, max: 10},
		{value: "-1.0", expected: 0, min: 0, max: 10},
		{value: "foo", expected: 0, min: 0, max: 10, err: true},
	}

	for _, tc := range testCases {
		i := boundFloat64Value{val: new(float64), min: tc.min, max: tc.max}
		err := i.Set(tc.value)

		if err != nil && tc.err != true {
			t.Error("Expected parse error, got none")
		}

		if *i.val != tc.expected {
			t.Errorf("Invalid value, expected: %v, got: %v", tc.expected, *i.val)
		}
	}
}

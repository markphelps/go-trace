package main

import "testing"

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

		if err != nil && !tc.err {
			t.Error("Expected parse error, got none")
		}

		if *i.val != tc.expected {
			t.Errorf("Invalid value, expected: %v, got: %v", tc.expected, *i.val)
		}
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

		if err != nil && !tc.err {
			t.Error("Expected parse error, got none")
		}

		if *i.val != tc.expected {
			t.Errorf("Invalid value, expected: %v, got: %v", tc.expected, *i.val)
		}
	}
}

func TestFilenameValueString(t *testing.T) {
	testCases := []struct {
		val      *string
		expected string
	}{
		{val: new(string), expected: ""},
		{val: nil, expected: ""},
	}

	for _, tc := range testCases {
		i := filenameValue{val: tc.val}
		if i.String() != tc.expected {
			t.Errorf("Invalid String() value, expected: %v, got: %v", tc.expected, i.String())
		}
	}
}

func TestFilenameValueSet(t *testing.T) {
	testCases := []struct {
		value   string
		allowed map[string]interface{}
		valid   bool
		err     string
	}{
		{value: "test.png", allowed: map[string]interface{}{".png": true}, valid: true},
		{value: "test.invalid", allowed: map[string]interface{}{".png": true}, valid: false, err: "Invalid extension: .invalid"},
		{value: "test.png", allowed: map[string]interface{}{".png": true, ".jpg": true}, valid: true},
		{value: "test.jpg", allowed: map[string]interface{}{".png": true, ".jpg": true}, valid: true},
	}

	for _, tc := range testCases {
		f := filenameValue{val: new(string), extensions: tc.allowed}
		err := f.Set(tc.value)

		if !tc.valid {
			if err == nil {
				t.Error("Expected error, got none")
			} else if err.Error() != tc.err {
				t.Errorf("Expected error: %v, got: %v", tc.err, err.Error())
			}
		} else {
			if *f.val != tc.value {
				t.Errorf("Expected: %v, got: %v", tc.value, *f.val)
			}
		}
	}
}

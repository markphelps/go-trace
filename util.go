package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

func writePNG(path string, img image.Image) {
	file, err := os.Create(path)
	check(err, "Error opening file: %v\n")
	defer file.Close()

	err = png.Encode(file, img)
	check(err, "Error writing to file: %v\n")
}

func check(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, msg, err)
		os.Exit(1)
	}
}

func boundFloat(min, max, value float64) float64 {
	if value > max {
		value = max
	} else if value < min {
		value = min
	}
	return value
}

func boundInt(min, max, value int) int {
	if value > max {
		value = max
	} else if value < min {
		value = min
	}
	return value
}

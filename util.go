package main

import (
	"image"
	"image/png"
	"os"
)

func writePNG(path string, img image.Image) (err error) {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	err = png.Encode(file, img)
	return err
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

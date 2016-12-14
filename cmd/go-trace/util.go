package main

import (
	"image"
	"image/png"
	"os"
)

func writePNG(path string, img image.Image) error {
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

func clampFloat(value, min, max float64) float64 {
	if value > max {
		value = max
	} else if value < min {
		value = min
	}
	return value
}

func clampInt(value, min, max int) int {
	if value > max {
		value = max
	} else if value < min {
		value = min
	}
	return value
}

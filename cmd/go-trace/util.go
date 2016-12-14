package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"strings"
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

func outputProgress(ch <-chan int, rows int) {
	fmt.Println()
	for i := 1; i <= rows; i++ {
		<-ch
		pct := 100 * float64(i) / float64(rows)
		filled := (progressBarWidth * i) / rows
		bar := strings.Repeat("=", filled) + strings.Repeat("-", progressBarWidth-filled)
		fmt.Printf("\r[%s] %.2f%%", bar, pct)
	}
	fmt.Println()
}

package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	acceptableBrightnessDifference = 125
	acceptableColorDifference      = 500
)

type color struct {
	Red   int
	Green int
	Blue  int
}

func (c *color) Set(input string) error {
	input = strings.TrimLeft(input, "#")
	if len(input) != 6 {
		return fmt.Errorf("invalid hex value: %s", input)
	}

	redBytes, err := hex.DecodeString(input[:2])
	if err != nil {
		return fmt.Errorf("error getting red value from hex: %s", input[:2])
	} else if len(redBytes) > 1 {
		return fmt.Errorf("hex value too large: %s", redBytes)
	}

	greenBytes, err := hex.DecodeString(input[2:4])
	if err != nil {
		return fmt.Errorf("error getting green value from hex: %s", input[2:4])
	} else if len(greenBytes) > 1 {
		return fmt.Errorf("hex value too large: %s", greenBytes)
	}

	blueBytes, err := hex.DecodeString(input[4:])
	if err != nil {
		return fmt.Errorf("error getting blue value from hex: %s", input[4:])
	} else if len(blueBytes) > 1 {
		return fmt.Errorf("hex value too large: %s", blueBytes)
	}

	c.Red = int(redBytes[0])
	if c.Red > 255 {
		return fmt.Errorf("red value too large: %d", c.Red)
	}

	c.Green = int(greenBytes[0])
	if c.Green > 255 {
		return fmt.Errorf("green value too large: %d", c.Green)
	}

	c.Blue = int(blueBytes[0])
	if c.Blue > 255 {
		return fmt.Errorf("blue value too large: %d", c.Blue)
	}
	return nil
}

func main() {
	var hex1, hex2 string
	if len(os.Args) < 3 {
		fmt.Printf("Enter the hex value for the first color: #")
		fmt.Scan(&hex1)
		fmt.Printf("Enter the hex value for the second color: #")
		fmt.Scan(&hex2)
	} else {
		hex1, hex2 = os.Args[1], os.Args[2]
	}

	color1 := &color{}
	err := color1.Set(hex1)
	if err != nil {
		log.Fatalf("Error with first hex: %v\n", err)
	}

	color2 := &color{}
	err = color2.Set(hex2)
	if err != nil {
		log.Fatalf("Error with second hex: %v\n", err)
	}

	bd := brightnessDifference(*color1, *color2)
	if bd < acceptableBrightnessDifference {
		fmt.Printf("Calculated brightness difference (%d) is lower than acceptable brightness difference (%d).\n",
			bd, acceptableBrightnessDifference)
		os.Exit(0)
	}

	cd := colorDifference(*color1, *color2)
	if cd < acceptableColorDifference {
		fmt.Printf("Calculated color difference (%d) is lower than acceptable color difference (%d).\n",
			cd, acceptableColorDifference)
		os.Exit(0)
	}

	fmt.Printf("Brightness: %d and color difference: %d OK!\n", bd, cd)
}

func brightnessDifference(color1, color2 color) int {
	diff := brightness(color1) - brightness(color2)
	if diff < 0 {
		return -diff
	}
	return diff
}

func brightness(c color) int {
	return ((c.Red * 299) + (c.Green * 587) + (c.Blue * 114)) / 1000
}

func colorDifference(color1, color2 color) int {
	return (max(color1.Red, color2.Red) - min(color1.Red, color2.Red)) +
		(max(color1.Green, color2.Green) - min(color1.Green, color2.Green)) +
		(max(color1.Blue, color2.Blue) - min(color1.Blue, color2.Blue))
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

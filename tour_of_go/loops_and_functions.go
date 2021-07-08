package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {
	z := 1.0
	newZ := (z - (z*z-x)/(2*z))
	for math.Abs(z-newZ) > 0.001 {
		z = newZ
		newZ = z - (z*z-x)/(2*z)
	}
	return z
}

func main() {
	fmt.Println(Sqrt(2))
}

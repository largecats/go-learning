package main

import (
	"fmt"
	"math"
)

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number: %v", float64(e))
	/* fmt.Sprint(e) will call e.Error() to convert the value e to a string. If the Error() method calls fmt.Sprint(e),
	then the program recurses until out of memory. You can break the recursion by converting the e to a value without a String or Error method. */
}

func Sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, ErrNegativeSqrt(x)
	}

	z := 1.0
	newZ := (z - (z*z-x)/(2*z))
	for math.Abs(z-newZ) > 0.001 {
		z = newZ
		newZ = z - (z*z-x)/(2*z)
	}
	return z, nil
}

func main() {
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(-2))
}

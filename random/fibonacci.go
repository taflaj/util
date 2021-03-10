// random/fibonacci.go

package random

import (
	"fmt"
)

var fib map[int]int

func init() {
	fib = make(map[int]int)
	fib[0] = 0
	fib[1] = 1
}

// Fibonacci returns the Fibonacci number with the given index
func Fibonacci(index int) (int, error) {
	if index < 0 {
		return index, fmt.Errorf("Unable to calculate the Fibonacci number for a negative index such as %v", index)
	}
	result, ok := fib[index]
	if ok {
		return result, nil
	}
	previous1, err := Fibonacci(index - 1)
	if err != nil {
		return 0, err
	}
	previous2, err := Fibonacci(index - 2)
	if err != nil {
		return 0, err
	}
	result = previous1 + previous2
	fib[index] = result
	return result, nil
}

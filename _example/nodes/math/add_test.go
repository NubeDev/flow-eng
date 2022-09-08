package math

import (
	"fmt"
	"testing"
)

func Test_math(t *testing.T) {
	fmt.Println(mathOperation(add, 100, 10, 10, 100))
	fmt.Println(mathOperation(sub, 100, 10, 10, 100))
	fmt.Println(mathOperation(multiply, 2, 2, 2, 20))
}

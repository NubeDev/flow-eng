package array

import (
	"fmt"
	"testing"
)

func TestAllTrueFloat64(t *testing.T) {

	fmt.Println("Subtract", Subtract([]float64{100, 10, 10}))
	fmt.Println("Add", Add([]float64{10, 10, 10, 1}))
	fmt.Println("Divide", Divide([]float64{10, 2, 2}))
	fmt.Println("Multiply", Multiply([]float64{2, 10, 2}))

	fmt.Println("AllTrueFloat64", AllTrueFloat64([]float64{1.1, 1, 1}))
	fmt.Println("OneIsTrueFloat6", OneIsTrueFloat64([]float64{0, 0, 1.1}))
	fmt.Println("MaxFloat64", MaxFloat64([]float64{111, 0, 1.1}))
	fmt.Println("MinFloat64", MinFloat64([]float64{111, 0, -111.11}))
	min, max := MinMaxFloat64([]float64{})
	fmt.Println("MinMaxFloat64", min, max)
	min, max = MinMaxFloat64([]float64{100, 0, -50})
	fmt.Println("MinMaxFloat64", min, max)
}

package main

import "fmt"

type Number struct {
	value interface{}
}

func add(n1, n2 Number) Number {
	switch v1 := n1.value.(type) {
	case int:
		if v2, ok := n2.value.(int); ok {
			return Number{v1 + v2}
		}
	case float64:
		if v2, ok := n2.value.(float64); ok {
			return Number{v1 + v2}
		}
	}
	panic("Unsupported types for addition")
}

func handleFailure() {
	if r := recover(); r != nil {
		fmt.Println("Recovered from panic:", r)
	}
}
func main2() {
	fmt.Println("Hello, World!")

	n1 := Number{10}
	n2 := Number{20}
	resultInt := add(n1, n2)
	defer func() {
		//if r := recover(); r != nil {
		//	fmt.Println("Recovered from panic:", r)
		//}
		handleFailure()
	}()
	fmt.Printf("Addition of integers: %v\n", resultInt.value)

	n3 := Number{15}
	n4 := Number{24.5}
	resultFloat := add(n3, n4)
	fmt.Printf("Addition of floats: %v\n", resultFloat.value)
}

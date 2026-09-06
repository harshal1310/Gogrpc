package generics

import (
	"cmp"
)

// golang 1.18 to 1.20
type vars interface {
	int | float64 | string
}

func add[T vars](a, b T) T {
	return a + b
}

func Add[T cmp.Ordered](a, b T) T {
	return a + b
}

//func main() {
//	fmt.Println(add(5.2, 3))
//	fmt.Println(Add(5, 3))
//}

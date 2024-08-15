package main

import (
	"fmt"
	"reflect"
)

func main() {

	a := [3]int{1, 2, 3}

	b := a[0:2:2]

	fmt.Println(reflect.TypeOf(b))

	b = append(b, 5)
	fmt.Println(b)

	fmt.Println(a)

}

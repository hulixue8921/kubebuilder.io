package main

import (
	"fmt"
)

type X interface {
	run()
}

type P struct {
	Name string
	Have []string
}

func (p P) run() {
	fmt.Println("run")
}

func T(x X) {

}

func main() {
	//	x, _ := os.ReadFile("./1.go")
	//	fmt.Println(string(x))

	x := `
	xxxx
	"cjvljdlj"
	cjvlljlj
	xxcjvldfjdf`

	fmt.Println(x)
}

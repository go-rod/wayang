package wayang_test

import (
	"fmt"

	"github.com/ysmood/kit"

	"github.com/go-rod/wayang"
)

//
func Example_RunProgram() {
	var program wayang.Program

	// Reading the json to program struct
	err := kit.ReadJSON("examples_test.json", &program)
	if err != nil {
		println(err.Error())
	}

	// last item executed is returned
	res, err := wayang.RunProgram(program)
	if err != nil {
		println(err.Error())
	}

	fmt.Println(res)

	// Output:
	// package main
	//
	// import (
	// 	"fmt"
	// 	"time"
	// )
	//
	// var c chan int
	//
	// func handle(int) {}
	//
	// func main() {
	// 	select {
	// 	case m := <-c:
	// 		handle(m)
	// 	case <-time.After(10 * time.Second):
	// 		fmt.Println("timed out")
	// 	}
	// }
}

package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	fmt.Println(os.Getenv("USERNAME"))
	for i := 1; i < len(os.Args); i++ {
		fmt.Println(os.Args[i])
	}
	fmt.Println(runtime.Version())
}
package main

import (
	"fmt"
	"os"
	"runtime"
)

func test() {
	wd, err := os.Getwd()
	fmt.Println(wd, err)

	exePath, err := os.Executable()
	fmt.Println(exePath, err)

	pc, file, line, ok := runtime.Caller(0)
	fmt.Println(pc, file, line, ok)

	root, ok := findProjectRoot(file)
	if !ok {
		fmt.Println("Failed to find project root")
		return
	}
	fmt.Println("Project root:", root)
}

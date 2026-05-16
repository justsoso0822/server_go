package main

import (
	"fmt"
	"strings"
)

func main() {
	str := "1001,1,1002,0"
	parts := strings.SplitN(str, ",", -1)
	fmt.Println(parts)
}
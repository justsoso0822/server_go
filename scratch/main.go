package main

import (
	"fmt"
	"time"
)

func main() {
	// str := "1001,1,1002,0"
	// parts := strings.SplitN(str, ",", -1)
	// fmt.Println(parts)

	tm := time.Now().Add(-time.Second * 120).Truncate(time.Second)
	fmt.Println(tm)
}

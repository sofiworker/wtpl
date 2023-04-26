package main

import (
	"fmt"
	"wtpl/cmd"
)

func main() {
	err := cmd.Run()
	if err != nil {
		fmt.Printf("run program failed:%s", err.Error())
	}
}

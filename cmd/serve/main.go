package main

import (
	"easycoding/cmd/serve/cmd"
	"fmt"
	"os"
)

func main() {
	if err := cmd.InitCmd(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

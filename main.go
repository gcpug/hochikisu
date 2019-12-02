package main

import (
	"fmt"
	"github.com/gcpug/hochikisu/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Printf("err %+v", err)
		os.Exit(1)
	}
	os.Exit(0)
}

package main

import (
	"os"

	pool "nacos-sdk-go-example/pkg/name-pool"
)

func main() {
	if err := pool.Execute(); err != nil {
		os.Exit(1)
	}
}

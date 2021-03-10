package main

import (
	//"distributed/coordinator"
	"fmt"

	"github.com/gopheramit/distributed-go-with-rabbitmq/src/distributed/coordinator"
)

func main() {
	ql := coordinator.NewQueueListener()
	go ql.ListenForNewSource()

	var a string
	fmt.Scanln(&a)
}

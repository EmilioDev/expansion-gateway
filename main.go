package main

import (
	"fmt"
)

func main() {
	fmt.Println("game gateway starting...")

	gateway := GetGateway()

	gateway.Start()
}

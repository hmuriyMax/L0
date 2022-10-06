package main

import (
	"github.com/hmuriyMax/L0/pkg/random_order"
	"log"
)

func main() {
	err := random_order.SendRandomOrder(log.Default())
	if err != nil {
		log.Fatal(err)
	}
}

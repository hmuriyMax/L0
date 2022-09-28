package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hmuriyMax/L0/cmd/internal/database"
	"github.com/hmuriyMax/L0/cmd/internal/order_receiver"
	"log"
	"os"
)

func insertTest(db *database.DataBase) {
	bytes, err := os.ReadFile("./model.json")
	if err != nil {
		log.Fatal(err)
	}
	var order order_receiver.Order
	err = json.Unmarshal(bytes, &order)
	if err != nil {
		log.Fatal(err)
	}
	db.Insert(context.TODO(), order)
}

func main() {
	db := database.New(log.Default())
	err := db.Start()
	defer db.Stop()
	if err != nil {
		log.Fatal(err)
	}
	for {
		var command string
		_, err := fmt.Fscanln(os.Stdin, &command)
		if command == "stop" || err != nil {
			break
		}
	}
}

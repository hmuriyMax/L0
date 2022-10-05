package main

import (
	"context"
	"fmt"
	"github.com/hmuriyMax/L0/internal/database"
	"github.com/hmuriyMax/L0/internal/order_receiver"
	"log"
	"os"
	"time"
)

//func insertTest(db *database.DataBase) {
//	bytes, err := os.ReadFile("./model.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//	var order order_receiver.Order
//	err = json.Unmarshal(bytes, &order)
//	if err != nil {
//		log.Fatal(err)
//	}
//	db.Insert(context.background(), order)
//}

func main() {
	newOrders := make(chan order_receiver.Order)
	db := database.New(log.Default())
	err := db.Start()
	defer db.Stop()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	conn, subs := order_receiver.Run(newOrders)
	go func() {
		for {
			select {
			case order, ok := <-newOrders:
				if !ok {
					return
				}
				db.Insert(ctx, order)
			}
		}

	}()
	for {
		var command string
		_, err := fmt.Fscanln(os.Stdin, &command)
		if command == "stop" || err != nil {
			_ = conn.Close()
			_ = subs.Unsubscribe()
			break
		}
	}
}

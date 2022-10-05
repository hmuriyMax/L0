package main

import (
	"encoding/json"
	"github.com/hmuriyMax/L0/internal/order_receiver"
	"github.com/nats-io/stan.go"
	"log"
	"math/rand"
	"time"
)

func GetJSON(order order_receiver.Order) (bytes []byte) {
	bytes, err := json.Marshal(order)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func randomString(length int) string {
	const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	res := make([]byte, length)
	for i := range res {
		res[i] = chars[rand.Intn(len(chars))]
	}
	return string(res)
}

func main() {
	conn, err := stan.Connect(order_receiver.ClusterID,
		"order-client",
		stan.NatsURL(stan.DefaultNatsURL),
		stan.ConnectWait(time.Second*5))
	if err != nil {
		log.Fatalln(err)
	}
	defer func() { _ = conn.Close() }()

	message := GetJSON(order_receiver.Order{
		OrderUid:    randomString(rand.Intn(5) + 5),
		TrackNumber: randomString(rand.Intn(5) + 5),
		Entry:       randomString(4),
		Delivery: order_receiver.OrderDelivery{
			Name:  randomString(10),
			Phone: randomString(10),
		},
		Payment: order_receiver.OrderPayment{
			Transaction: randomString(rand.Intn(5) + 5),
		},
		Items: []order_receiver.OrderItem{
			{
				ChartID:     rand.Int63(),
				TrackNumber: randomString(rand.Intn(5) + 5),
			},
		},
		Locale:            randomString(2),
		InternalSignature: randomString(rand.Intn(3)),
		CustomerID:        randomString(rand.Intn(5) + 3),
		DeliveryService:   randomString(10),
		ShardKey:          randomString(rand.Intn(5) + 5),
		SmID:              rand.Intn(100),
		DateCreated:       time.Now().Format(time.RFC3339),
		OofShard:          randomString(1),
	})

	err = conn.Publish(order_receiver.Channel, message)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Published message on channel: " + order_receiver.Channel)
}

package order_receiver

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
)

const (
	ClusterID = "test-cluster"
	ClientID  = "order-service"
	Channel   = "id-channel"
)

type OrderDelivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type OrderPayment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDT    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type OrderItem struct {
	ChartID     int64  `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int64  `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

type Order struct {
	OrderUid          string        `json:"order_uid"`
	TrackNumber       string        `json:"track_number"`
	Entry             string        `json:"entry"`
	Delivery          OrderDelivery `json:"delivery"`
	Payment           OrderPayment  `json:"payment"`
	Items             []OrderItem   `json:"items"`
	Locale            string        `json:"locale"`
	InternalSignature string        `json:"internal_signature"`
	CustomerID        string        `json:"customer_id"`
	DeliveryService   string        `json:"delivery_service"`
	ShardKey          string        `json:"shardkey"`
	SmID              int           `json:"sm_id"`
	DateCreated       string        `json:"date_created"`
	OofShard          string        `json:"oof_shard"`
}

func Run(orders chan Order) (stan.Conn, stan.Subscription, error) {
	connect, err := stan.Connect(ClusterID, ClientID, stan.NatsURL(stan.DefaultNatsURL))
	if err != nil {
		return nil, nil, err
	}
	subs, err := connect.Subscribe(Channel, func(msg *stan.Msg) {
		var order Order
		log.Printf("received message from: %s", msg.Subject)
		err = json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Println(err)
			return
		}
		orders <- order
	}, stan.StartWithLastReceived())
	if err != nil {
		return nil, nil, err
	}
	return connect, subs, nil
}

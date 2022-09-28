package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/hmuriyMax/L0/cmd/internal/order_receiver"
	"log"
)

const (
	connection = "user=maxim password=fuck2022 port=5432 database=l0 sslmode=disable"
)

var dbase *sql.DB

func init() {
	var err error
	dbase, err = sql.Open("database", connection)
	if err != nil || dbase.Ping() != nil {
		log.Fatal(err)
	}
}

func Insert(ctx context.Context, order order_receiver.Order) {
	command := fmt.Sprintf(""+"INSERT "+
		"INTO orders (order_uid, track_number, "+
		"entry, delivery, payment, items, locale, internal_signature, customer_id, "+
		"delivery_service, shardkey, sm_id, date_created, oof_shard)"+
		"VALUES ('%s', '%s', '%s', '%vs', '%vs', '%vs', '%s', '%s', '%s', '%s', '%s', '%d', '%s', '%s')",
		order.OrderUid, order.TrackNumber, order.Entry, order.Delivery, order.Payment, order.Items,
		order.Locale, order.InternalSignature, order.CustomerId, order.DeliveryService, order.Shardkey,
		order.SmID, order.DateCreated, order.OofShard)
	_, err := dbase.ExecContext(ctx, command)
	if err != nil {
		log.Println(err)
	}
}

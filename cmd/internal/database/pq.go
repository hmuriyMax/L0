package database

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hmuriyMax/L0/cmd/internal/order_receiver"
	"log"
)

func (db *DataBase) importToCache(ctx context.Context) {
	//for i := 0; i < 60; i++ {
	//	fmt.Println(i)
	//	time.Sleep(time.Second)
	//}
	rows, err := db.pg.QueryContext(ctx, "SELECT * FROM orders")
	if err != nil {
		db.lg.Println(err)
	}
	for rows.Next() {
		var (
			tmp      order_receiver.Order
			delivery []byte
			payment  []byte
			items    []byte
		)

		err := rows.Scan(&tmp.OrderUid, &tmp.TrackNumber, &tmp.Entry, &delivery, &payment, &items,
			&tmp.Locale, &tmp.InternalSignature, &tmp.CustomerId, &tmp.DeliveryService, &tmp.Shardkey,
			&tmp.SmID, &tmp.DateCreated, &tmp.OofShard)
		if err != nil {
			db.lg.Println(err)
		}
		_ = json.Unmarshal(delivery, &tmp.Delivery)
		_ = json.Unmarshal(payment, &tmp.Payment)
		_ = json.Unmarshal(items, &tmp.Items)

		db.cache.insert(tmp, db.lg)
	}
	log.Println("import finished successfully!")
}

func (db *DataBase) insertInDB(ctx context.Context, order order_receiver.Order) {
	delivery, _ := json.Marshal(order.Delivery)
	payment, _ := json.Marshal(order.Payment)
	items, _ := json.Marshal(order.Items)
	command := fmt.Sprintf(""+"INSERT "+
		"INTO orders (order_uid, track_number, "+
		"entry, delivery, payment, items, locale, internal_signature, customer_id, "+
		"delivery_service, shardkey, sm_id, date_created, oof_shard)"+
		"VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%d', '%s', '%s')",
		order.OrderUid, order.TrackNumber, order.Entry, delivery, payment, items,
		order.Locale, order.InternalSignature, order.CustomerId, order.DeliveryService, order.Shardkey,
		order.SmID, order.DateCreated, order.OofShard)
	_, err := db.pg.ExecContext(ctx, command)
	if err != nil {
		log.Println(err)
	}
}

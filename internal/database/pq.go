package database

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hmuriyMax/L0/internal/order_receiver"
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
			&tmp.Locale, &tmp.InternalSignature, &tmp.CustomerID, &tmp.DeliveryService, &tmp.ShardKey,
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

func (db *DataBase) insertInDB(ctx context.Context, order order_receiver.Order) error {
	delivery, err := json.Marshal(order.Delivery)
	if err != nil {
		return err
	}
	payment, err := json.Marshal(order.Payment)
	if err != nil {
		return err
	}
	items, err := json.Marshal(order.Items)
	if err != nil {
		return err
	}
	command := fmt.Sprintf(""+"INSERT "+
		"INTO orders (order_uid, track_number, "+
		"entry, delivery, payment, items, locale, internal_signature, customer_id, "+
		"delivery_service, shardkey, sm_id, date_created, oof_shard)"+
		"VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%d', '%s', '%s')",
		order.OrderUid, order.TrackNumber, order.Entry, delivery, payment, items,
		order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.ShardKey,
		order.SmID, order.DateCreated, order.OofShard)
	_, err = db.pg.ExecContext(ctx, command)
	if err != nil {
		return err
	}
	return nil
}

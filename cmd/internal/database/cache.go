package database

import (
	"fmt"
	"github.com/hmuriyMax/L0/cmd/internal/order_receiver"
	"log"
)

type Cache map[string]order_receiver.Order

func (c Cache) len() int { return len(c) }

func (c Cache) insert(order order_receiver.Order, lg *log.Logger) {
	if c.checkIn(order) {
		lg.Println(fmt.Errorf("order '%s' already exists in cache", order.OrderUid))
		return
	}
	c[order.OrderUid] = order
	lg.Printf("value %s cashed", order.OrderUid)
}

func (c Cache) checkIn(order order_receiver.Order) bool {
	_, ok := c[order.OrderUid]
	return ok
}

func (c Cache) all() (res []order_receiver.Order) {
	for _, val := range c {
		res = append(res, val)
	}
	return
}

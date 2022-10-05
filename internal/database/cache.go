package database

import (
	"fmt"
	"log"
)

type Cache map[string]Order

func (c Cache) len() int { return len(c) }

func (c Cache) insert(order Order, lg *log.Logger) {
	if c.checkIn(order) {
		lg.Println(fmt.Errorf("order '%s' already exists in cache", order.OrderUid))
		return
	}
	c[order.OrderUid] = order
	lg.Printf("value %s cashed", order.OrderUid)
}

func (c Cache) checkIn(order Order) bool {
	_, ok := c[order.OrderUid]
	return ok
}

func (c Cache) all() (res []Order) {
	for _, val := range c {
		res = append(res, val)
	}
	return
}

func (c Cache) getById(id string) (ord Order, ok bool) {
	ord, ok = c[id]
	return
}

func (c Cache) removeById(id string, lg *log.Logger) (ok bool) {
	_, ok = c[id]
	if !ok {
		delete(c, id)
	}
	lg.Println("removed id:", id)
	return
}

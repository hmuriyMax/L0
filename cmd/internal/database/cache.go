package database

import "github.com/hmuriyMax/L0/cmd/internal/order_receiver"

type Cache map[string]order_receiver.Order

func (c Cache) Len() int { return len(c) }

func (c Cache) Import() {

}

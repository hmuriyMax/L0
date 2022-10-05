package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/hmuriyMax/L0/internal/order_receiver"
	_ "github.com/lib/pq"
	"log"
	"sync"
)

const (
	connection = "user=maxim password=fuck2022 port=5432 database=l0 sslmode=disable"
)

type DataBase struct {
	pg    *sql.DB
	cache Cache
	lg    *log.Logger
	wg    sync.WaitGroup
}

func New(lg *log.Logger) *DataBase {
	var d DataBase
	d.lg = lg
	d.wg = sync.WaitGroup{}
	d.cache = make(map[string]order_receiver.Order)
	return &d
}

func (db *DataBase) Start() (err error) {
	db.pg, err = sql.Open("postgres", connection)
	if err != nil || db.pg.Ping() != nil {
		return
	}
	db.wg.Add(1)
	go func() {
		db.importToCache(context.TODO())
		db.wg.Done()
	}()
	db.lg.Println("database started")
	return
}

func (db *DataBase) Insert(ctx context.Context, order order_receiver.Order) {
	if db.cache.checkIn(order) {
		db.lg.Println(fmt.Errorf("order '%s' already exists", order.OrderUid))
		return
	}
	db.cache.insert(order, db.lg)
	db.wg.Add(1)
	go func() {
		db.insertInDB(ctx, order)
		db.wg.Done()
	}()
}

func (db *DataBase) GetAll() []order_receiver.Order {
	return db.cache.all()
}

func (db *DataBase) Stop() {
	db.wg.Wait()
	db.lg.Println("database stopped")
	return
}

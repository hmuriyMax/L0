package database

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"sync"
)

const (
	connection = "user=%s password=%s port=%s database=%s sslmode=disable"
)

type configStruct struct {
	Name     string `yaml:"user"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
}

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
	d.cache = make(map[string]Order)
	return &d
}

func (db *DataBase) Start() (err error) {
	var config configStruct
	configBytes, err := os.ReadFile(configPath)
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return
	}

	db.pg, err = sql.Open("postgres",
		fmt.Sprintf(connection, config.Name, config.Password, config.Port, config.Database))
	if err != nil {
		return
	}
	if err = db.pg.Ping(); err != nil {
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

func (db *DataBase) Insert(ctx context.Context, order Order) {
	if db.cache.checkIn(order) {
		db.lg.Println(fmt.Errorf("order '%s' already exists", order.OrderUid))
		return
	}
	db.cache.insert(order, db.lg)
	db.wg.Add(1)
	go func() {
		defer db.wg.Done()
		err := db.insertInDB(ctx, order)
		if err != nil {
			db.lg.Println(err)
			db.cache.removeById(order.OrderUid, db.lg)
		}
	}()
}

func (db *DataBase) GetAll() []Order {
	return db.cache.all()
}

func (db *DataBase) GetById(id string) (Order, error) {
	ord, ok := db.cache.getById(id)
	if !ok {
		return Order{}, fmt.Errorf("no order with id '%s'", id)
	}
	return ord, nil
}

func (db *DataBase) Stop() {
	db.wg.Wait()
	db.lg.Println("database stopped")
	return
}

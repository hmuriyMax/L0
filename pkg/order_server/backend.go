package order_server

import (
	"context"
	"encoding/json"
	"fmt"
	db "github.com/hmuriyMax/L0/internal/database"
	"github.com/nats-io/stan.go"
	"log"
	"os"
	"time"
)

const (
	DefaultExportFile = "./api/cache_log.json"
)

type App struct {
	db     *db.DataBase
	conn   stan.Conn
	subs   stan.Subscription
	logger *log.Logger

	ctx        context.Context
	cancelFunc func()
}

func (a *App) Start(lg *log.Logger) error {
	a.db = db.New(lg)
	a.logger = lg
	err := a.db.Start()
	a.ctx, a.cancelFunc = context.WithCancel(context.Background())
	return err
}

func (a *App) Stop() {
	a.cancelFunc()
	a.db.Stop()
	_ = a.conn.Close()
	_ = a.subs.Unsubscribe()
	a.logger.Println("app stopped")
}

func (a *App) InitListener() (err error) {
	newOrders := make(chan db.Order)
	a.conn, a.subs, err = a.runNutsListener(newOrders)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case order, ok := <-newOrders:
				if !ok {
					return
				}
				ctx, _ := context.WithTimeout(a.ctx, time.Second*5)
				a.db.Insert(ctx, order)
			case <-a.ctx.Done():
				a.logger.Printf("context is canceled. Finishing newOrders listener")
				return
			}
		}
	}()
	return
}

func (a *App) ExportCache(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	bytes, err := json.MarshalIndent(a.db.GetAll(), "", "\t")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(file, string(bytes))
	if err != nil {
		return err
	}
	a.logger.Printf("Successfully exported cache to %s", filename)
	return nil
}

func (a *App) GetById(id string) ([]byte, error) {
	order, err := a.db.GetById(id)
	if err != nil {
		return nil, err
	}
	return json.Marshal(order)
}

func (a *App) runNutsListener(orders chan db.Order) (stan.Conn, stan.Subscription, error) {
	connect, err := stan.Connect(db.ClusterID, "order-server", stan.NatsURL(stan.DefaultNatsURL))
	if err != nil {
		return nil, nil, err
	}
	subs, err := connect.Subscribe(db.ChannelID, func(msg *stan.Msg) {
		var order db.Order
		a.logger.Printf("received message from: %s", msg.Subject)
		err = json.Unmarshal(msg.Data, &order)
		if err != nil {
			a.logger.Println(err)
			return
		}
		orders <- order
	}, stan.StartWithLastReceived())
	if err != nil {
		return nil, nil, err
	}
	return connect, subs, nil
}

func InsertTest(a App) error {
	bytes, err := os.ReadFile("./api/model.json")
	if err != nil {
		return err
	}
	var order db.Order
	err = json.Unmarshal(bytes, &order)
	if err != nil {
		return err
	}
	a.db.Insert(a.ctx, order)
	return nil
}

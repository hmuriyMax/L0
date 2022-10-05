package back

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hmuriyMax/L0/internal/database"
	"github.com/hmuriyMax/L0/internal/order_receiver"
	"github.com/nats-io/stan.go"
	"log"
	"os"
	"time"
)

const (
	DefaultExportFile = "./cache_log.json"
)

type App struct {
	db     *database.DataBase
	conn   stan.Conn
	subs   stan.Subscription
	logger *log.Logger

	ctx        context.Context
	cancelFunc func()
}

func (a *App) Start(lg *log.Logger) error {
	a.db = database.New(lg)
	err := a.db.Start()
	a.ctx, a.cancelFunc = context.WithCancel(context.Background())
	return err
}

func (a *App) Stop() {
	a.cancelFunc()
	a.db.Stop()
	_ = a.conn.Close()
	_ = a.subs.Unsubscribe()

}

func (a *App) InitListener() (err error) {
	newOrders := make(chan order_receiver.Order)
	a.conn, a.subs, err = order_receiver.Run(newOrders)
	if err != nil {
		return
	}
	go func() {
		var (
			cf  func()
			ctx context.Context
		)
		for {
			select {
			case order, ok := <-newOrders:
				if !ok {
					return
				}
				ctx, _ = context.WithTimeout(a.ctx, time.Second*5)
				a.db.Insert(ctx, order)
			case <-a.ctx.Done():
				log.Printf("context is canceled. Finishing newOrders listener")
				cf()
				return
			}
		}
	}()
	return
}

func (a *App) ExportCached(filename string) error {
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
	return nil
}

func (a *App) GetById(id string) ([]byte, error) {
	order, err := a.db.GetById(id)
	if err != nil {
		return nil, err
	}
	return json.Marshal(order)
}

//func main() {
//	newOrders := make(chan order_receiver.Order)
//	db := database.New(log.Default())
//	err := db.Start()
//	defer db.Stop()
//	if err != nil {
//		log.Fatal(err)
//	}
//	conn, subs := order_receiver.Run(newOrders)
//	go func() {
//		for {
//			select {
//			case order, ok := <-newOrders:
//				if !ok {
//					return
//				}
//				ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
//				db.Insert(ctx, order)
//			}
//		}
//
//	}()
//	for {
//		var command string
//		_, err := fmt.Fscanln(os.Stdin, &command)
//		if command == "stop" || err != nil {
//			_ = conn.Close()
//			_ = subs.Unsubscribe()
//			break
//		} else if command == "cached" {
//
//		}
//	}
//}

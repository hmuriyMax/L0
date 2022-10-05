package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hmuriyMax/L0/internal/database"
	"github.com/hmuriyMax/L0/pkg/order_server"
	"html/template"
	"log"
	"net/http"
	"os"
)

const (
	httpAddress = ":8080"
	HTMLPath    = "./web/template/"
)

var (
	backend = order_server.App{}
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	bytes, err := backend.GetById(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var anm database.Order
	err = json.Unmarshal(bytes, &anm)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageTemlate := template.Must(template.ParseFiles(HTMLPath + "index.html"))

	err = pageTemlate.Execute(w, anm)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	err := backend.Start(logger)
	if err != nil {
		log.Fatal(err)
	}
	err = backend.InitListener()
	if err != nil {
		log.Fatal(err)
	}
	defer backend.Stop()

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	srv := &http.Server{
		Addr:    httpAddress,
		Handler: mux,
	}
	go func() {
		err = srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Fatalln(err)
		}
	}()
	logger.Printf("server listening on http://%s\n", srv.Addr)
	for {
		var command string
		_, err := fmt.Fscanln(os.Stdin, &command)
		if err != nil {
			logger.Fatal(err)
		}
		switch command {
		case "stop":
			_ = srv.Shutdown(context.Background())
			return
		case "export":
			err := backend.ExportCache(order_server.DefaultExportFile)
			if err != nil {
				logger.Println(err)
			}
		case "import":
			err := order_server.InsertTest(backend)
			if err != nil {
				logger.Println(err)
			}
		default:
			logger.Println("Unknown command:", command)
		}

	}
}

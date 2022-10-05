package main

import (
	"encoding/json"
	"fmt"
	"github.com/hmuriyMax/L0/cmd/back"
	"github.com/hmuriyMax/L0/internal/order_receiver"
	"html/template"
	"log"
	"net/http"
	"os"
)

const (
	httpAddress = ":8080"
	HTMLPath    = "./cmd/httpserver/html/"
)

var (
	backend = back.App{}
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	bytes, err := backend.GetById(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var anm order_receiver.Order
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
	err := backend.Start(log.Default())
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
	defer func() { _ = srv.Close() }()
	go func() {
		err = srv.ListenAndServe()
		if err != nil {
			log.Fatalln(err)
		}
	}()
	log.Printf("server listening on http://%s\n", srv.Addr)
	for {
		var command string
		_, err := fmt.Fscanln(os.Stdin, &command)
		if command == "stop" || err != nil {
			return
		} else if command == "cached" {
			err := backend.ExportCached(back.DefaultExportFile)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

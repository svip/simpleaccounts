package main

import (
	"api"
	"db"
	"log"
	"net/http"
)

func main() {
	db.Init(false)
	http.HandleFunc("/account/", api.AccountEntry)
	http.HandleFunc("/transaction/", api.TransactionEntry)
	http.Handle("/", http.FileServer(http.Dir("media/")))
	log.Println("Server listening on 8080.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

package main

import (
	"api"
	"db"
	"net/http"
	"log"
)

func main() {
	db.Init(false)
	http.HandleFunc("/account/", api.AccountEntry)
	http.HandleFunc("/transaction/", api.TransactionEntry)
	log.Println("Server listening on 8080.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

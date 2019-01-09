package api

import (
	"encoding/json"
	"fmt"
	"log"
	"model"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const LOGFILE = "access.log"

func doLog(r *http.Request) {
	s := fmt.Sprintf("%s %s %s\n", time.Now().Format(time.RFC3339), r.Method, r.URL)
	fmt.Print(s)
	file, err := os.OpenFile(LOGFILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Println("Error logging", err)
		return
	}
	defer file.Close()

	file.Write([]byte(s))
}

func getURLPath(r *http.Request, index int) (string, bool) {
	s := strings.Split(r.URL.Path, "/")
	if index >= len(s) {
		return "", false
	}
	return s[index], true
}

func getURLPathInt(r *http.Request, index int) (int, bool) {
	s, ok := getURLPath(r, index)
	if !ok {
		return 0, false
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return i, true
}

func handleRawError(w http.ResponseWriter, statuscode int, message string) {
	w.WriteHeader(statuscode)
	s := struct {
		Error string
	}{
		message,
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(&s); err != nil {
		w.Write([]byte(message))
		return
	}
	w.Header().Add("Content-Type", "application/json")
}

func handleError(w http.ResponseWriter, statuscode int, err error) {
	handleRawError(w, statuscode, err.Error())
}

func returnJson(w http.ResponseWriter, data interface{}) {
	enc := json.NewEncoder(w)
	if err := enc.Encode(&data); err != nil {
		handleError(w, 500, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
}

func AccountEntry(w http.ResponseWriter, r *http.Request) {
	doLog(r)
	if r.Method == "GET" {
		list := model.ListAccounts()
		returnJson(w, struct {
			Accounts []model.ResultAccount
		}{
			list,
		})
	} else if r.Method == "PUT" {
		name, ok := getURLPath(r, 2)
		if !ok || name == "" {
			handleRawError(w, 400, "No name provided.")
			return
		}
		acc := model.CreateAccount(name)
		returnJson(w, acc)
	} else {
		handleRawError(w, http.StatusMethodNotAllowed, "Method not allowed.  Only supports GET and PUT.")
	}
}

func TransactionEntry(w http.ResponseWriter, r *http.Request) {
	doLog(r)
	if r.Method == "GET" {
		accountid, ok := getURLPathInt(r, 2)
		if !ok {
			handleRawError(w, 400, "Bad or missing account id.")
			return
		}
		list, err := model.ListTransactions(accountid)
		if err != nil {
			handleError(w, 400, err)
			return
		}
		returnJson(w, struct {
			Transactions []model.ResultTransaction
		}{
			list,
		})
	} else if r.Method == "PUT" {
		accountid, ok := getURLPathInt(r, 2)
		if !ok {
			handleRawError(w, 400, "Bad or missing account id.")
			return
		}
		dec := json.NewDecoder(r.Body)
		var input struct {
			Description string
			Amount      float64
		}
		if err := dec.Decode(&input); err != nil {
			handleError(w, 400, err)
			return
		}
		tid, err := model.CreateTransaction(accountid, input.Amount, input.Description)
		if err != nil {
			handleError(w, 400, err)
			return
		}
		returnJson(w, struct {
			TransactionId string
		}{
			string(tid),
		})
	} else if r.Method == "DELETE" {
		accountid, ok := getURLPathInt(r, 2)
		if !ok {
			handleRawError(w, 400, "Bad or missing account id.")
			return
		}
		tid, ok := getURLPath(r, 3)
		if !ok {
			handleRawError(w, 400, "Bad or missing transaction ID.")
			return
		}
		err := model.DeleteTransaction(accountid, model.TransactionId(tid))
		if err != nil {
			handleError(w, 400, err)
			return
		}
		returnJson(w, struct {
			TransactionId string
		}{
			string(tid),
		})
	} else {
		handleRawError(w, http.StatusMethodNotAllowed, "Method not allowed.  Only supports GET, PUT and DELETE.")
	}
}

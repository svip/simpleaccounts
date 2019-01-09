package model

import (
	"db"
	"fmt"
	"time"
)

type TransactionId string

type ResultAccount struct {
	Id int
	Name string
	Sum float64
}

type ResultTransaction struct {
	Id TransactionId
	Amount float64
	Description string
}

func toTransactionId(t time.Time) TransactionId {
	return TransactionId(t.Format(time.RFC3339Nano))
}

func fromTransactionId(id TransactionId) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, string(id))
}

func CreateAccount(name string) ResultAccount {
	acc := db.CreateAccount(name)
	return ResultAccount{acc.Id, acc.Name, acc.Sum().Float()}
}

func ListAccounts() []ResultAccount {
	dblist := db.GetAccounts()
	var list []ResultAccount
	for _, acc := range dblist {
		list = append(list, ResultAccount{acc.Id, acc.Name, acc.Sum().Float()})
	}
	return list
}

func ListTransactions(accountid int) ([]ResultTransaction, error) {
	acc, err := db.GetAccount(accountid)
	var list []ResultTransaction
	if err != nil {
		return list, err
	}
	for _, trans := range acc.Transactions {
		list = append(list, ResultTransaction{toTransactionId(trans.Time), trans.Amount.Float(), trans.Description})
	}
	return list, nil
}

func CreateTransaction(accountid int, amount float64, description string) (TransactionId, error) {
	m := db.NewMoney(amount)
	if m.IsZero() {
		return "", fmt.Errorf("Transaction cannot have an amount of 0.00")
	}
	tid, err := db.CreateTransaction(accountid, m, description)
	if err != nil {
		return "", err
	}
	return toTransactionId(tid), nil
}

func DeleteTransaction(accountid int, transid TransactionId) error {
	tid, err := fromTransactionId(transid)
	if err != nil {
		return err
	}
	
	err = db.DeleteTransaction(accountid, tid)
	if err != nil {
		return err
	}
	
	return nil
}


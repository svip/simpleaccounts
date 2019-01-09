package model

import (
	"db"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	db.Init(true)
	testacc := ResultAccount{1, "TestName", 0.00}
	acc := CreateAccount(testacc.Name)
	if acc != testacc {
		t.Error("Expected account ID 1, name 'TestName' and a sum of 0.00, got ", acc)
	}
}

func TestListAccounts(t *testing.T) {
	db.Init(true)
	testacc := ResultAccount{1, "TestName", 0.00}
	CreateAccount(testacc.Name)
	list := ListAccounts()
	if len(list) != 1 {
		t.Error("Expected 1 account, got ", len(list))
	}
	if list[0] != testacc {
		t.Error("Expected account ID 1, name 'TestName' and a sum of 0.00, got ", list[0])
	}
}

func TestCreateTransaction(t *testing.T) {
	db.Init(true)
	CreateAccount("TestName")
	_, err := CreateTransaction(1, 10.50, "Text1")
	if err != nil {
		t.Error("Expected no error, got ", err)
	}
	list, err := ListTransactions(1)
	if err != nil {
		t.Error("Expected no error, got ", err)
	}
	if len(list) != 1 {
		t.Error("Expected 1 transaction, got ", len(list))
	}
	if list[0].Amount != 10.50 {
		t.Error("Expected an amount of 10.50, got ", list[0].Amount)
	}
	if list[0].Description != "Text1" {
		t.Error("Expected a description of Text1, got ", list[0].Description)
	}
	acc := ListAccounts()[0]
	if acc.Sum != 10.50 {
		t.Error("Expected a sum of 10.50, got ", acc.Sum)
	}
	_, err = CreateTransaction(1, 0.00, "Text2")
	if err == nil {
		t.Error("Expected error, got nothing")
	}
}

func TestListTransactions(t *testing.T) {
	db.Init(true)
	CreateAccount("TestName")
	_, err := CreateTransaction(1, 10.50, "Text1")
	if err != nil {
		t.Error("Expected no error, got ", err)
	}
	_, err = CreateTransaction(1, 19.25, "Text2")
	if err != nil {
		t.Error("Expected no error, got ", err)
	}
	list, err := ListTransactions(1)
	if err != nil {
		t.Error("Expected no error, got ", err)
	}
	if len(list) != 2 {
		t.Error("Expected 2 transactions, got ", len(list))
	}
	if list[1].Amount != 19.25 {
		t.Error("Expected an amount of 19.25, got ", list[1].Amount)
	}
	if list[1].Description != "Text2" {
		t.Error("Expected a description of Text2, got ", list[1].Description)
	}
	acc := ListAccounts()[0]
	if acc.Sum != 29.75 {
		t.Error("Expected a sum of 29.75, got ", acc.Sum)
	}
}

func TestDeleteTransaction(t *testing.T) {
	db.Init(true)
	CreateAccount("TestName")
	tid1, err := CreateTransaction(1, 10.50, "Text1")
	if err != nil {
		t.Error("Expected no error, got ", err)
	}
	err = DeleteTransaction(1, tid1)
	if err != nil {
		t.Error("Expected no error, got ", err)
	}
	list1, _ := ListTransactions(1)
	if len(list1) != 0 {
		t.Error("Expected no transactions, got ", len(list1))
	}
	_, err = CreateTransaction(1, 19.25, "Text2")
	if err != nil {
		t.Error("Expected no error, got ", err)
	}
	list2, _ := ListTransactions(1)
	tid2 := list2[0].Id
	err = DeleteTransaction(1, tid2)
	if err != nil {
		t.Error("Expected no error, got ", err)
	}
	list3, _ := ListTransactions(1)
	if len(list3) != 0 {
		t.Error("Expected no transactions, got ", len(list3))
	}
}

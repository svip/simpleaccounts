package db

import (
	"testing"
)

// Right now, the tests are very basic.  But should we run into further errors,
// we can now easier those tests to these.

func TestCreateAccount(t *testing.T) {
	Init(true)
	testacc := Account{1, "TestName", []Transaction{}}
	acc1 := CreateAccount(testacc.Name)
	if acc1.Id != testacc.Id || acc1.Name != testacc.Name {
		t.Error("Expected an account with Id 1 and name 'TestName', got ", acc1)
	}
	acc2, err := GetAccount(1)
	if err != nil {
		t.Error("Expected no error, got ", err)
	}
	if acc2.Id != testacc.Id || acc2.Name != testacc.Name {
		t.Error("Expected an account with Id 1 and name 'TestName', got ", acc2)
	}
}

func TestGetAccounts(t *testing.T) {
	Init(true)
	testacc := Account{1, "TestName", []Transaction{}}
	CreateAccount(testacc.Name)
	list := GetAccounts()
	if len(list) != 1 {
		t.Error("Expected a list of 1 accounts, got ", len(list))
	}
	acc1 := list[0]
	if acc1.Id != testacc.Id || acc1.Name != testacc.Name {
		t.Error("Expected an account with Id 1 and name 'TestName', got ", acc1)
	}
}

func TestGetAccount(t *testing.T) {
	Init(true)
	testacc := Account{1, "TestName", []Transaction{}}
	acc1 := CreateAccount(testacc.Name)
	if acc1.Id != testacc.Id || acc1.Name != testacc.Name {
		t.Error("Expected an account with Id 1 and name 'TestName', got ", acc1)
	}
	acc2, err := GetAccount(1)
	if err != nil {
		t.Error("Expected no error, got ", err)
	}
	if acc2.Id != testacc.Id || acc2.Name != testacc.Name {
		t.Error("Expected an account with Id 1 and name 'TestName', got ", acc2)
	}
}

func TestCreateTransaction(t *testing.T) {
	Init(true)
	testacc := Account{1, "TestName", []Transaction{}}
	CreateAccount(testacc.Name)
	_, err := CreateTransaction(1, NewMoney(10.00), "Text1")
	if err != nil {
		t.Error("Expected no error, got ", err)
	}
	acc1, _ := GetAccount(1)
	if len(acc1.Transactions) != 1 {
		t.Error("Expected just one transaction, got ", len(acc1.Transactions))
	}
	if acc1.Transactions[0].Amount != NewMoney(10.00) {
		t.Error("Expected an amount of 10.00, got ", acc1.Transactions[0].Amount)
	}
	if acc1.Transactions[0].Description != "Text1" {
		t.Error("Expected description of Text1, got ", acc1.Transactions[0].Description)
	}
	if acc1.Sum() != NewMoney(10.00) {
		t.Error("Expected a sum of 10.00, got ", acc1.Sum())
	}
	
	CreateTransaction(1, NewMoney(20.00), "Text2")
	acc2, _ := GetAccount(1)
	if acc2.Transactions[1].Amount != NewMoney(20.00) {
		t.Error("Expected an amount of 20.00, got ", acc2.Transactions[1].Amount)
	}
	if acc2.Transactions[1].Description != "Text2" {
		t.Error("Expected description of Text2, got ", acc2.Transactions[0].Description)
	}
	if acc2.Sum() != NewMoney(30.00) {
		t.Error("Expected a sum of 30.00, got ", acc2.Sum())
	}
	
	CreateTransaction(1, NewMoney(10.009), "Text3")
	acc3, _ := GetAccount(1)
	if acc3.Transactions[2].Amount != NewMoney(10.00) {
		t.Error("Expected an amount of 10.00, got ", acc3.Transactions[2].Amount)
	}
	if acc3.Sum() != NewMoney(40.00) {
		t.Error("Expected a sum of 40.00, got ", acc3.Sum())
	}
}

func TestDeleteTransaction(t *testing.T) {
	Init(true)
	testacc := Account{1, "TestName", []Transaction{}}
	CreateAccount(testacc.Name)
	CreateTransaction(1, NewMoney(10.00), "Text1")
	
	acc1, _ := GetAccount(1)
	id := acc1.Transactions[0].Time
	err := DeleteTransaction(1, id)
	if err != nil {
		t.Error("Expected no error, got ", err)
	}
	
	acc2, _ := GetAccount(1)
	if len(acc2.Transactions) != 0 {
		t.Error("Expected no transactions, got ", len(acc2.Transactions))
	}
}


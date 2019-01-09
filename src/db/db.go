package db

import (
	"encoding/json"
	"time"
	"fmt"
	"sync"
	"sort"
	"os"
	"log"
	"io"
)

// TYPES

// Money is generally accepted to have two decimal numbers, but storing floats
// can have some unforseen consequences (like 0.1+0.2!=0.3).
type Money int

func NewMoney(a float64) Money {
	// We'll just any further precision.  So yes, 1.505 becomes 1.50.
	return Money(int(a*100))
}

func (m Money) ToFloat() float64 {
	return float64(m)/100
}

func (m Money) String() string {
	return fmt.Sprintf("%0.2f", m.ToFloat())
} 

type Transaction struct {
	Time time.Time
	Amount Money
	Description string
}

type Account struct {
	Id int
	Name string
	Transactions []Transaction
}

func (a Account) Sum() Money {
	var sum Money
	for _, trans := range a.Transactions {
		sum += trans.Amount
	}
	return sum
}

func (a Account) String() string {
	return fmt.Sprintf("%d: %s (%s)", a.Id, a.Name, a.Sum().String())
}

// GENERAL VARIABLES AND HELPER FUNCTIONS

// Our database is merely a list of accounts.
// Stored in a map for easy lookup.
var database map[int]Account
var mutex = &sync.Mutex{} // Our general database mutex.
var savetofile bool
const FILENAME = "database.json"

func save() error {
	if !savetofile {
		return nil
	}
	file, err := os.Create(FILENAME)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	if err := enc.Encode(&database); err != nil {
		return err
	}
	return nil
}

// EXPORTED FUNCTIONS

func Init(nofile bool) {
	savetofile = !nofile
	if savetofile {
		file, err := os.Create(FILENAME)
		if err != nil {
			// If we cannot even open the file, then we might as well not run at all
			log.Fatal(err)
		}
		defer file.Close()
		dec := json.NewDecoder(file)
		mutex.Lock()
		if err = dec.Decode(&database); err != nil {
			if err == io.EOF {
				database = make(map[int]Account)
			} else {
				log.Fatal(err)
			}
		}
		mutex.Unlock()
	} else {
		database = make(map[int]Account)
	}
}

func EmptyDatabase() {
	mutex.Lock()
	database = make(map[int]Account)
	mutex.Unlock()
}

func CreateAccount(name string) Account {
	// We are keeping it simple, with just incrementing IDs.
	mutex.Lock()
	newid := len(database)+1
	
	account := Account{newid, name, []Transaction{}}
	
	database[newid] = account
	
	save()
	
	mutex.Unlock()
	
	return account
}

type byId []Account

func (a byId) Len() int           { return len(a) }
func (a byId) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byId) Less(i, j int) bool { return a[i].Id < a[j].Id }

func GetAccounts() []Account {
	var list []Account
	mutex.Lock()
	for _, acc := range database {
		list = append(list, acc)
	}
	mutex.Unlock()
	sort.Sort(byId(list))
	return list
}

func GetAccount(id int) (Account, error) {
	mutex.Lock()
	defer mutex.Unlock()
	if acc, ok := database[id]; ok {
		return acc, nil
	} else {
		return Account{}, fmt.Errorf("Account id %d not found.", id)
	}
}

func CreateTransaction(accountid int, amount Money, description string) error {
	acc, err := GetAccount(accountid)
	if err != nil {
		return err
	}
	
	acc.Transactions = append(acc.Transactions, Transaction{time.Now(), amount, description})
	
	mutex.Lock()
	database[accountid] = acc
	save()
	mutex.Unlock()
	
	return nil
}

func DeleteTransaction(accountid int, transid time.Time) error {
	acc, err := GetAccount(accountid)
	if err != nil {
		return err
	}
	
	id := -1
	for i, trans := range acc.Transactions {
		if trans.Time == transid {
			id = i
		}
	}
	if id == -1 {
		return fmt.Errorf("No transaction with id (time) %x", transid)
	}
	
	acc.Transactions = append(acc.Transactions[:id], acc.Transactions[id+1:]...)
	
	mutex.Lock()
	database[accountid] = acc
	save()
	mutex.Unlock()
	
	return nil
}


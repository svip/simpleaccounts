Simple account and transaction database
=======================================

Set up
------

Simply start `main.go`:

```
go run src/main.go
```

It will run on port `:8080`.  It will save the database to a file called
`database.json` in the same directory it is being run from.  If there is no
file, or it is empty, it will initialise a new database without anything.

Supported calls
---------------

### GET /account

Returns the list of accounts.

#### Example

Input:
`GET /account/`

Output:
```json
{
  "Accounts": [
    {
      "Id": 1,
      "Name": "Person Name",
      "Sum": 25.25
    }
  ]
} 
```

### PUT /account/{Name}

Creates a new account with name as the first positional parameter.  No body
expected.  Then returns the new account.

#### Example

Input:
`PUT /account/Person%20Name`

Output:
```json
{
  "Id": 1,
  "Name": "Person Name",
  "Sum": 0
}
```

### GET /transaction/{AccountId}

Get the list of transactions for a given account.  Will return error if the
account does not exist.

#### Example

Input:
`GET /account/1`

Output:
```json
{
  "Transactions": [
    {
      "Id": "2019-01-09T19:52:10.038329163+01:00",
      "Amount": 10.25,
      "Description": "Money for something."
    }
  ]
}
```

### PUT /transaction/{AccountId}

Adds another transaction to the account behind the account ID.  Expects a JSON
body as follows:

```json
{
  "Amount": 0.0,
  "Description": "Text"
}
```

Will return error if the account does not exist.  Will return an error, if
either of the parameters are missing.  Will return an error, if the amount is 0.

On success, it returns the transaction ID.

#### Example

Input:
`PUT /transaction/1`
```json
{
  "Amount": 25.5,
  "Description": "Here is money."
}
```

Output:
```json
{
  "TransactionId": "2019-01-09T19:52:10.038329163+01:00"
}
```

### DELETE /transaction/{AccountId}/{TransactionId}

Deletes an existing transaction.  Returns an error if the account or transaction
does not exist.  Returns the transaction id it just deleted if successful.

#### Example

Input
`DELETE /transaction/1/2019-01-09T19:52:10.038329163+01:00`

Output:
```json
{
  "TransactionId": "2019-01-09T19:52:10.038329163+01:00"
}
```


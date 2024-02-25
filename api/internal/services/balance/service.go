package balance

import (
	"database/sql"
	"log"
	"rinha-backend-2024/api/internal/database"
	"rinha-backend-2024/api/internal/model"
	"rinha-backend-2024/api/internal/model/exception"
	"time"
)

func InsertTransaction(clientId int, transaction model.Transaction) (model.Client, interface{}) {
	conn := database.GetConnection()

	client := model.Client{}

	tx, err := conn.Begin()
	if err != nil {
		log.Println("Error beginning transaction:", err)
		return client, exception.TransactionError{Message: "Error begin transaction"}
	}

	defer tx.Rollback()

	// Pessimist locking with `FOR UPDATE`
	err = tx.QueryRow(`SELECT id, balanceLimit, balance FROM client WHERE id=$1 FOR UPDATE`, clientId).Scan(&client.Id, &client.BalanceLimit, &client.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return client, exception.UserNotFound{Message: "User not found"}
		}
		return client, exception.TransactionError{Message: "Error retrieving client information"}
	}

	// Validate transaction value and type
	newBalance := transaction.Value
	if transaction.Type == "d" {
		limit := client.Balance - newBalance

		if limit < (client.BalanceLimit * -1) {
			return client, exception.UnprocessableEntity{
				Message: "Value not accepted for 'd'",
			}
		}

		newBalance = -transaction.Value
	}

	_, err = tx.Exec(`INSERT INTO transaction (client_id, value, type, description, date) VALUES ($1, $2, $3, $4, $5)`, clientId, transaction.Value, transaction.Type, transaction.Description, time.Now())
	if err != nil {
		log.Println("Error inserting transaction:", err)
		return client, exception.TransactionError{Message: "Error insert transaction"}
	}

	// fix concurrency
	//the sum at the time of the update guarantees that the "balance" value will be the most updated
	err = tx.QueryRow(`UPDATE client SET balance = balance + $1 WHERE id = $2 RETURNING balance`, newBalance, clientId).Scan(&client.Balance)
	if err != nil {
		log.Println("Error updating client balance:", err)
		return client, exception.TransactionError{Message: "Error updating client balance"}
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Error committing transaction:", err)
		return client, exception.TransactionError{Message: "Error committing transaction"}
	}

	return client, nil
}

func GetExtractByUserId(clientId int) (model.Extract, interface{}) {
	conn := database.GetConnection()

	sqlResult := conn.QueryRow(`SELECT id, balanceLimit, balance FROM client WHERE id=$1`, clientId)

	client := model.Client{}

	err := sqlResult.Scan(&client.Id, &client.BalanceLimit, &client.Balance)

	if err != nil && err == sql.ErrNoRows {
		return model.Extract{}, exception.UserNotFound{
			Message: "User not found",
		}

	}

	transactionRows, err := conn.Query(`
	SELECT t.value, t.type, t.description, t.date 
	FROM transaction t WHERE t.client_id = $1 
	ORDER BY t.date DESC LIMIT 10;`, clientId)

	if err != nil && err == sql.ErrNoRows {
		return model.Extract{
			Saldo: model.Balance{
				Total: client.Balance,
				Limit: client.BalanceLimit,
				Date:  time.Now(),
			},
		}, nil

	}

	transactions := []model.Transaction{}

	for transactionRows.Next() {

		transaction := model.Transaction{}
		errT := transactionRows.Scan(&transaction.Value, &transaction.Type, &transaction.Description, &transaction.Date)

		if errT != nil {
			log.Println("Error convert transaction on extract")
			continue
		}

		transactions = append(transactions, transaction)
	}

	return model.Extract{
		Saldo: model.Balance{
			Total: client.Balance,
			Limit: client.BalanceLimit,
			Date:  time.Now(),
		},
		Transactions: transactions,
	}, nil
}

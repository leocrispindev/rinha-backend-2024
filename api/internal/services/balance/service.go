package balance

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"rinha-backend-2024/api/internal/database"
	"rinha-backend-2024/api/internal/model"
	"rinha-backend-2024/api/internal/model/exception"
	"time"

	"github.com/jackc/pgx/v4"
)

func InsertTransaction(clientId int, transaction model.Transaction) (model.Balance, interface{}) {
	conn, err := database.GetConnection()

	if err != nil {
		return model.Balance{}, exception.TransactionError{
			Message: "Error on acquire connection",
		}
	}

	balance := model.Balance{}

	tx, err := conn.Begin(context.Background())
	if err != nil {
		log.Println("Error beginning transaction:", err)
		return balance, exception.TransactionError{Message: "Error begin transaction"}
	}

	defer tx.Rollback(context.Background())
	defer database.ReleaseConnection(conn)

	// Pessimist locking with `FOR UPDATE`
	err = tx.QueryRow(context.Background(), `SELECT balanceLimit, balance FROM client WHERE id=$1 FOR UPDATE`, clientId).Scan(&balance.BalanceLimit, &balance.Balance)
	if err != nil {
		if err == pgx.ErrNoRows {
			return balance, exception.UserNotFound{Message: "User not found"}
		}
		fmt.Printf("userId=%d", clientId)
		return balance, exception.TransactionError{Message: "Error retrieving client information"}
	}

	// Validate transaction value and type
	newBalance := transaction.Value
	if transaction.Type == "d" {
		limit := balance.Balance - newBalance

		if limit < (balance.BalanceLimit * -1) {
			return balance, exception.UnprocessableEntity{
				Message: "Value not accepted for 'd'",
			}
		}

		newBalance = -transaction.Value
	}

	_, err = tx.Exec(context.Background(), `INSERT INTO transaction (client_id, value, type, description, date) VALUES ($1, $2, $3, $4, $5)`, clientId, transaction.Value, transaction.Type, transaction.Description, time.Now())
	if err != nil {
		log.Println("Error inserting transaction:", err)
		return balance, exception.TransactionError{Message: "Error insert transaction"}
	}

	// fix concurrency
	//the sum at the time of the update guarantees that the "balance" value will be the most updated
	err = tx.QueryRow(context.Background(), `UPDATE client SET balance = balance + $1 WHERE id = $2 RETURNING balance`, newBalance, clientId).Scan(&balance.Balance)
	if err != nil {
		log.Println("Error updating client balance:", err)
		return balance, exception.TransactionError{Message: "Error updating client balance"}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		log.Println("Error committing transaction:", err)
		return balance, exception.TransactionError{Message: "Error committing transaction"}
	}

	return balance, nil
}

func GetExtractByUserId(clientId int) (model.Extract, interface{}) {
	conn, err := database.GetConnection()

	if err != nil {
		return model.Extract{}, exception.TransactionError{
			Message: "Error on acquire connection",
		}
	}

	defer database.ReleaseConnection(conn)

	sqlResult := conn.QueryRow(context.Background(), `SELECT balanceLimit, balance FROM client WHERE id=$1`, clientId)

	balance := model.Balance{}

	err = sqlResult.Scan(&balance.BalanceLimit, &balance.Balance)

	if err != nil {
		if err == pgx.ErrNoRows {
			return model.Extract{}, exception.UserNotFound{Message: "User not found"}
		}
		fmt.Printf("userId=%d", clientId)
		return model.Extract{}, exception.TransactionError{Message: "Error retrieving client information"}
	}

	transactionRows, err := conn.Query(context.Background(), `
	SELECT t.value, t.type, t.description, t.date 
	FROM transaction t WHERE t.client_id = $1 
	ORDER BY t.date DESC LIMIT 10;`, clientId)

	if err != nil && err == sql.ErrNoRows {
		return model.Extract{
			Saldo: model.ExtractBalance{
				Total:        balance.Balance,
				BalanceLimit: balance.BalanceLimit,
				Date:         time.Now(),
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
		Saldo: model.ExtractBalance{
			Total:        balance.Balance,
			BalanceLimit: balance.BalanceLimit,
			Date:         time.Now(),
		},
		Transactions: transactions,
	}, nil
}

package balance

import (
	"database/sql"
	"fmt"
	"log"
	"rinha-backend-2024/api/internal/database"
	"rinha-backend-2024/api/internal/model"
	"rinha-backend-2024/api/internal/model/exception"
	"time"
)

func InsertTransaction(clientId int, transaction model.Transaction) (model.Client, interface{}) {
	conn := database.GetConnection()

	client := model.Client{}

	sqlTransaction, err := conn.Begin() // TODO  tratar o erro na abertura da transaction

	if err != nil {
		log.Println("error on begin transaction clientID: " + fmt.Sprintf("%d", clientId))
		return client, exception.TransactionError{
			Message: "Error on begin transaction",
		}
	}

	defer sqlTransaction.Rollback() // Se a transação não for commitada, ocorre o rollback

	// usando o id como chave para bloqueio
	// será liberado no commit ou no rollback

	sqlTransaction.Exec("SELECT pg_advisory_xact_lock($1)", clientId)

	clientRow := sqlTransaction.QueryRow(`SELECT id, balanceLimit, balance FROM client WHERE id=$1`, clientId)

	err = clientRow.Scan(&client.Id, &client.BalanceLimit, &client.Balance)

	if err != nil || err == sql.ErrNoRows {
		log.Println(err)
		return client, exception.TransactionError{
			Message: "Internal Server Error",
		}

	}

	newBalance := transaction.Value

	if transaction.Type == "d" {
		newBalance = -transaction.Value
	}

	if transaction.Type == "d" && client.Balance-transaction.Value < client.BalanceLimit*-1 {
		return client, exception.UnprocessableEntity{
			Message: "You have no limit for this transaction",
		}
	}

	_, err = sqlTransaction.Exec(`INSERT INTO transaction (client_id, value, type, description, date)  VALUES 
	($1, $2, $3, $4, $5)`, client.Id, transaction.Value, transaction.Type, transaction.Description, time.Now())

	if err != nil {
		log.Println("error on insert transaction clientID: " + fmt.Sprintf("%d", clientId))

		return client, exception.TransactionError{
			Message: "Error insert transaction",
		}
	}

	updateRow := sqlTransaction.QueryRow(`UPDATE client SET balance=balance + $1 WHERE id=$2 RETURNING balance`, newBalance, client.Id)

	err = updateRow.Scan(&client.Balance)

	if err != nil {
		return client, exception.TransactionError{
			Message: "Error on update client balance",
		}
	}

	err = sqlTransaction.Commit()

	if err != nil {
		log.Println("error on commit transaction clientID: " + fmt.Sprintf("%d", clientId))
		return client, exception.TransactionError{
			Message: "Error on commit message",
		}
	}

	return client, nil
}

func GetExtractByUserId(clientId int) (model.Extract, interface{}) {
	conn := database.GetConnection()

	sqlTransaction, err := conn.Begin() // TODO  tratar o erro na abertura da transaction

	if err != nil {
		log.Println("Error on begin transaction (Statement) clientID: " + fmt.Sprintf("%d", clientId))
		return model.Extract{}, exception.TransactionError{
			Message: "Error on begin transaction",
		}
	}

	defer sqlTransaction.Rollback() // Se a transação não for commitada, ocorre o rollback

	sqlResult := sqlTransaction.QueryRow(`SELECT id, balanceLimit, balance FROM client WHERE id=$1`, clientId)

	client := model.Client{}

	err = sqlResult.Scan(&client.Id, &client.BalanceLimit, &client.Balance)

	if err != nil && err == sql.ErrNoRows {
		return model.Extract{}, exception.UserNotFound{
			Message: "User not found",
		}

	}

	transactionRows, err := sqlTransaction.Query(`
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

	err = sqlTransaction.Commit()

	if err != nil {
		log.Println("error on commit transaction (statement) clientID: " + fmt.Sprintf("%d", clientId))
		return model.Extract{}, exception.TransactionError{
			Message: "Error on commit message",
		}
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

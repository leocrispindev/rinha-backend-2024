package balance

import (
	"database/sql"
	"fmt"
	"rinha-backend-2024/api/internal/database"
	"rinha-backend-2024/api/internal/model"
	"time"
)

func InsertTransaction(clientId int, transaction model.Transaction) (model.Client, error) {
	conn := database.GetConnection()

	client := model.Client{}

	sqlTransaction, err := conn.Begin() // TODO  tratar o erro na abertura da transaction

	defer sqlTransaction.Rollback() // Se a transação não for commitada, ocorre o rollback

	clientRow := sqlTransaction.QueryRow(`SELECT id, balanceLimit, balance FROM client WHERE id=$1 FOR UPDATE`, clientId)

	err = clientRow.Scan(&client.Id, &client.BalanceLimit, &client.Balance)

	if err != nil || err == sql.ErrNoRows {
		return client, fmt.Errorf("Client not found")

	}

	newBalance := transaction.Value + client.Balance
	if transaction.Type == "d" && newBalance < client.BalanceLimit {

		return client, fmt.Errorf("Operation not ok")
	}

	_, err = sqlTransaction.Exec(`INSERT INTO transaction (client_id, value, type, description, date)  VALUES 
	($1, $2, $3, $4, $5)`, client.Id, transaction.Value, transaction.Type, transaction.Description, time.Now())

	if err != nil {
		return client, fmt.Errorf("Error on insert transaction")
	}

	_, err = sqlTransaction.Exec(`UPDATE client SET balance=$1 WHERE id=$2`, newBalance, client.Id)

	if err != nil {
		return client, fmt.Errorf("Error on update client balance")
	}

	client.Balance = newBalance

	err = sqlTransaction.Commit()

	if err != nil {
		return client, fmt.Errorf("Error on commi transaction")
	}

	return client, nil
}

func GetExtractByUserId(clientId int) (model.Extract, error) {
	conn := database.GetConnection()

	sqlTransaction, err := conn.Begin() // TODO  tratar o erro na abertura da transaction

	defer sqlTransaction.Rollback() // Se a transação não for commitada, ocorre o rollback

	clientRow := sqlTransaction.QueryRow(`SELECT id, balanceLimit, balance FROM client WHERE id=$1 FOR UPDATE`, clientId)

	client := model.Client{}

	err = clientRow.Scan(&client.Id, &client.BalanceLimit, &client.Balance)

	if err != nil || err == sql.ErrNoRows {
		return model.Extract{}, fmt.Errorf("Client not found")

	}

	transactionRows, err := sqlTransaction.Query(`
	SELECT t.value, t.type, t.description, t.date 
	FROM transaction t WHERE t.client_id = $1 
	ORDER BY t.date DESC LIMIT 10;`, clientId)

	if err != nil || err == sql.ErrNoRows {
		return model.Extract{}, fmt.Errorf("Error retrieve client transactions: " + err.Error())

	}

	transactions := []model.Transaction{}

	for transactionRows.Next() {

		transaction := model.Transaction{}
		errT := transactionRows.Scan(&transaction.Value, &transaction.Type, &transaction.Description, &transaction.Date)

		if errT != nil {
			continue
		}

		transactions = append(transactions, transaction)
	}

	err = sqlTransaction.Commit()
	if err != nil {
		return model.Extract{}, fmt.Errorf("Error execute extract: " + err.Error())
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

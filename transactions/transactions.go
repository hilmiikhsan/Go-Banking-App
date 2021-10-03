package transactions

import (
	"banking-app/database"
	"banking-app/helpers"
	"banking-app/interfaces"
)

func CreateTransaction(From uint, To uint, Amount int) {
	// db := helpers.ConnectDB()
	transaction := &interfaces.Transaction{From: From, To: To}
	database.DB.Create(&transaction)

	// defer db.Close()
}

func GetTransactionsByAccount(id uint) []interfaces.ResponseTransaction {
	transactions := []interfaces.ResponseTransaction{}
	database.DB.Table("transactions").Select("id, transactions.from, transactions.to, amount").Where(interfaces.Transaction{From: id}).Or(interfaces.Transaction{To: id}).Scan(&transactions)
	return transactions
}

func GetMyTransactions(id string, jwt string) map[string]interface{} {
	isValid := helpers.ValidateToken(id, jwt)

	if isValid {
		accounts := []interfaces.ResponseAccount{}

		database.DB.Table("accounts").Select("id, name, balance").Where("user_id = ?", id).Scan(&accounts)

		transactions := []interfaces.ResponseTransaction{}

		for i := 0; i < len(accounts); i++ {
			accountTransaction := GetTransactionsByAccount(accounts[i].ID)
			transactions = append(transactions, accountTransaction...)
		}

		var response = map[string]interface{}{"message": "Semuanya Sukses!"}
		response["data"] = transactions
		return response

	} else {
		return map[string]interface{}{"message": "Token Tidak Valid"}
	}
}
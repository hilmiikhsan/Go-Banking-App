package useraccounts

import (
	"banking-app/database"
	"banking-app/helpers"
	"banking-app/interfaces"
	"banking-app/transactions"
	"fmt"
)

func updateAccount(id uint, amount int) interfaces.ResponseAccount {
	// db := helpers.ConnectDB()
	account := interfaces.Account{}
	responseAccount := interfaces.ResponseAccount{}
	// db.Where("id = ? ", id).First(&account)
	database.DB.Where("id = ? ", id).First(&account)
	account.Balance = uint(amount)
	// db.Save(&account)
	database.DB.Save(&account)

	responseAccount.ID = account.ID
	responseAccount.Name = account.Name
	responseAccount.Balance = int(account.Balance)

	// defer db.Close()
	return responseAccount
}

func getAccount(id uint) *interfaces.Account {
	// db := helpers.ConnectDB()
	account := &interfaces.Account{}
	if database.DB.Where("id = ?", id).First(&account).RecordNotFound() {
		return nil
	}

	// defer db.Close()
	return account
}

func Transaction(userId uint, from uint, to uint, amount int, jwt string) map[string]interface{} {
	userIdString := fmt.Sprint(userId)
	isValid := helpers.ValidateToken(userIdString, jwt)

	if isValid {
		fromAccount := getAccount(from)
		toAccount := getAccount(to)

		if fromAccount == nil || toAccount == nil {
			return map[string]interface{}{"message": "Akun Tidak Ditemukan"}
		} else if fromAccount.UserID != userId {
			return map[string]interface{}{"message": "Anda Bukan Pemilik Akun Ini"}
		} else if int(fromAccount.Balance) < amount {
			return map[string]interface{}{"message": "Saldo Akun Anda Terlalu Kecil"}
		}

		updatedAccount := updateAccount(from, int(fromAccount.Balance) - amount)
		updateAccount(to, int(toAccount.Balance) + amount)

		transactions.CreateTransaction(from, to, amount)

		var response = map[string]interface{}{"message": "Semuanya Sukses!"}
		response["data"] = updatedAccount

		return response

	} else {
		return map[string]interface{}{"message": "Token Tidak Valid"}
	}
}
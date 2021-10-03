package migrations

import (
	"banking-app/database"
	"banking-app/helpers"
	"banking-app/interfaces"
)

func createAccounts() {
	// db := helpers.ConnectDB()

	users := &[2]interfaces.User {
		{Username: "Ikhsan", Email: "robbinhood22012@gmail.com"},
		{Username: "Risti", Email: "ristiindria@gmail.com"},
	}

	for i := 0; i < len(users); i++ {
		generatePassword := helpers.HashAndSalt([]byte(users[i].Username))
		user := &interfaces.User{Username: users[i].Username, Email: users[i].Email, Password: generatePassword}
		// db.Create(&user)
		database.DB.Create(&user)

		account := &interfaces.Account{Type: "Daily Account", Name: string(users[i].Username + "'s" + " account"), Balance: uint(10000 * int(i+1)), UserID: user.ID}
		// db.Create(&account)
		database.DB.Create(&account)
	}
	// defer db.Close()
}

func Migrate() {
	User := &interfaces.User{}
	Account := &interfaces.Account{}
	Transactions := &interfaces.Transaction{}
	// db := helpers.ConnectDB()
	// db.AutoMigrate(&User, &Account)
	database.DB.AutoMigrate(&User, &Account, &Transactions)
	// defer db.Close()

	createAccounts()
}

// func MigrateTransactions() {
// 	Transactions := &interfaces.Transaction{}

// 	db := helpers.ConnectDB()
// 	db.AutoMigrate(&Transactions)
// 	defer db.Close()
// }
package main

import (
	"banking-app/api"
	"banking-app/database"
)

func main() {
	// Do migration
	// migrations.MigrateTransactions()
	database.InitDatabase()
	api.StartApi()
}
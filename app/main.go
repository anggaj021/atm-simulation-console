package main

import (
	atm_handler "atm-simulation-console/internal/handler/atm"
	account_csv "atm-simulation-console/internal/repository/account/csv"
	transaction_csv "atm-simulation-console/internal/repository/transaction/csv"
	atm_service "atm-simulation-console/internal/service/atm"
)

func main() {
	// accountRepo := in_memory.NewInMemoryAccount()
	accpath := "./data/userdata.csv"
	trxpath := "./data/transaction.csv"
	accountRepo := account_csv.NewCSVAccountRepository(accpath)
	trxRepo := transaction_csv.NewCSVTransactionRepository(trxpath)
	atmSvc := atm_service.NewATMService(accountRepo, trxRepo)

	atmController := atm_handler.NewATMController(atmSvc)

	atmController.Start()
}

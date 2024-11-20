package main

import (
	atm_handler "atm-simulation-console/internal/handler/atm"
	"atm-simulation-console/internal/repository/account/csv"
	atm_service "atm-simulation-console/internal/service/atm"
)

func main() {
	// accountRepo := in_memory.NewInMemoryAccount()
	filepath := "./data/userdata.csv"
	accountRepo := csv.NewCSVAccountRepository(filepath)
	atmSvc := atm_service.NewATMService(accountRepo)

	atmController := atm_handler.NewATMController(atmSvc)

	atmController.Start()
}

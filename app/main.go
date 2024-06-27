package main

import (
	atm_handler "atm-simulation-console/internal/handler/atm"
	account_repository "atm-simulation-console/internal/repository/account"
	atm_service "atm-simulation-console/internal/service/atm"
)

func main() {
	accountRepo := account_repository.NewAccountRepository()
	atmSvc := atm_service.NewATMService(accountRepo)

	atmController := atm_handler.NewATMController(atmSvc)

	atmController.Start()
}

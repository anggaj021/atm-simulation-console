package main

import (
	account_repository "atm-simulation-console/internal/account/repository"
	atm_controller "atm-simulation-console/internal/atm/controller"
	atm_service "atm-simulation-console/internal/atm/service"
)

func main() {
	accountRepo := account_repository.NewAccountRepository()
	atmSvc := atm_service.NewATMService(accountRepo)

	atmController := atm_controller.NewATMController(atmSvc)

	atmController.Start()
}

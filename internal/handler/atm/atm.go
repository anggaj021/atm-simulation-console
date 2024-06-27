package atm_controller

import (
	account_repository "atm-simulation-console/internal/repository/account"
	atm_service "atm-simulation-console/internal/service/atm"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"

	"atm-simulation-console/internal/util/formatter"
	"atm-simulation-console/internal/util/generator"
)

type ATMData struct {
	AccNumber string
	AccDest   string
	Amount    string
	Ref       string
}

type ATMController struct {
	service *atm_service.ATMService
}

func NewATMController(svc *atm_service.ATMService) *ATMController {
	return &ATMController{
		service: svc,
	}
}

func (c *ATMController) Start() {
	c.initSampleAccounts()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("enter Account Number: ")
	accNumber := c.service.GetInputString(reader)

	fmt.Print("enter PIN: ")
	pin := c.service.GetInputString(reader)

	account, err := c.service.ValidateAccount(accNumber)
	if account == nil {
		formatter.ErrorMessage(err.Error())
		return
	}

	validated, err := c.service.ValidatePIN(account, pin)
	if validated == nil {
		formatter.ErrorMessage(err.Error())
		return
	}

	for {
		if !c.displayTrxScreen(reader, accNumber) {
			break
		}
	}
}

// ==================================== PROCESSOR ====================================

func (c *ATMController) processMainMenu(reader *bufio.Reader, accNumber string, option string) bool {
	switch option {
	case "1":
		return c.displayWithdrawScreen(reader, accNumber)
	case "2":
		return c.displayTrfDestNumScreen(reader, accNumber)
	case "3", "":
		formatter.ErrorMessage("exiting...")
		return false
	default:
		formatter.ErrorMessage("invalid option")
	}

	return true
}

func (c *ATMController) processWithdrawMenu(reader *bufio.Reader, accNumber string, option string) bool {
	switch option {
	case "1":
		amount := 10
		if ok := c.checkBalanceBoolResult(accNumber, amount); !ok {
			return false
		}
		c.service.Withdraw(accNumber, amount)
		return c.displayWdSummaryScreen(reader, accNumber, amount)

	case "2":
		amount := 50
		if ok := c.checkBalanceBoolResult(accNumber, amount); !ok {
			return false
		}
		c.service.Withdraw(accNumber, amount)
		return c.displayWdSummaryScreen(reader, accNumber, amount)
	case "3":
		amount := 100
		if ok := c.checkBalanceBoolResult(accNumber, amount); !ok {
			return false
		}
		c.service.Withdraw(accNumber, amount)
		return c.displayWdSummaryScreen(reader, accNumber, amount)
	case "4":
		return c.displayOtherWithdrawScreen(reader, accNumber)
	case "5", "":
		c.displayTrxScreen(reader, accNumber)
		return false
	default:
		formatter.ErrorMessage("invalid option")
	}
	return true
}

func (c *ATMController) processWdSummary(reader *bufio.Reader, accNumber string, option string) bool {
	switch option {
	case "1":
		c.displayWithdrawScreen(reader, accNumber)
	case "2":
		return false
	default:
		formatter.ErrorMessage("invalid option")
	}

	return true
}

func (c *ATMController) processTrfDestNumber(reader *bufio.Reader, accNumber string, val string) bool {
	if val == accNumber {
		formatter.ErrorMessage("invalid account number")
		return true
	}

	switch val {
	case "", "0":
		c.displayTrxScreen(reader, accNumber)
		return false
	default:
		account, err := c.service.ValidateAccount(val)
		if account == nil {
			formatter.ErrorMessage(err.Error())
			return true
		}
		detail := ATMData{
			AccNumber: accNumber,
			AccDest:   val,
		}
		return c.displayTrfAmountScreen(reader, detail)
	}
}

func (c *ATMController) processTrfAmount(reader *bufio.Reader, detail ATMData) bool {
	switch detail.Amount {
	case "", "0":
		c.displayTrxScreen(reader, detail.AccNumber)
		return false
	default:
		intAmount, err := strconv.Atoi(detail.Amount)
		if err != nil {
			formatter.ErrorMessage("invalid amount")
		}
		err = c.service.ValidateTransferAmount(detail.AccNumber, intAmount)
		if err != nil {
			formatter.ErrorMessage(err.Error())
			return true
		}
		return c.displayTransferConfirmScreen(reader, detail)
	}
}

func (c *ATMController) processTrfConfirm(reader *bufio.Reader, detail ATMData, option string) bool {
	switch option {
	case "1":
		intAmount, _ := strconv.Atoi(detail.Amount)
		err := c.service.Transfer(detail.AccNumber, detail.AccDest, intAmount)
		if err != nil {
			formatter.ErrorMessage(err.Error())
			return true
		}
		c.displayTransferSummaryScreen(reader, detail)
	case "2":
		return c.displayTrxScreen(reader, detail.AccNumber)
	default:
		formatter.ErrorMessage("invalid option")
	}

	return true
}

func (c *ATMController) processTrxSummary(reader *bufio.Reader, accNumber string, option string) bool {
	switch option {
	case "1":
		c.displayTrxScreen(reader, accNumber)
	case "2":
		return false
	default:
		formatter.ErrorMessage("invalid option")
	}

	return true
}

// ==================================== DISPLAY SCREEN ====================================

func (c *ATMController) displayTrxScreen(reader *bufio.Reader, accNumber string) bool {
	fmt.Println("1. Withdraw")
	fmt.Println("2. Fund Transfer")
	fmt.Println("3. Exit")
	fmt.Print("Please choose option[3]: ")

	option := c.service.GetInputString(reader)
	return c.processMainMenu(reader, accNumber, option)
}

func (c *ATMController) displayWithdrawScreen(reader *bufio.Reader, accNumber string) bool {
	fmt.Println("1. $10")
	fmt.Println("2. $50")
	fmt.Println("3. $100")
	fmt.Println("4. Other")
	fmt.Println("5. Back")
	fmt.Print("Please choose option[5]: ")

	option := c.service.GetInputString(reader)
	return c.processWithdrawMenu(reader, accNumber, option)
}

func (c *ATMController) displayOtherWithdrawScreen(reader *bufio.Reader, accNumber string) bool {
	fmt.Println("Other Withdraw")
	fmt.Print("Enter amount to withdraw: ")

	amount, err := c.service.GetInputNumber(reader)
	if err != nil {
		formatter.ErrorMessage(err.Error())
		return true
	}

	err = c.service.ValidateOtherWithdraw(accNumber, amount)
	if err != nil {
		formatter.ErrorMessage(err.Error())
		return true
	}

	c.service.Withdraw(accNumber, amount)
	return c.displayWdSummaryScreen(reader, accNumber, amount)
}

func (c *ATMController) displayWdSummaryScreen(reader *bufio.Reader, accNumber string, amount int) bool {

	time := formatter.DateFormatter(time.Now())
	balance := c.service.GetBalance(accNumber)

	fmt.Println("Summary")
	fmt.Println("Date		: " + time)
	fmt.Println("Withdraw	: " + strconv.Itoa(amount))
	fmt.Println("Balance	: " + strconv.Itoa(balance))
	fmt.Println("")
	fmt.Println("1. Transaction")
	fmt.Println("2. Exit")
	fmt.Print("Choose option[2]: ")

	option := c.service.GetInputString(reader)
	return c.processWdSummary(reader, accNumber, option)
}

func (c *ATMController) displayTrfDestNumScreen(reader *bufio.Reader, accNumber string) bool {

	fmt.Println("Please enter destination account")
	fmt.Println("or enter 0 to go back to Transaction")
	fmt.Print("Destination account[0]: ")

	accDest := c.service.GetInputString(reader)
	if accDest != "" {
		return c.processTrfDestNumber(reader, accNumber, accDest)
	}

	c.displayTrxScreen(reader, accNumber)
	return false
}

func (c *ATMController) displayTrfAmountScreen(reader *bufio.Reader, detail ATMData) bool {

	fmt.Println("Please enter transfer amount")
	fmt.Println("or enter 0 to go back to Transaction")
	fmt.Print("Transfer amount[0]: ")

	amount := c.service.GetInputString(reader)
	detail.Amount = amount
	if detail.AccDest != "" {
		return c.processTrfAmount(reader, detail)
	}

	c.displayTrxScreen(reader, detail.AccNumber)
	return false
}

func (c *ATMController) displayTransferConfirmScreen(reader *bufio.Reader, detail ATMData) bool {

	refNum := generator.GenerateRandomNDigitNumber(6)
	stringRef := strconv.Itoa(refNum)

	fmt.Println("Transfer Confirmation")
	fmt.Println("Destination Account : " + detail.AccDest)
	fmt.Println("Transfer Amount     : " + detail.Amount)
	fmt.Println("Reference Number    : " + stringRef)
	fmt.Println("")
	fmt.Println("1. Confirm Trx")
	fmt.Println("2. Cancel Trx")
	fmt.Print("Choose option[2]: ")

	option := c.service.GetInputString(reader)
	details := ATMData{
		AccNumber: detail.AccNumber,
		AccDest:   detail.AccDest,
		Amount:    detail.Amount,
		Ref:       stringRef,
	}
	return c.processTrfConfirm(reader, details, option)
}

func (c *ATMController) displayTransferSummaryScreen(reader *bufio.Reader, detail ATMData) bool {

	balance := c.service.GetBalance(detail.AccNumber)

	fmt.Println("Fund Transfer Summary")
	fmt.Println("Destination Account : " + detail.AccDest)
	fmt.Println("Transfer Amount     : " + detail.Amount)
	fmt.Println("Reference Number    : " + detail.Ref)
	fmt.Println("Balance             : " + strconv.Itoa(balance))
	fmt.Println("")
	fmt.Println("1. Transaction")
	fmt.Println("2. Exit")
	fmt.Print("Choose option[2]: ")

	option := c.service.GetInputString(reader)
	return c.processTrxSummary(reader, detail.AccNumber, option)
}

// ==================================== OTHER ====================================

func (c *ATMController) checkBalanceBoolResult(accNumber string, amount int) bool {
	err := c.service.CheckBalance(accNumber, amount)
	if err != nil {
		formatter.ErrorMessage(err.Error())
		return false
	}
	return true
}

// ==================================== ACCOUNT SEEDER ====================================

func (c *ATMController) initSampleAccounts() {
	// Add sample accounts
	account1 := account_repository.Account{
		AccountNumber: "112233",
		Name:          "John Doe",
		Pin:           "123123",
		Balance:       100,
	}
	account2 := account_repository.Account{
		AccountNumber: "112244",
		Name:          "Jane Doe",
		Pin:           "123123",
		Balance:       30,
	}

	c.service.AddAccount(account1)
	c.service.AddAccount(account2)
}

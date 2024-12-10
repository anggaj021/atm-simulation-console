package atm_service

import (
	account_repository "atm-simulation-console/internal/repository/account"
	transaction_repository "atm-simulation-console/internal/repository/transaction"
	"atm-simulation-console/internal/util/formatter"
	"atm-simulation-console/internal/util/generator"
	"bufio"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ATMService struct {
	accRepo account_repository.AccountRepository
	trxRepo transaction_repository.TransactionRepository
}

func NewATMService(
	accRepo account_repository.AccountRepository,
	trxRepo transaction_repository.TransactionRepository,
) *ATMService {
	return &ATMService{
		accRepo: accRepo,
		trxRepo: trxRepo,
	}
}

func (s *ATMService) AddAccount(account account_repository.Account) bool {
	return s.accRepo.Add(account)
}

func (s *ATMService) ValidateAccount(accNumber string) (*account_repository.Account, error) {
	if err := validateLength(accNumber, 6, "account number"); err != nil {
		return nil, err
	}
	if err := validateDigitsOnly(accNumber, "account number"); err != nil {
		return nil, err
	}

	acc := s.accRepo.Find(accNumber)
	if acc == nil {
		return nil, errors.New("invalid account number")
	}

	return acc, nil
}

func (s *ATMService) ValidatePIN(account *account_repository.Account, pin string) (*account_repository.Account, error) {
	if err := validateLength(pin, 6, "PIN"); err != nil {
		return nil, err
	}
	if err := validateDigitsOnly(pin, "PIN"); err != nil {
		return nil, err
	}

	if account.Pin != pin {
		return nil, errors.New("invalid account number/PIN")
	}

	return account, nil
}

func (s *ATMService) GetBalance(accNumber string) int {
	return s.accRepo.GetBalance(accNumber)
}

func (s *ATMService) CheckBalance(accNumber string, amount int) error {
	currentBalance := s.GetBalance(accNumber)
	if currentBalance < amount {
		return errors.New("insufficient balance " + "$" + strconv.Itoa(amount))
	}

	return nil
}

func (s *ATMService) ValidateOtherWithdraw(accNumber string, amount int) error {
	if amount%10 != 0 {
		return errors.New("invalid amount: must be a multiple of 10")
	}

	if amount > 1000 {
		return errors.New("maximum amount to withdraw is $1000")
	}

	currentBalance := s.GetBalance(accNumber)
	if currentBalance < amount {
		return errors.New("insufficient balance " + "$" + strconv.Itoa(amount))
	}

	return s.CheckBalance(accNumber, amount)
}

func (s *ATMService) ValidateTransferAmount(accNumber string, amount int) error {
	if amount < 0 {
		return errors.New("minimum amount to transfer is $1")
	}

	if amount > 1000 {
		return errors.New("maximum amount to transfer is $1000")
	}

	currentBalance := s.GetBalance(accNumber)
	if currentBalance < amount {
		return errors.New("insufficient balance " + "$" + strconv.Itoa(amount))
	}

	return s.CheckBalance(accNumber, amount)
}

func (s *ATMService) Withdraw(accNumber string, amount int) bool {
	success := s.accRepo.Withdraw(accNumber, amount)
	if success {
		trx := transaction_repository.Transaction{
			TransactionID: generateTrxNumber("TRX-WD"),
			SourceID:      "",
			DestinationID: accNumber,
			Type:          "withdraw",
			TrxType:       "db",
			Amount:        amount,
			Date:          formatter.DateFormatter(time.Now()),
		}

		return s.trxRepo.Store(trx)
	}

	return false
}

func (s *ATMService) Deposit(accNumber string, amount int) bool {
	success := s.accRepo.Deposit(accNumber, amount)
	if success {
		trx := transaction_repository.Transaction{
			TransactionID: generateTrxNumber("TRX-DP"),
			SourceID:      "",
			DestinationID: accNumber,
			Type:          "deposit",
			TrxType:       "cr",
			Amount:        amount,
			Date:          formatter.DateFormatter(time.Now()),
		}

		return s.trxRepo.Store(trx)
	}

	return false
}

func (s *ATMService) Transfer(srcNumber, destNumber string, amount int) error {

	destNum := s.accRepo.Find(destNumber)
	if destNum == nil || srcNumber == destNumber {
		return errors.New("invalid destination account")
	}

	transactions := []transaction_repository.Transaction{}

	if s.accRepo.Withdraw(srcNumber, amount) {
		trxNum := generateTrxNumber("TRX-TF")
		wd := transaction_repository.Transaction{
			TransactionID: trxNum,
			SourceID:      srcNumber,
			DestinationID: destNumber,
			Type:          "transfer",
			TrxType:       "cr",
			Amount:        amount,
			Date:          formatter.DateFormatter(time.Now()),
		}

		transactions = append(transactions, wd)

		success := s.accRepo.Deposit(destNumber, amount)
		if !success {
			return errors.New("error transfering balance")
		}

		dp := transaction_repository.Transaction{
			TransactionID: trxNum,
			SourceID:      destNumber,
			DestinationID: srcNumber,
			Type:          "transfer",
			TrxType:       "db",
			Amount:        amount,
			Date:          formatter.DateFormatter(time.Now()),
		}

		transactions = append(transactions, dp)

		for _, row := range transactions {
			s.trxRepo.Store(row)
		}

		return nil
	}
	return errors.New("insufficient balance " + "$" + strconv.Itoa(amount))
}

func (s *ATMService) GetInputNumber(reader *bufio.Reader) (int, error) {
	amountStr, err := reader.ReadString('\n')
	if err != nil {
		return 0, errors.New("invalid input")
	}
	amountStr = strings.TrimSpace(amountStr)

	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		return 0, errors.New("invalid input: please enter a valid number")
	}
	return amount, nil
}

func (s *ATMService) GetTransactionHistory(accNumber string) []transaction_repository.Transaction {
	return s.trxRepo.Get(accNumber, 10)
}

func (s *ATMService) GetInputString(reader *bufio.Reader) string {
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	return input
}

func validateLength(input string, length int, fieldName string) error {
	if len(input) != length {
		return errors.New(fieldName + " should have " + strconv.Itoa(length) + " digits length")
	}
	return nil
}

func validateDigitsOnly(input string, fieldName string) error {

	if matched, _ := regexp.MatchString(`^\d{`+strconv.Itoa(len(input))+`}$`, input); !matched {
		return errors.New(fieldName + " should only contain numbers")
	}
	return nil
}

func generateTrxNumber(prefix string) string {
	return prefix + "-" + strconv.Itoa(generator.GenerateRandomNDigitNumber(6))
}

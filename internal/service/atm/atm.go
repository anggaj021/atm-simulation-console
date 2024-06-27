package atm_service

import (
	account_repository "atm-simulation-console/internal/repository/account"
	"bufio"
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type ATMService struct {
	repo *account_repository.AccountRepository
}

func NewATMService(repo *account_repository.AccountRepository) *ATMService {
	return &ATMService{
		repo: repo,
	}
}

func (s *ATMService) AddAccount(account account_repository.Account) bool {
	return s.repo.AddAccount(account)
}

func (s *ATMService) ValidateAccount(accNumber string) (*account_repository.Account, error) {
	if err := validateLength(accNumber, 6, "account number"); err != nil {
		return nil, err
	}
	if err := validateDigitsOnly(accNumber, "account number"); err != nil {
		return nil, err
	}

	acc := s.repo.FindAccount(accNumber)
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
	return s.repo.GetBalance(accNumber)
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
	return s.repo.Withdraw(accNumber, amount)
}

func (s *ATMService) Deposit(accNumber string, amount int) bool {
	return s.repo.Deposit(accNumber, amount)
}

func (s *ATMService) Transfer(srcNumber, destNumber string, amount int) error {

	destNum := s.repo.FindAccount(destNumber)
	if destNum == nil || srcNumber == destNumber {
		return errors.New("invalid destination account")
	}

	if s.repo.Withdraw(srcNumber, amount) {
		s.repo.Deposit(destNumber, amount)
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

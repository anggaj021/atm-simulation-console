package account_csv

import (
	account_repository "atm-simulation-console/internal/repository/account"
	"encoding/csv"
	"os"
	"strconv"
)

type CSVAccountRepository struct {
	filePath string
}

func NewCSVAccountRepository(filePath string) *CSVAccountRepository {
	return &CSVAccountRepository{
		filePath: filePath,
	}
}

func (r *CSVAccountRepository) Add(account account_repository.Account) bool {
	file, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return false
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return writer.Write([]string{account.AccountNumber, account.Name, account.Pin, strconv.Itoa(account.Balance)}) == nil
}

func (r *CSVAccountRepository) Find(accountNumber string) *account_repository.Account {
	accounts := r.readAllAccounts()
	for _, account := range accounts {
		if account.AccountNumber == accountNumber {
			return &account
		}
	}
	return nil
}

func (r *CSVAccountRepository) GetBalance(accountNumber string) int {
	account := r.Find(accountNumber)
	if account != nil {
		return account.Balance
	}
	return 0
}

func (r *CSVAccountRepository) Withdraw(number string, amount int) bool {
	accounts := r.readAllAccounts()

	if len(accounts) <= 0 {
		return false
	}

	for i := range accounts {
		if accounts[i].AccountNumber == number {
			if accounts[i].Balance < amount {
				return false
			}
			accounts[i].Balance -= amount
		}
	}

	return r.writeAllAccounts(accounts)
}

func (r *CSVAccountRepository) Deposit(number string, amount int) bool {
	accounts := r.readAllAccounts()

	if len(accounts) <= 0 {
		return false
	}

	for i := range accounts {
		if accounts[i].AccountNumber == number {
			accounts[i].Balance += amount
		}
	}

	return r.writeAllAccounts(accounts)
}

func (r *CSVAccountRepository) readAllAccounts() []account_repository.Account {
	file, err := os.Open(r.filePath)
	if err != nil {
		return []account_repository.Account{}
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return []account_repository.Account{}
	}

	var accounts = []account_repository.Account{}
	for _, line := range lines {
		balance, _ := strconv.Atoi(line[3])
		account := account_repository.Account{
			AccountNumber: line[0],
			Name:          line[1],
			Pin:           line[2],
			Balance:       balance,
		}

		accounts = append(accounts, account)
	}

	return accounts
}

func (r *CSVAccountRepository) writeAllAccounts(accounts []account_repository.Account) bool {
	file, err := os.Create(r.filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, account := range accounts {
		err := writer.Write([]string{
			account.AccountNumber,
			account.Name,
			account.Pin,
			strconv.Itoa(account.Balance),
		})

		if err != nil {
			return false
		}

	}
	return true
}

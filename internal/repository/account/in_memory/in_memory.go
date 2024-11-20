package in_memory

import (
	account_repository "atm-simulation-console/internal/repository/account"
)

type InMemoryAccount struct {
	accounts map[string]account_repository.Account
}

func NewInMemoryAccount() *InMemoryAccount {
	return &InMemoryAccount{
		accounts: make(map[string]account_repository.Account),
	}
}

func (r *InMemoryAccount) Add(account account_repository.Account) bool {
	r.accounts[account.AccountNumber] = account
	return true
}

func (r *InMemoryAccount) Find(number string) *account_repository.Account {
	account, ok := r.accounts[number]
	if !ok {
		return nil
	}
	return &account
}

func (r *InMemoryAccount) GetBalance(number string) int {
	return r.accounts[number].Balance
}

func (r *InMemoryAccount) Withdraw(number string, amount int) bool {
	account, ok := r.accounts[number]
	if !ok || account.Balance < amount {
		return false
	}

	account.Balance -= amount
	r.accounts[number] = account
	return true
}

func (r *InMemoryAccount) Deposit(number string, amount int) bool {
	account := r.accounts[number]
	account.Balance += amount
	r.accounts[number] = account
	return true
}

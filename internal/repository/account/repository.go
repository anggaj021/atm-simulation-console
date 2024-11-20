package account_repository

// Account structure
type Account struct {
	AccountNumber string
	Name          string
	Pin           string
	Balance       int
}

type AccountRepository interface {
	Add(account Account) bool
	Find(number string) *Account
	GetBalance(number string) int
	Withdraw(number string, amount int) bool
	Deposit(number string, amount int) bool
}

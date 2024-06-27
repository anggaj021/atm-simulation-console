package account_repository

type Account struct {
	AccountNumber string
	Name          string
	Pin           string
	Balance       int
}

type AccountRepository struct {
	accounts map[string]Account
}

func NewAccountRepository() *AccountRepository {
	return &AccountRepository{
		accounts: make(map[string]Account),
	}
}

func (r *AccountRepository) AddAccount(account Account) bool {
	r.accounts[account.AccountNumber] = account
	return true
}

func (r *AccountRepository) FindAccount(number string) *Account {
	account, ok := r.accounts[number]
	if !ok {
		return nil
	}
	return &account
}

func (r *AccountRepository) GetBalance(number string) int {
	return r.accounts[number].Balance
}

func (r *AccountRepository) Withdraw(number string, amount int) bool {
	account, ok := r.accounts[number]
	if !ok || account.Balance < amount {
		return false
	}

	account.Balance -= amount
	r.accounts[number] = account
	return true
}

func (r *AccountRepository) Deposit(number string, amount int) bool {
	account := r.accounts[number]
	account.Balance += amount
	r.accounts[number] = account
	return true
}

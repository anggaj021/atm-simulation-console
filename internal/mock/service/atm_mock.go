package mock

// Define the Account struct and the AccountRepository interface here if not defined already

type Account struct {
	Number string
	Pin    string
}

type AccountRepository interface {
	AddAccount(account Account) bool
	FindAccount(accNumber string) *Account
	GetBalance(accNumber string) int
	Withdraw(accNumber string, amount int) bool
	Deposit(accNumber string, amount int) bool
}

type ATMService struct {
	repo AccountRepository
}

func NewATMService(repo AccountRepository) *ATMService {
	return &ATMService{
		repo: repo,
	}
}

// Implement methods of ATMService

// Define MockAccountRepository implementing AccountRepository for testing purposes
type MockAccountRepository struct{}

func (m MockAccountRepository) AddAccount(account Account) bool {
	return true
}

func (m MockAccountRepository) FindAccount(accNumber string) *Account {
	return &Account{}
}

func (m MockAccountRepository) GetBalance(accNumber string) int {
	return 1000
}

func (m MockAccountRepository) Withdraw(accNumber string, amount int) bool {
	return true
}

func (m MockAccountRepository) Deposit(accNumber string, amount int) bool {
	return true
}

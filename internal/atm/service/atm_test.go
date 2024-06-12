package atm_service

import (
	account_repository "atm-simulation-console/internal/account/repository"
	"bufio"
	"strings"
	"testing"
)

func TestAddAccount(t *testing.T) {
	repo := account_repository.NewAccountRepository()
	atmSvc := NewATMService(repo)

	// Add test account
	testAccount := account_repository.Account{
		AccountNumber: "123456",
		Pin:           "1234",
		Balance:       1000,
	}

	// Test adding account
	if !atmSvc.AddAccount(testAccount) {
		t.Error("Expected true for successful addition of account, got false")
	}
}

func TestValidateAccount(t *testing.T) {
	repo := account_repository.NewAccountRepository()
	atmSvc := NewATMService(repo)

	// Add test account
	testAccount := account_repository.Account{
		AccountNumber: "123456",
		Pin:           "111111",
		Balance:       1000,
	}
	repo.AddAccount(testAccount)

	// Test valid account
	validAccount, _ := atmSvc.ValidateAccount("123456")
	if validAccount == nil {
		t.Error("Expected valid account, got nil")
	}

	// Test invalid account : less than 6 digits
	acc, _ := atmSvc.ValidateAccount("123")
	if acc != nil {
		t.Error("Expected nil, got account")
	}

	// Test invalid account : invalid input
	acc, _ = atmSvc.ValidateAccount("aaabbb")
	if acc != nil {
		t.Error("Expected nil, got account")
	}

	// Test invalid account
	acc, _ = atmSvc.ValidateAccount("123111")
	if acc != nil {
		t.Error("Expected nil, got account")
	}
}

func TestValidatePIN(t *testing.T) {
	repo := account_repository.NewAccountRepository()
	atmSvc := NewATMService(repo)

	// Add test account
	testAccount := account_repository.Account{
		AccountNumber: "123456",
		Pin:           "111111",
		Balance:       1000,
	}
	repo.AddAccount(testAccount)

	// Test valid account
	validAccount, _ := atmSvc.ValidatePIN(&testAccount, "111111")
	if validAccount == nil {
		t.Error("Expected valid account, got nil")
	}

	// Test invalid account : less than 6 digits
	acc, _ := atmSvc.ValidatePIN(&testAccount, "123")
	if acc != nil {
		t.Error("Expected nil, got account")
	}

	// Test invalid account : invalid input
	acc, _ = atmSvc.ValidatePIN(&testAccount, "aaabbb")
	if acc != nil {
		t.Error("Expected nil, got account")
	}

	// Test invalid account
	acc, _ = atmSvc.ValidatePIN(&testAccount, "123111")
	if acc != nil {
		t.Error("Expected nil, got account")
	}
}

func TestTransfer(t *testing.T) {
	repo := account_repository.NewAccountRepository()
	atmSvc := NewATMService(repo)

	// Add test accounts
	srcAccount := account_repository.Account{
		AccountNumber: "123456",
		Pin:           "1234",
		Balance:       500,
	}
	destAccount := account_repository.Account{
		AccountNumber: "987654",
		Pin:           "5678",
		Balance:       2000,
	}
	repo.AddAccount(srcAccount)
	repo.AddAccount(destAccount)

	err := atmSvc.CheckBalance("123456", 500)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = atmSvc.CheckBalance("123456", 600)
	if err == nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = atmSvc.ValidateTransferAmount("123456", 500)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// balance less than 0
	err = atmSvc.ValidateTransferAmount("123456", -1)
	if err == nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// balance more than 1000
	err = atmSvc.ValidateTransferAmount("123456", 1500)
	if err == nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// insufficient balance
	err = atmSvc.ValidateTransferAmount("123456", 700)
	if err == nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test successful transfer
	err = atmSvc.Transfer("123456", "987654", 500)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if repo.GetBalance("123456") != 0 {
		t.Errorf("Expected balance of 0 for source account, got %d", repo.GetBalance("123456"))
	}
	if repo.GetBalance("987654") != 2500 {
		t.Errorf("Expected balance of 2500 for destination account, got %d", repo.GetBalance("987654"))
	}

	// Test failed transfer due to invalid destination account
	err = atmSvc.Transfer("123456", "999999", 500)
	if err == nil {
		t.Error("Expected false for failed transfer (invalid destination account), got true")
	}

	// Test failed transfer due to insufficient balance
	err = atmSvc.Transfer("123456", "987654", 1500)
	if err == nil {
		t.Errorf("Expected 'insufficient balance' error message, got %s", err)
	}
}

func TestValidateOtherWithdraw(t *testing.T) {
	repo := account_repository.NewAccountRepository()
	atmSvc := NewATMService(repo)

	// Add test accounts
	srcAccount := account_repository.Account{
		AccountNumber: "123456",
		Pin:           "1234",
		Balance:       500,
	}
	repo.AddAccount(srcAccount)

	err := atmSvc.ValidateOtherWithdraw("123456", 500)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = atmSvc.ValidateOtherWithdraw("123456", 600)
	if err == nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = atmSvc.ValidateOtherWithdraw("123456", 15)
	if err == nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = atmSvc.ValidateOtherWithdraw("123456", 1600)
	if err == nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGetInputNumber(t *testing.T) {
	repo := account_repository.NewAccountRepository()
	atmSvc := NewATMService(repo)
	tests := []struct {
		input     string
		expected  int
		expectErr bool
	}{
		{input: "100\n", expected: 100, expectErr: false},
		{input: "invalid\n", expected: 0, expectErr: true},
		{input: "invalid", expected: 0, expectErr: true},
	}

	for _, test := range tests {
		reader := bufio.NewReader(strings.NewReader(test.input))
		result, _ := atmSvc.GetInputNumber(reader)

		if result != test.expected && !test.expectErr {
			t.Errorf("For input %q, expected %d but got %d", test.input, test.expected, result)
		}
	}
}

func TestGetInputString(t *testing.T) {
	repo := account_repository.NewAccountRepository()
	atmSvc := NewATMService(repo)
	tests := []struct {
		input     string
		expected  string
		expectErr bool
	}{
		{input: "100\n", expected: "100", expectErr: false},
		{input: "invalid\n", expected: "invalid", expectErr: true},
	}

	for _, test := range tests {
		reader := bufio.NewReader(strings.NewReader(test.input))
		result := atmSvc.GetInputString(reader)

		if result != test.expected && !test.expectErr {
			t.Errorf("For input %q, expected %q but got %q", test.input, test.expected, result)
		}
	}
}

func TestGetBalance(t *testing.T) {
	repo := account_repository.NewAccountRepository()
	atmSvc := NewATMService(repo)

	// Add test account
	testAccount := account_repository.Account{
		AccountNumber: "123456",
		Pin:           "1234",
		Balance:       1000,
	}
	repo.AddAccount(testAccount)

	// Test getting balance
	balance := atmSvc.GetBalance("123456")
	if balance != 1000 {
		t.Errorf("Expected balance of 1000, got %d", balance)
	}
}

func TestWithdraw(t *testing.T) {
	repo := account_repository.NewAccountRepository()
	atmSvc := NewATMService(repo)

	// Add test account
	testAccount := account_repository.Account{
		AccountNumber: "123456",
		Pin:           "1234",
		Balance:       1000,
	}
	repo.AddAccount(testAccount)

	// Test successful withdrawal
	if !atmSvc.Withdraw("123456", 500) {
		t.Error("Expected true for successful withdrawal, got false")
	}
	if repo.GetBalance("123456") != 500 {
		t.Errorf("Expected balance of 500 after withdrawal, got %d", repo.GetBalance("123456"))
	}

	// Test failed withdrawal due to insufficient balance
	if atmSvc.Withdraw("123456", 600) {
		t.Error("Expected false for failed withdrawal due to insufficient balance, got true")
	}
}

func TestDeposit(t *testing.T) {
	repo := account_repository.NewAccountRepository()
	atmSvc := NewATMService(repo)

	// Add test account
	testAccount := account_repository.Account{
		AccountNumber: "123456",
		Pin:           "1234",
		Balance:       1000,
	}
	repo.AddAccount(testAccount)

	// Test successful deposit
	if !atmSvc.Deposit("123456", 500) {
		t.Error("Expected true for successful deposit, got false")
	}
	if repo.GetBalance("123456") != 1500 {
		t.Errorf("Expected balance of 1500 after deposit, got %d", repo.GetBalance("123456"))
	}
}

package account_repository

import "testing"

func TestFindAccount(t *testing.T) {
	repo := NewAccountRepository()
	repo.AddAccount(Account{
		AccountNumber: "123456",
		Pin:           "1234",
		Balance:       10000,
	})

	if repo.FindAccount("123456") == nil {
		t.Errorf("expected account, got nil")
	}
	if repo.FindAccount("123451") != nil {
		t.Errorf("expected nil, got account")
	}
}

func TestGetBalance(t *testing.T) {
	repo := NewAccountRepository()
	repo.AddAccount(Account{
		AccountNumber: "123456",
		Pin:           "1234",
		Balance:       10000,
	})

	if repo.GetBalance("123456") != 10000 {
		t.Errorf("expected 10000, got %d", repo.GetBalance("123456"))
	}
}

func TestWithdraw(t *testing.T) {
	repo := NewAccountRepository()
	repo.AddAccount(Account{
		AccountNumber: "123456",
		Pin:           "1234",
		Balance:       10000,
	})

	if !repo.Withdraw("123456", 5000) {
		t.Errorf("expected true, got false")
	}
	if repo.GetBalance("123456") != 5000 {
		t.Errorf("expected 5000, got %d", repo.GetBalance("123456"))
	}
	if repo.Withdraw("123456", 6000) {
		t.Errorf("expected false, got true")
	}
}

func TestDeposit(t *testing.T) {
	repo := NewAccountRepository()
	repo.AddAccount(Account{
		AccountNumber: "123456",
		Pin:           "1234",
		Balance:       10000,
	})

	repo.Deposit("123456", 5000)
	if repo.GetBalance("123456") != 15000 {
		t.Errorf("expected 15000, got %d", repo.GetBalance("123456"))
	}
}

package account_csv

import (
	account_repository "atm-simulation-console/internal/repository/account"
	"os"
	"reflect"
	"testing"
)

func setupTempFile(t *testing.T, initialData string) (string, func()) {
	tmpFile, err := os.CreateTemp("", "test_repo_acc_*.csv")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	if initialData != "" {
		if _, err := tmpFile.Write([]byte(initialData)); err != nil {
			t.Fatalf("failed to write initial data to temp file: %v", err)
		}
	}

	tmpFile.Close()

	return tmpFile.Name(), func() {
		os.Remove(tmpFile.Name())
	}
}

func TestNewCSVAccountRepository(t *testing.T) {
	filePath := "testpath.csv"

	repo := NewCSVAccountRepository(filePath)

	if repo == nil {
		t.Fatalf("NewCSVAccountRepository() returned nil")
	}
	if repo.filePath != filePath {
		t.Errorf("NewCSVAccountRepository().filePath = %v, want %v", repo.filePath, filePath)
	}
}

func TestCSVAccountRepository_Add(t *testing.T) {
	type args struct {
		account account_repository.Account
	}

	tests := []struct {
		name      string
		setupFile func() (string, func())
		args      args
		want      bool
	}{
		{
			name: "successfully add account to csv",
			setupFile: func() (string, func()) {
				return setupTempFile(t, "")
			},
			args: args{
				account_repository.Account{
					AccountNumber: "123123",
					Name:          "user test",
					Pin:           "123123",
					Balance:       1000,
				},
			},
			want: true,
		},
		{
			name: "failed to open file due to incorrect filepath",
			setupFile: func() (string, func()) {
				return "/tmp/err/user_add.csv", func() {}
			},
			args: args{
				account_repository.Account{
					AccountNumber: "123123",
					Name:          "user test",
					Pin:           "123123",
					Balance:       1000,
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath, cleanup := tt.setupFile()
			defer cleanup()

			r := &CSVAccountRepository{
				filePath: filePath,
			}

			if got := r.Add(tt.args.account); got != tt.want {
				t.Errorf("CSVAccountRepository.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVAccountRepository_Find(t *testing.T) {
	type args struct {
		accountNumber string
	}
	tests := []struct {
		name      string
		setupFile func() (string, func())
		args      args
		want      *account_repository.Account
	}{
		{
			name: "successfully find user in the csv",
			setupFile: func() (string, func()) {
				return setupTempFile(t, "112233,user test,123123,1000\n")
			},
			args: args{
				accountNumber: "112233",
			},
			want: &account_repository.Account{
				AccountNumber: "112233",
				Name:          "user test",
				Pin:           "123123",
				Balance:       1000,
			},
		},
		{
			name: "failed to find user listed in the csv",
			setupFile: func() (string, func()) {
				return "/tmp/err/user_add.csv", func() {}
			},
			args: args{
				accountNumber: "123123",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath, cleanup := tt.setupFile()
			defer cleanup()

			r := &CSVAccountRepository{
				filePath: filePath,
			}
			if got := r.Find(tt.args.accountNumber); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CSVAccountRepository.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVAccountRepository_GetBalance(t *testing.T) {
	type args struct {
		accountNumber string
	}
	tests := []struct {
		name      string
		setupFile func() (string, func())
		args      args
		want      int
	}{
		{
			name: "successfully get balance from listed user in csv",
			setupFile: func() (string, func()) {
				return setupTempFile(t, "112233,user test,123123,1000\n")
			},
			args: args{
				accountNumber: "112233",
			},
			want: 1000,
		},
		{
			name: "failed get balance user",
			setupFile: func() (string, func()) {
				return setupTempFile(t, "")
			},
			args: args{
				accountNumber: "123123",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath, cleanup := tt.setupFile()
			defer cleanup()
			r := &CSVAccountRepository{
				filePath: filePath,
			}
			if got := r.GetBalance(tt.args.accountNumber); got != tt.want {
				t.Errorf("CSVAccountRepository.GetBalance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVAccountRepository_Withdraw(t *testing.T) {
	type args struct {
		number string
		amount int
	}
	tests := []struct {
		name      string
		setupFile func() (string, func())
		args      args
		want      bool
	}{
		{
			name: "successfully withdraw",
			setupFile: func() (string, func()) {
				return setupTempFile(t, "112233,user test,123123,1000\n")
			},
			args: args{
				number: "112233",
				amount: 10,
			},
			want: true,
		},
		{
			name: "failed to withdraw due to insufficient balance",
			setupFile: func() (string, func()) {
				return setupTempFile(t, "223344,user test,123123,1000\n")
			},
			args: args{
				number: "223344",
				amount: 1010,
			},
			want: false,
		},
		{
			name: "failed to withdraw due to invalid account",
			setupFile: func() (string, func()) {
				return setupTempFile(t, "")
			},
			args: args{
				number: "223345",
				amount: 1010,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath, cleanup := tt.setupFile()
			defer cleanup()

			r := &CSVAccountRepository{
				filePath: filePath,
			}
			if got := r.Withdraw(tt.args.number, tt.args.amount); got != tt.want {
				t.Errorf("CSVAccountRepository.Withdraw() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVAccountRepository_Deposit(t *testing.T) {
	type args struct {
		number string
		amount int
	}
	tests := []struct {
		name      string
		setupFile func() (string, func())
		args      args
		want      bool
	}{
		{
			name: "successfully deposit",
			setupFile: func() (string, func()) {
				return setupTempFile(t, "112233,user test,123123,1000\n")
			},
			args: args{
				number: "112233",
				amount: 10,
			},
			want: true,
		},
		{
			name: "failed deposit due to invalid account",
			setupFile: func() (string, func()) {
				return setupTempFile(t, "")
			},
			args: args{
				number: "123123",
				amount: 10,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath, cleanup := tt.setupFile()
			defer cleanup()

			r := &CSVAccountRepository{
				filePath: filePath,
			}
			if got := r.Deposit(tt.args.number, tt.args.amount); got != tt.want {
				t.Errorf("CSVAccountRepository.Deposit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVAccountRepository_ReadAllAccounts(t *testing.T) {
	tests := []struct {
		name      string
		setupFile func() (string, func())
		want      []account_repository.Account
	}{
		{
			name: "valid accounts file",
			setupFile: func() (string, func()) {
				return setupTempFile(t, "112233,user test,123123,1000\n445566,another user,321321,2000\n")
			},
			want: []account_repository.Account{
				{AccountNumber: "112233", Name: "user test", Pin: "123123", Balance: 1000},
				{AccountNumber: "445566", Name: "another user", Pin: "321321", Balance: 2000},
			},
		},
		{
			name: "empty file",
			setupFile: func() (string, func()) {
				return setupTempFile(t, "")
			},
			want: []account_repository.Account{},
		},
		{
			name: "file does not exist",
			setupFile: func() (string, func()) {
				return "nonexistent.csv", func() {}
			},
			want: []account_repository.Account{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath, cleanup := tt.setupFile()
			defer cleanup()

			r := &CSVAccountRepository{
				filePath: filePath,
			}
			if got := r.ReadAllAccounts(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CSVAccountRepository.ReadAllAccounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

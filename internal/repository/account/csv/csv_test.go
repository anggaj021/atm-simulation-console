package account_csv

import (
	account_repository "atm-simulation-console/internal/repository/account"
	"os"
	"testing"
)

func TestCSVAccountRepository_Add(t *testing.T) {
	type fields struct {
		filePath string
	}
	type args struct {
		account account_repository.Account
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "successfully add account to csv",
			fields: fields{
				filePath: "",
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
		{
			name: "failed to open file due to incorrect filepath",
			fields: fields{
				filePath: "/tmp/err_user_add.csv",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.fields.filePath != "" {
				fileTmp, err := os.CreateTemp("", "test_add_account.csv")

				if err != nil {
					t.Errorf("failed to create temporary csv file")
				}

				defer os.Remove(fileTmp.Name())
				tt.fields.filePath = fileTmp.Name()
				fileTmp.Close()
			}

			r := &CSVAccountRepository{
				filePath: tt.fields.filePath,
			}

			if got := r.Add(tt.args.account); got != tt.want {
				t.Errorf("CSVAccountRepository.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestCSVAccountRepository_Find(t *testing.T) {
// 	type fields struct {
// 		filePath string
// 	}
// 	type args struct {
// 		accountNumber string
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   *account_repository.Account
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			r := &CSVAccountRepository{
// 				filePath: tt.fields.filePath,
// 			}
// 			if got := r.Find(tt.args.accountNumber); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("CSVAccountRepository.Find() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestCSVAccountRepository_GetBalance(t *testing.T) {
// 	type fields struct {
// 		filePath string
// 	}
// 	type args struct {
// 		accountNumber string
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   int
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			r := &CSVAccountRepository{
// 				filePath: tt.fields.filePath,
// 			}
// 			if got := r.GetBalance(tt.args.accountNumber); got != tt.want {
// 				t.Errorf("CSVAccountRepository.GetBalance() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestCSVAccountRepository_Withdraw(t *testing.T) {
// 	type fields struct {
// 		filePath string
// 	}
// 	type args struct {
// 		number string
// 		amount int
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			r := &CSVAccountRepository{
// 				filePath: tt.fields.filePath,
// 			}
// 			if got := r.Withdraw(tt.args.number, tt.args.amount); got != tt.want {
// 				t.Errorf("CSVAccountRepository.Withdraw() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestCSVAccountRepository_Deposit(t *testing.T) {
// 	type fields struct {
// 		filePath string
// 	}
// 	type args struct {
// 		number string
// 		amount int
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			r := &CSVAccountRepository{
// 				filePath: tt.fields.filePath,
// 			}
// 			if got := r.Deposit(tt.args.number, tt.args.amount); got != tt.want {
// 				t.Errorf("CSVAccountRepository.Deposit() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestCSVAccountRepository_readAllAccounts(t *testing.T) {
// 	type fields struct {
// 		filePath string
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		want   []account_repository.Account
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			r := &CSVAccountRepository{
// 				filePath: tt.fields.filePath,
// 			}
// 			if got := r.readAllAccounts(); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("CSVAccountRepository.readAllAccounts() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestCSVAccountRepository_writeAllAccounts(t *testing.T) {
// 	type fields struct {
// 		filePath string
// 	}
// 	type args struct {
// 		accounts []account_repository.Account
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			r := &CSVAccountRepository{
// 				filePath: tt.fields.filePath,
// 			}
// 			if got := r.writeAllAccounts(tt.args.accounts); got != tt.want {
// 				t.Errorf("CSVAccountRepository.writeAllAccounts() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

package atm_service

import (
	account_repository "atm-simulation-console/internal/repository/account"
	"atm-simulation-console/internal/repository/account/in_memory"
	transaction_repository "atm-simulation-console/internal/repository/transaction"
	transaction_csv "atm-simulation-console/internal/repository/transaction/csv"
	"bufio"
	"strings"
	"testing"
)

type MockAccountRepository struct {
	AddFn        func(account account_repository.Account) bool
	FindFn       func(number string) *account_repository.Account
	GetBalanceFn func(number string) int
	WithdrawFn   func(number string, amount int) bool
	DepositFn    func(number string, amount int) bool
}

func (m *MockAccountRepository) Add(account account_repository.Account) bool {
	return m.AddFn(account)
}

func (m *MockAccountRepository) Find(accountNumber string) *account_repository.Account {
	return m.FindFn(accountNumber)
}

func (m *MockAccountRepository) GetBalance(number string) int {
	return m.GetBalanceFn(number)
}

func (m *MockAccountRepository) Withdraw(accountNumber string, amount int) bool {
	return m.WithdrawFn(accountNumber, amount)
}

func (m *MockAccountRepository) Deposit(accountNumber string, amount int) bool {
	return m.DepositFn(accountNumber, amount)
}

type MockTransactionRepository struct {
	StoreFn func(transaction transaction_repository.Transaction) bool
	GetFn   func(userID string, limit int) []transaction_repository.Transaction
}

func (m *MockTransactionRepository) Store(transaction transaction_repository.Transaction) bool {
	return m.StoreFn(transaction)
}

func (m *MockTransactionRepository) Get(userID string, limit int) []transaction_repository.Transaction {
	return m.GetFn(userID, limit)
}

func TestAddAccount(t *testing.T) {
	repo := in_memory.NewInMemoryAccount()
	trxRepo := transaction_csv.NewCSVTransactionRepository("test")

	atmSvc := NewATMService(repo, trxRepo)

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
	repo := in_memory.NewInMemoryAccount()
	trxRepo := transaction_csv.NewCSVTransactionRepository("test")
	atmSvc := NewATMService(repo, trxRepo)

	// Add test account
	testAccount := account_repository.Account{
		AccountNumber: "123456",
		Pin:           "111111",
		Balance:       1000,
	}
	repo.Add(testAccount)

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
	repo := in_memory.NewInMemoryAccount()
	trxRepo := transaction_csv.NewCSVTransactionRepository("test")

	atmSvc := NewATMService(repo, trxRepo)

	// Add test account
	testAccount := account_repository.Account{
		AccountNumber: "123456",
		Pin:           "111111",
		Balance:       1000,
	}
	repo.Add(testAccount)

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

func TestValidateOtherWithdraw(t *testing.T) {
	repo := in_memory.NewInMemoryAccount()
	trxRepo := transaction_csv.NewCSVTransactionRepository("test")

	atmSvc := NewATMService(repo, trxRepo)

	// Add test accounts
	srcAccount := account_repository.Account{
		AccountNumber: "123456",
		Pin:           "1234",
		Balance:       500,
	}
	repo.Add(srcAccount)

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
	repo := in_memory.NewInMemoryAccount()
	trxRepo := transaction_csv.NewCSVTransactionRepository("test")

	atmSvc := NewATMService(repo, trxRepo)
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
	repo := in_memory.NewInMemoryAccount()
	trxRepo := transaction_csv.NewCSVTransactionRepository("test")

	atmSvc := NewATMService(repo, trxRepo)
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
	repo := in_memory.NewInMemoryAccount()
	trxRepo := transaction_csv.NewCSVTransactionRepository("test")

	atmSvc := NewATMService(repo, trxRepo)

	// Add test account
	testAccount := account_repository.Account{
		AccountNumber: "123456",
		Pin:           "1234",
		Balance:       1000,
	}
	repo.Add(testAccount)

	// Test getting balance
	balance := atmSvc.GetBalance("123456")
	if balance != 1000 {
		t.Errorf("Expected balance of 1000, got %d", balance)
	}
}

func TestATMService_Transfer(t *testing.T) {
	type fields struct {
		accRepo *MockAccountRepository
		trxRepo *MockTransactionRepository
	}
	type args struct {
		srcNumber  string
		destNumber string
		amount     int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Successful transfer",
			fields: fields{
				accRepo: &MockAccountRepository{
					FindFn: func(accountNumber string) *account_repository.Account {
						if accountNumber == "654321" {
							return &account_repository.Account{}
						}
						return nil
					},
					WithdrawFn: func(accountNumber string, amount int) bool {
						return true
					},
					DepositFn: func(accountNumber string, amount int) bool {
						return true
					},
				},
				trxRepo: &MockTransactionRepository{
					StoreFn: func(transaction transaction_repository.Transaction) bool {
						return true
					},
					GetFn: func(userID string, limit int) []transaction_repository.Transaction {
						return []transaction_repository.Transaction{}
					},
				},
			},
			args: args{
				srcNumber:  "123456",
				destNumber: "654321",
				amount:     100,
			},
			wantErr: false,
		},
		{
			name: "Invalid destination account",
			fields: fields{
				accRepo: &MockAccountRepository{
					FindFn: func(accountNumber string) *account_repository.Account {
						return nil
					},
				},
				trxRepo: &MockTransactionRepository{},
			},
			args: args{
				srcNumber:  "123456",
				destNumber: "654321",
				amount:     100,
			},
			wantErr: true,
		},
		{
			name: "Insufficient balance",
			fields: fields{
				accRepo: &MockAccountRepository{
					FindFn: func(accountNumber string) *account_repository.Account {
						return &account_repository.Account{}
					},
					WithdrawFn: func(accountNumber string, amount int) bool {
						return false
					},
				},
				trxRepo: &MockTransactionRepository{},
			},
			args: args{
				srcNumber:  "123456",
				destNumber: "654321",
				amount:     100,
			},
			wantErr: true,
		},
		{
			name: "Deposit failure",
			fields: fields{
				accRepo: &MockAccountRepository{
					FindFn: func(accountNumber string) *account_repository.Account {
						return &account_repository.Account{}
					},
					WithdrawFn: func(accountNumber string, amount int) bool {
						return true
					},
					DepositFn: func(accountNumber string, amount int) bool {
						return false
					},
				},
				trxRepo: &MockTransactionRepository{
					GetFn: func(userID string, limit int) []transaction_repository.Transaction {
						return []transaction_repository.Transaction{}
					},
				},
			},
			args: args{
				srcNumber:  "123456",
				destNumber: "654321",
				amount:     100,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ATMService{
				accRepo: tt.fields.accRepo,
				trxRepo: tt.fields.trxRepo,
			}
			if err := s.Transfer(tt.args.srcNumber, tt.args.destNumber, tt.args.amount); (err != nil) != tt.wantErr {
				t.Errorf("ATMService.Transfer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestATMService_ValidateTransferAmount(t *testing.T) {
	type fields struct {
		accRepo *MockAccountRepository
		trxRepo *MockTransactionRepository
	}
	type args struct {
		accNumber string
		amount    int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Minimum amount error",
			fields: fields{},
			args: args{
				accNumber: "123456",
				amount:    -5,
			},
			wantErr: true,
		},
		{
			name:   "Maximum amount error",
			fields: fields{},
			args: args{
				accNumber: "123456",
				amount:    1500,
			},
			wantErr: true,
		},
		{
			name: "Insufficient balance",
			fields: fields{
				accRepo: &MockAccountRepository{
					FindFn: func(accountNumber string) *account_repository.Account {
						return &account_repository.Account{}
					},
					GetBalanceFn: func(number string) int {
						return 100
					},
				},
			},
			args: args{
				accNumber: "123456",
				amount:    500,
			},
			wantErr: true,
		},
		{
			name: "Successful validation",
			fields: fields{
				accRepo: &MockAccountRepository{
					FindFn: func(accountNumber string) *account_repository.Account {
						return &account_repository.Account{}
					},
					GetBalanceFn: func(number string) int {
						return 500
					},
				},
			},
			args: args{
				accNumber: "123456",
				amount:    200,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ATMService{
				accRepo: tt.fields.accRepo,
				trxRepo: tt.fields.trxRepo,
			}
			if err := s.ValidateTransferAmount(tt.args.accNumber, tt.args.amount); (err != nil) != tt.wantErr {
				t.Errorf("ATMService.ValidateTransferAmount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestATMService_Withdraw(t *testing.T) {
	type fields struct {
		accRepo *MockAccountRepository
		trxRepo *MockTransactionRepository
	}
	type args struct {
		accNumber string
		amount    int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Successful withdrawal",
			fields: fields{
				accRepo: &MockAccountRepository{
					WithdrawFn: func(accountNumber string, amount int) bool {
						return true
					},
				},
				trxRepo: &MockTransactionRepository{
					StoreFn: func(transaction transaction_repository.Transaction) bool {
						return true
					},
				},
			},
			args: args{
				accNumber: "123456",
				amount:    200,
			},
			want: true,
		},
		{
			name: "Withdrawal failure due to account repository error",
			fields: fields{
				accRepo: &MockAccountRepository{
					WithdrawFn: func(accountNumber string, amount int) bool {
						return false
					},
				},
				trxRepo: &MockTransactionRepository{},
			},
			args: args{
				accNumber: "123456",
				amount:    200,
			},
			want: false,
		},
		{
			name: "Withdrawal failure due to transaction storage error",
			fields: fields{
				accRepo: &MockAccountRepository{
					WithdrawFn: func(accountNumber string, amount int) bool {
						return true
					},
				},
				trxRepo: &MockTransactionRepository{
					StoreFn: func(transaction transaction_repository.Transaction) bool {
						return false
					},
				},
			},
			args: args{
				accNumber: "123456",
				amount:    200,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ATMService{
				accRepo: tt.fields.accRepo,
				trxRepo: tt.fields.trxRepo,
			}
			if got := s.Withdraw(tt.args.accNumber, tt.args.amount); got != tt.want {
				t.Errorf("ATMService.Withdraw() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestATMService_Deposit(t *testing.T) {
	type fields struct {
		accRepo *MockAccountRepository
		trxRepo *MockTransactionRepository
	}
	type args struct {
		accNumber string
		amount    int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Successful deposit",
			fields: fields{
				accRepo: &MockAccountRepository{
					DepositFn: func(accountNumber string, amount int) bool {
						return true
					},
				},
				trxRepo: &MockTransactionRepository{
					StoreFn: func(transaction transaction_repository.Transaction) bool {
						return true
					},
				},
			},
			args: args{
				accNumber: "123456",
				amount:    200,
			},
			want: true,
		},
		{
			name: "Deposit failure due to account repository error",
			fields: fields{
				accRepo: &MockAccountRepository{
					DepositFn: func(accountNumber string, amount int) bool {
						return false
					},
				},
				trxRepo: &MockTransactionRepository{},
			},
			args: args{
				accNumber: "123456",
				amount:    200,
			},
			want: false,
		},
		{
			name: "Deposit failure due to transaction storage error",
			fields: fields{
				accRepo: &MockAccountRepository{
					DepositFn: func(accountNumber string, amount int) bool {
						return true
					},
				},
				trxRepo: &MockTransactionRepository{
					StoreFn: func(transaction transaction_repository.Transaction) bool {
						return false
					},
				},
			},
			args: args{
				accNumber: "123456",
				amount:    200,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ATMService{
				accRepo: tt.fields.accRepo,
				trxRepo: tt.fields.trxRepo,
			}
			if got := s.Deposit(tt.args.accNumber, tt.args.amount); got != tt.want {
				t.Errorf("ATMService.Deposit() = %v, want %v", got, tt.want)
			}
		})
	}
}

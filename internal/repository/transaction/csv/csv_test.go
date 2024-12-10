package transaction_csv

import (
	transaction_repository "atm-simulation-console/internal/repository/transaction"
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

func TestNewCSVTransactionRepository(t *testing.T) {
	filePath := "testpath.csv"

	repo := NewCSVTransactionRepository(filePath)

	if repo == nil {
		t.Fatalf("NewCSVTransactionRepository() returned nil")
	}
	if repo.filePath != filePath {
		t.Errorf("NewCSVTransactionRepository().filePath = %v, want %v", repo.filePath, filePath)
	}
}

func TestCSVTransactionRepository_ReadAllTransaction(t *testing.T) {
	tests := []struct {
		name      string
		setupFile func() (string, func())
		want      []transaction_repository.Transaction
	}{
		{
			name: "valid accounts file",
			setupFile: func() (string, func()) {
				return setupTempFile(t, "123123,123123,222222,deposit,cr,100,2024-12-10\n")
			},
			want: []transaction_repository.Transaction{
				{
					TransactionID: "123123",
					SourceID:      "123123",
					DestinationID: "222222",
					Type:          "deposit",
					TrxType:       "cr",
					Amount:        100,
					Date:          "2024-12-10",
				},
			},
		},
		{
			name: "empty file",
			setupFile: func() (string, func()) {
				return setupTempFile(t, "")
			},
			want: []transaction_repository.Transaction{},
		},
		{
			name: "file does not exist",
			setupFile: func() (string, func()) {
				return "notexists.csv", func() {}
			},
			want: []transaction_repository.Transaction{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath, cleanup := tt.setupFile()
			defer cleanup()
			r := &CSVTransactionRepository{
				filePath: filePath,
			}
			if got := r.ReadAllTransaction(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CSVTransactionRepository.ReadAllTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVTransactionRepository_WriteAllTransaction(t *testing.T) {
	type args struct {
		transactions []transaction_repository.Transaction
	}
	tests := []struct {
		name      string
		setupFile func() (string, func())
		args      args
		want      bool
	}{
		{
			name: "success write file",
			setupFile: func() (string, func()) {
				return setupTempFile(t, "")
			},
			args: args{
				[]transaction_repository.Transaction{
					{
						TransactionID: "123123",
						SourceID:      "123123",
						DestinationID: "222222",
						Type:          "deposit",
						TrxType:       "cr",
						Amount:        100,
						Date:          "2024-12-10",
					},
				},
			},
			want: true,
		},
		{
			name: "file does not exist",
			setupFile: func() (string, func()) {
				return "notexistdir/file.csv", func() {}
			},
			args: args{
				[]transaction_repository.Transaction{
					{
						TransactionID: "123123",
						SourceID:      "123123",
						DestinationID: "222222",
						Type:          "deposit",
						TrxType:       "cr",
						Amount:        100,
						Date:          "2024-12-10",
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath, cleanup := tt.setupFile()
			defer cleanup()
			r := &CSVTransactionRepository{
				filePath: filePath,
			}
			if got := r.WriteAllTransaction(tt.args.transactions); got != tt.want {
				t.Errorf("CSVTransactionRepository.WriteAllTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVTransactionRepository_Store(t *testing.T) {
	type args struct {
		transaction transaction_repository.Transaction
	}
	tests := []struct {
		name      string
		setupFile func() (string, func())
		args      args
		want      bool
	}{
		{
			name: "success store record",
			setupFile: func() (string, func()) {
				return setupTempFile(t, "123123,123123,222222,deposit,cr,100,2024-12-10\n")
			},
			args: args{
				transaction_repository.Transaction{
					TransactionID: "123123",
					SourceID:      "123123",
					DestinationID: "222222",
					Type:          "deposit",
					TrxType:       "cr",
					Amount:        100,
					Date:          "2024-12-10",
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath, cleanup := tt.setupFile()
			defer cleanup()
			r := &CSVTransactionRepository{
				filePath: filePath,
			}
			if got := r.Store(tt.args.transaction); got != tt.want {
				t.Errorf("CSVTransactionRepository.Store() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVTransactionRepository_Get(t *testing.T) {
	type args struct {
		accNumber string
		limit     int
	}
	tests := []struct {
		name      string
		setupFile func() (string, func())
		args      args
		want      []transaction_repository.Transaction
	}{
		{
			name: "success get record",
			setupFile: func() (string, func()) {
				return setupTempFile(t, "123123,123123,222222,deposit,cr,100,2024-12-10\n123123,123123,222222,transfer,cr,100,2024-12-10\n123123,123123,222222,deposit,cr,100,2024-12-10")
			},
			args: args{
				accNumber: "222222",
				limit:     2,
			},
			want: []transaction_repository.Transaction{
				{
					TransactionID: "123123",
					SourceID:      "123123",
					DestinationID: "222222",
					Type:          "deposit",
					TrxType:       "cr",
					Amount:        100,
					Date:          "2024-12-10",
				},
				{
					TransactionID: "123123",
					SourceID:      "123123",
					DestinationID: "222222",
					Type:          "transfer",
					TrxType:       "cr",
					Amount:        100,
					Date:          "2024-12-10",
				},
			},
		},
		{
			name: "no record found",
			setupFile: func() (string, func()) {
				return setupTempFile(t, "")
			},
			args: args{
				accNumber: "222222",
				limit:     2,
			},
			want: []transaction_repository.Transaction{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath, cleanup := tt.setupFile()
			defer cleanup()
			r := &CSVTransactionRepository{
				filePath: filePath,
			}
			if got := r.Get(tt.args.accNumber, tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CSVTransactionRepository.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

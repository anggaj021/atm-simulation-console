package transaction_csv

import (
	transaction_repository "atm-simulation-console/internal/repository/transaction"
	"encoding/csv"
	"os"
	"strconv"
)

type CSVTransactionRepository struct {
	filePath string
}

func NewCSVTransactionRepository(filePath string) *CSVTransactionRepository {
	return &CSVTransactionRepository{
		filePath: filePath,
	}
}

func (r *CSVTransactionRepository) Get(accNumber string, limit int) []transaction_repository.Transaction {
	history := r.ReadAllTransaction()
	trxHistory := []transaction_repository.Transaction{}
	counter := 0
	for _, row := range history {
		switch row.Type {
		case "deposit", "withdraw":
			if row.DestinationID == accNumber {
				trxHistory = append(trxHistory, row)
				counter++
			}
		case "transfer":
			if row.DestinationID == accNumber || row.SourceID == accNumber {
				trxHistory = append(trxHistory, row)
				counter++
			}
		}
		if counter >= limit {
			break
		}
	}

	return trxHistory
}

func (r *CSVTransactionRepository) Store(transaction transaction_repository.Transaction) bool {
	history := r.ReadAllTransaction()
	history = append(history[:1], append([]transaction_repository.Transaction{transaction}, history[1:]...)...)

	return r.WriteAllTransaction(history)
}

func (r *CSVTransactionRepository) ReadAllTransaction() []transaction_repository.Transaction {
	file, err := os.Open(r.filePath)
	if err != nil {
		return []transaction_repository.Transaction{}
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return []transaction_repository.Transaction{}
	}

	var transactions = []transaction_repository.Transaction{}
	for _, line := range lines {
		amount, _ := strconv.Atoi(line[5])
		transaction := transaction_repository.Transaction{
			TransactionID: line[0],
			SourceID:      line[1],
			DestinationID: line[2],
			Type:          line[3],
			TrxType:       line[4],
			Amount:        amount,
			Date:          line[6],
		}

		transactions = append(transactions, transaction)
	}

	return transactions
}

func (r *CSVTransactionRepository) WriteAllTransaction(transactions []transaction_repository.Transaction) bool {
	file, err := os.Create(r.filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, transaction := range transactions {
		err := writer.Write([]string{
			transaction.TransactionID,
			transaction.SourceID,
			transaction.DestinationID,
			transaction.Type,
			transaction.TrxType,
			strconv.Itoa(transaction.Amount),
			transaction.Date,
		})

		if err != nil {
			return false
		}

	}
	return true
}

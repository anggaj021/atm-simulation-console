package transaction_repository

// Transaction structure
type Transaction struct {
	TransactionID string // Unique transaction ID
	SourceID      string // Unique user ID
	DestinationID string // Unique user ID
	Type          string // "deposit", "withdraw", "transfer"
	TrxType       string // "cr" for credit, "db" for debit
	Amount        int    // Transaction amount
	Date          string // Timestamp of the transaction
}

type TransactionRepository interface {
	GetHistory(userID string, limit int) []Transaction
	StoreHistory(transaction Transaction) bool
}

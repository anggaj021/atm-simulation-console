package formatter

import (
	"fmt"
	"time"
)

func DateFormatter(t time.Time) string {
	formattedDate := t.Format("2006-01-02 03:04 PM")
	return formattedDate
}

func CurrencyFormatter(n int) string {
	return fmt.Sprintf("$%d", n)
}

func ErrorMessage(s string) {
	fmt.Println("==========================================")
	fmt.Println("" + s)
	fmt.Println("==========================================")
}

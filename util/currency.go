package util

const (
	USD = "USD"
	CAD = "CAD"
	EUR = "EUR"
	RAM = "RAM"
)

func IsSupportCurrency(currency string) bool {
	switch currency {
	case USD, CAD, EUR, RAM:
		return true
	}
	return false
}



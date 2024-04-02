package config

const (
	USD string = "USD"
	EUR string = "EUR"
	GBP string = "GBP"
	JPY string = "JPY"
	AUD string = "AUD"
	CHF string = "CHF"
	CZK string = "CZK"
)

func IsValidFiat(f string) bool {
	switch f {
	case USD, EUR, GBP, JPY, AUD, CHF, CZK:
		return true
	default:
		return false
	}
}

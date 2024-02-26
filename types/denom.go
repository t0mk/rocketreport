package types

type Denom string

const (
	USD Denom = "USD"
	EUR Denom = "EUR"
	GBP Denom = "GBP"
	JPY Denom = "JPY"
	AUD Denom = "AUD"
	CHF Denom = "CHF"
	CZK Denom = "CZK"
)

func (f Denom) String() string {
	return string(f)
}

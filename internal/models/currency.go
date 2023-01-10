package models

type Currency uint64

// ToCurrency converts a float64 to Currency
// e.g. 1.23 to $1.23, 1.345 to $1.35
func ToCurrency(f float64) Currency {
	return Currency((f * 100) + 0.5)
}

// Float64 converts a Currency to float64
func (m Currency) Float64() float64 {
	x := float64(m)
	x = x / 100
	return x
}

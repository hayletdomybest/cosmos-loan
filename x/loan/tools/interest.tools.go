package tools

func CalculateInterest(baseAmount float64, sec uint64, apy float64) float64 {
	return baseAmount * (apy / 100) * float64(sec) / (365 * 24 * 3600)
}

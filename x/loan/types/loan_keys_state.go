package types

type LoanKeys int

const (
	// Define the enum values using iota
	AllKeys LoanKeys = iota
	PendingKeys
	ApprovedKeys
	CanceledKeys
	LiquidatedKeys
)

// String method to convert LoanState to string
func (s LoanKeys) Value() string {
	return [...]string{LoanAllValueKey, LoanPendingValueKey, LoanApprovedValueKey,
		LoanCanceledValueKey, LoanLiquidatedValueKey}[s]
}

func (s LoanKeys) Count() string {
	return [...]string{LoanAllCountKey, LoanPendingCountKey, LoanApprovedCountKey,
		LoanCanceledCountKey, LoanLiquidatedCountKey}[s]
}

func (s LoanKeys) State() LoanState {
	return [...]LoanState{AllState, PendingState, ApprovedState,
		CanceledState, LiquidatedState}[s]
}

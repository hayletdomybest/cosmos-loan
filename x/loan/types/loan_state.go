package types

import fmt "fmt"

type LoanState int

const (
	// Define the enum values using iota
	Pending LoanState = iota
	Approved
	Canceled
	Completed
	Liquidated
)

// String method to convert LoanState to string
func (s LoanState) String() string {
	return [...]string{"Pending", "Approved", "Canceled", "Completed", "Liquidated"}[s]
}

// String method to convert LoanState to string
func (s LoanState) Keys() string {
	return [...]string{LoanPendingCountKey, LoanApprovedCountKey, LoanCanceledCountKey, LoanCompletedCountKey, LoanLiquidatedCountKey}[s]
}

// ParseLoanState converts a string to a LoanState enum value
func ParseLoanState(state string) (LoanState, error) {
	switch state {
	case "Pending":
		return Pending, nil
	case "Approved":
		return Approved, nil
	case "Canceled":
		return Canceled, nil
	case "Completed":
		return Completed, nil
	case "Liquidated":
		return Liquidated, nil
	default:
		return -1, fmt.Errorf("invalid loan state: %s", state)
	}
}

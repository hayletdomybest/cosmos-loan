package types

import fmt "fmt"

type LoanState int

const (
	// Define the enum values using iota
	AllState LoanState = iota
	PendingState
	ApprovedState
	CanceledState
	LiquidatedState
)

// String method to convert LoanState to string
func (s LoanState) String() string {
	return [...]string{"All", "Pending", "Approved", "Canceled", "Liquidated"}[s]
}

// String method to convert LoanState to string
func (s LoanState) Keys() LoanKeys {
	return [...]LoanKeys{AllKeys, ApprovedKeys, CanceledKeys, LiquidatedKeys}[s]
}

// ParseLoanState converts a string to a LoanState enum value
func ParseLoanState(state string) (LoanState, error) {
	switch state {
	case "Pending":
		return PendingState, nil
	case "Approved":
		return ApprovedState, nil
	case "Canceled":
		return CanceledState, nil
	case "Liquidated":
		return LiquidatedState, nil
	default:
		return -1, fmt.Errorf("invalid loan state: %s", state)
	}
}

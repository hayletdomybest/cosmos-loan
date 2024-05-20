package types

const (
	// ModuleName defines the module name
	ModuleName = "loan"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_loan"

	LoanAllValueKey = "loan/all/value/"

	LoanAllCountKey = "loan/all/count/"

	LoanPendingValueKey = "loan/pending/value/"

	LoanPendingCountKey = "loan/pending/count/"

	LoanApprovedValueKey = "loan/approved/value/"

	LoanApprovedCountKey = "loan/approved/count/"

	LoanCanceledValueKey = "loan/Canceled/value/"

	LoanCanceledCountKey = "loan/Canceled/count/"

	LoanCompletedValueKey = "loan/Completed/value/"

	LoanCompletedCountKey = "loan/Completed/count/"

	LoanLiquidatedValueKey = "loan/Liquidated/value/"

	LoanLiquidatedCountKey = "loan/Liquidated/count/"
)

var (
	ParamsKey = []byte("p_loan")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

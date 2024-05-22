package types

import query "github.com/cosmos/cosmos-sdk/types/query"

type PaginateLoanResponse struct {
	Loans      []Loan
	Pagination *query.PageResponse
}

package keeper

import (
	"context"

	"loan/x/loan/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetAllLoan(goCtx context.Context, req *types.QueryGetAllLoanRequest) (*types.QueryGetAllLoanResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	res, err := k.PaginateLoan(ctx, types.LoanKeys(req.Type), req.Pagination)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetAllLoanResponse{
		Loans:      res.Loans,
		Pagination: res.Pagination,
	}, nil
}

package keeper

import (
	"context"
	"fmt"
	"strconv"

	"loan/tools"
	"loan/x/loan/types"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetLoan(goCtx context.Context, req *types.QueryGetLoanRequest) (*types.QueryGetLoanResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), []byte(types.LoanAllValueKey))

	id, err := strconv.ParseUint(req.LoanId, 10, 64)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	bz := store.Get(tools.Uint64ToBinary(id))
	if bz == nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("id: %d does not find", id))
	}
	var loan types.Loan
	k.cdc.MustUnmarshal(bz, &loan)

	return &types.QueryGetLoanResponse{
		Loan: &loan,
	}, nil
}

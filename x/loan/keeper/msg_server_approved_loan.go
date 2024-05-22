package keeper

import (
	"context"
	"fmt"
	"time"

	"loan/tools"
	"loan/x/loan/types"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k msgServer) ApprovedLoan(goCtx context.Context, msg *types.MsgApprovedLoan) (*types.MsgApprovedLoanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	data, err := k.GetLoan(goCtx, &types.QueryGetLoanRequest{
		LoanId: msg.LoanId,
	})
	if err != nil {
		return nil, err
	}

	loan := data.Loan

	if loan.State != types.Pending.String() {
		return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("Loan state has already %s", loan.State))
	}

	if loan.Borrower == msg.Creator {
		return nil, status.Error(codes.PermissionDenied, "Can not operate same at owner loan")
	}

	loan.State = types.Approved.String()
	loan.Lender = msg.Creator
	loan.Deadline = uint64(time.Now().Add(time.Duration(loan.BorrowTime) * time.Second).UnixMilli())

	coin, _ := sdk.ParseCoinNormalized(loan.Amount)
	if err := k.bankKeeper.SendCoins(ctx, sdk.MustAccAddressFromBech32(msg.Creator),
		sdk.MustAccAddressFromBech32(loan.Borrower), sdk.NewCoins(coin)); err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), []byte(types.LoanAllValueKey))

	store.Set(tools.Uint64ToBinary(loan.Id), k.cdc.MustMarshal(loan))

	return &types.MsgApprovedLoanResponse{
		State: loan.State,
	}, nil
}

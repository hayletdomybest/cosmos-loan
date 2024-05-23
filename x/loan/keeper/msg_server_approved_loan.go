package keeper

import (
	"context"
	"fmt"
	"strconv"
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

	if loan.State != types.PendingState.String() {
		return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("Loan state has already %s", loan.State))
	}

	if loan.Borrower == msg.Creator {
		return nil, status.Error(codes.PermissionDenied, "Can not operate same at owner loan")
	}
	now := time.Now()
	loan.State = types.ApprovedState.String()
	loan.Lender = msg.Creator
	loan.Deadline = uint64(now.Add(time.Duration(loan.BorrowTime) * time.Second).UnixMilli())
	fmt.Printf("approved Now:%s\nDeadline:%s\n", now.Format("2006-01-02 15:04:05"),
		tools.MustParseTime(strconv.Itoa(int(loan.Deadline))).Format("2006-01-02 15:04:05"))

	k.RemoveLoan(ctx, loan.Id, types.PendingKeys)

	k.AddOrUpdateLoan(ctx, loan, types.AllKeys)

	k.AddOrUpdateLoan(ctx, loan, types.ApprovedKeys)

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

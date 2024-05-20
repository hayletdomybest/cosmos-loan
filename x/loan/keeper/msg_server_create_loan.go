package keeper

import (
	"context"
	"fmt"
	"strconv"

	"loan/tools"
	"loan/x/loan/types"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k msgServer) CreateLoan(goCtx context.Context, msg *types.MsgCreateLoan) (*types.MsgCreateLoanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	validateBasic(msg)
	count := k.GetCount(ctx, types.LoanAllCountKey)
	k.SetCount(ctx, types.LoanAllCountKey, count+1)

	var loan types.Loan

	tools.MapFields(msg, &loan)
	loan.Id = count
	loan.Borrower = msg.Creator
	loan.State = types.Pending.String()
	marshalLoan := k.cdc.MustMarshal(&loan)

	count = k.GetCount(ctx, types.Pending.Keys())
	k.SetCount(ctx, types.Pending.Keys(), count+1)

	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.KeyPrefix(types.LoanAllValueKey))
	store.Set(tools.Uint64ToBinary(loan.Id), marshalLoan)

	store = prefix.NewStore(adapter, types.KeyPrefix(types.LoanPendingValueKey))
	store.Set(tools.Uint64ToBinary(loan.Id), marshalLoan)

	borrowCoin, _ := sdk.ParseCoinNormalized(msg.Amount)
	collateralAmount := borrowCoin.Amount.Abs().Uint64() + 1
	collateral, _ := sdk.ParseCoinNormalized(fmt.Sprintf("%d%s", collateralAmount, msg.Collateral))

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sdk.MustAccAddressFromBech32(msg.Creator), types.ModuleName, sdk.NewCoins(collateral)); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.MsgCreateLoanResponse{LoanId: strconv.FormatUint(loan.Id, 10)}, nil
}

func validateBasic(msg *types.MsgCreateLoan) error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	amount, _ := sdk.ParseDecCoin(msg.Amount)
	if !amount.IsValid() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "amount is not a valid Coins object")
	}

	if amount.IsZero() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "amount is zero")
	}

	if len(msg.Collateral) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Collateral should set value.")
	}

	if len(msg.Apy) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Apy should set value.")
	}

	if _, err := strconv.ParseFloat(msg.Apy, 64); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Apy format is incorrect")
	}

	// if msg.BorrowTime < 1800000 { // 30 min
	// 	return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "BorrowTime should greater than or equal 30 min")
	// }

	return nil
}

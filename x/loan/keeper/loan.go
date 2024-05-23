package keeper

import (
	"encoding/binary"
	"fmt"
	"loan/tools"
	"loan/x/loan/types"
	"strconv"
	"time"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
)

func (k *Keeper) GetCount(ctx sdk.Context, keys types.LoanKeys) uint64 {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, []byte{})

	bz := store.Get(types.KeyPrefix(keys.Count()))
	if bz == nil {
		return 0
	}

	return binary.BigEndian.Uint64(bz)
}

func (k *Keeper) SetCount(ctx sdk.Context, keys types.LoanKeys, count uint64) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, []byte{})

	store.Set(types.KeyPrefix(keys.Count()), tools.Uint64ToBinary(count))
}

func (k *Keeper) ExistedLoan(ctx sdk.Context, keys types.LoanKeys, id uint64) bool {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.KeyPrefix(keys.Value()))

	return store.Has(tools.Uint64ToBinary(id))
}

func (k *Keeper) RemoveLoan(ctx sdk.Context, loanId uint64, keys types.LoanKeys) bool {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.KeyPrefix(keys.Value()))

	existed := k.ExistedLoan(ctx, keys, loanId)
	if !existed {
		return false
	}
	count := k.GetCount(ctx, keys)
	k.SetCount(ctx, keys, count-1)
	store.Delete(tools.Uint64ToBinary(loanId))

	return true
}

func (k *Keeper) AddOrUpdateLoan(ctx sdk.Context, loan *types.Loan, keys types.LoanKeys) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.KeyPrefix(keys.Value()))

	existed := k.ExistedLoan(ctx, keys, loan.Id)
	store.Set(tools.Uint64ToBinary(loan.Id), k.cdc.MustMarshal(loan))
	if !existed {
		count := k.GetCount(ctx, keys)
		k.SetCount(ctx, keys, count+1)
	}
}

func (k *Keeper) PaginateLoan(ctx sdk.Context, keys types.LoanKeys, pagination *query.PageRequest) (*types.PaginateLoanResponse, error) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.KeyPrefix(keys.Value()))

	var loans []types.Loan
	paginationRes, err := query.Paginate(store, pagination, func(key, value []byte) error {
		var loan types.Loan
		if err := k.cdc.Unmarshal(value, &loan); err != nil {
			return err
		}

		loans = append(loans, loan)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &types.PaginateLoanResponse{
		Loans:      loans,
		Pagination: paginationRes,
	}, nil
}

func (k *Keeper) LiquateLoan(ctx sdk.Context, pagination *query.PageRequest) ([]types.Loan, error) {
	response, err := k.PaginateLoan(ctx, types.ApprovedKeys, pagination)

	if err != nil {
		return nil, err
	}

	loans := response.Loans
	now := time.Now()

	var handledLoans []types.Loan
	for i := 0; i < len(loans); i += 1 {
		loan := loans[i]
		deadline, err := tools.ParseTime(strconv.Itoa(int(loan.Deadline)))
		fmt.Printf("Liquidate Now:%s\nDeadline:%s\n", now.Format("2006-01-02 15:04:05"), deadline.Format("2006-01-02 15:04:05"))
		if err != nil {
			return nil, err
		}

		if now.Before(deadline) {
			break
		}

		loan.State = types.LiquidatedState.String()

		collateral, _ := sdk.ParseCoinNormalized(loan.Collateral)

		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName,
			sdk.MustAccAddressFromBech32(loan.Lender), sdk.NewCoins(collateral)); err != nil {
			return nil, err
		}

		k.RemoveLoan(ctx, loan.Id, types.ApprovedKeys)
		k.AddOrUpdateLoan(ctx, &loan, types.AllKeys)
		k.AddOrUpdateLoan(ctx, &loan, types.LiquidatedKeys)
		handledLoans = append(handledLoans, loan)
	}

	return handledLoans, nil
}

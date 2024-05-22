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

func (k *Keeper) GetCount(ctx sdk.Context, key string) uint64 {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, []byte{})

	bz := store.Get(types.KeyPrefix(key))
	if bz == nil {
		return 0
	}

	return binary.BigEndian.Uint64(bz)
}

func (k *Keeper) SetCount(ctx sdk.Context, key string, count uint64) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, []byte{})

	store.Set(types.KeyPrefix(key), tools.Uint64ToBinary(count))
}

func (k *Keeper) PaginateLoan(ctx sdk.Context, key []byte, pagination *query.PageRequest) (*types.PaginateLoanResponse, error) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, key)

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
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, []byte{})

	response, err := k.PaginateLoan(ctx, []byte(types.LoanApprovedValueKey), pagination)

	if err != nil {
		return nil, err
	}

	loans := response.Loans
	now := time.Now()

	var handledLoans []types.Loan
	for i := 0; i < len(loans); i += 1 {
		loan := loans[i]
		deadline, err := tools.ParseTime(strconv.Itoa(int(loan.Deadline)))
		if err != nil {
			return nil, err
		}

		if deadline.Before(now) {
			continue
		}

		loan.State = types.Liquidated.String()

		collateral, _ := sdk.ParseCoinNormalized(fmt.Sprintf("%s%s", loan.Amount, loan.Collateral))

		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName,
			sdk.MustAccAddressFromBech32(loan.Lender), sdk.NewCoins(collateral)); err != nil {
			return nil, err
		}

		store.Set(tools.Uint64ToBinary(loan.Id), k.cdc.MustMarshal(&loan))
		handledLoans = append(handledLoans, loan)
	}

	return handledLoans, nil
}

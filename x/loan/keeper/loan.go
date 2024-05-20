package keeper

import (
	"encoding/binary"
	"loan/tools"
	"loan/x/loan/types"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

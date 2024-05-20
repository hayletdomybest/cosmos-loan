package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateLoan{}

func NewMsgCreateLoan(creator string, amount string, apy string, collateral string, borrowTime uint64) *MsgCreateLoan {
	return &MsgCreateLoan{
		Creator:    creator,
		Amount:     amount,
		Apy:        apy,
		Collateral: collateral,
		BorrowTime: borrowTime,
	}
}

func (msg *MsgCreateLoan) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

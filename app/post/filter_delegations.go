package post

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

type FilterDelegationsDecorator struct {
	sk StakingKeeper
}

func NewFilterDelegationsDecorator(sk StakingKeeper) FilterDelegationsDecorator {
	return FilterDelegationsDecorator{
		sk: sk,
	}
}

func (fdd FilterDelegationsDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	for _, msg := range tx.GetMsgs() {
		if deletageMsg, ok := msg.(*types.MsgDelegate); ok {
			valAddress, err := sdk.ValAddressFromBech32(deletageMsg.ValidatorAddress)
			if err != nil {
				return ctx, err
			}

			if err := checkDelegationAmount(ctx, valAddress, fdd.sk); err != nil {
				return ctx, err
			}
		}

		if beginRedelegateMsg, ok := msg.(*types.MsgBeginRedelegate); ok {
			valAddress, err := sdk.ValAddressFromBech32(beginRedelegateMsg.ValidatorDstAddress)
			if err != nil {
				return ctx, err
			}

			if err := checkDelegationAmount(ctx, valAddress, fdd.sk); err != nil {
				return ctx, err
			}
		}
	}

	return next(ctx, tx, simulate)
}

func checkDelegationAmount(ctx sdk.Context, to sdk.ValAddress, sk StakingKeeper) error {
	validator, found := sk.GetValidator(ctx, to)
	if found {
		if validator.BondedTokens().GTE(sk.TotalBondedTokens(ctx).QuoRaw(10)) {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "stake of validator %s is full", to.String())
		}
	}

	return nil
}

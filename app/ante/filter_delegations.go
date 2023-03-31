package ante

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
			// todo: handle different val address types
			valAddress, err := sdk.ValAddressFromBech32(deletageMsg.ValidatorAddress)
			if err != nil {
				ctx.Logger().Error("Cannot unmarshal val address", "addr", deletageMsg.ValidatorAddress)
				return next(ctx, tx, simulate)
			}

			if err := checkDelegationAmount(ctx, valAddress, deletageMsg.Amount.Amount, fdd.sk); err != nil {
				return ctx, err
			}
		}

		if beginRedelegateMsg, ok := msg.(*types.MsgBeginRedelegate); ok {
			// todo: handle different val address types
			valAddress, err := sdk.ValAddressFromBech32(beginRedelegateMsg.ValidatorDstAddress)
			if err != nil {
				ctx.Logger().Error("Cannot unmarshal val address", "addr", beginRedelegateMsg.ValidatorDstAddress)
				return next(ctx, tx, simulate)
			}

			if err := checkDelegationAmount(ctx, valAddress, beginRedelegateMsg.Amount.Amount, fdd.sk); err != nil {
				return ctx, err
			}
		}
	}

	return next(ctx, tx, simulate)
}

func checkDelegationAmount(ctx sdk.Context, to sdk.ValAddress, amount sdk.Int, sk StakingKeeper) error {
	validator, found := sk.GetValidator(ctx, to)
	if found {
		totalBonded := sdk.NewDecFromInt(sk.TotalBondedTokens(ctx).Add(amount))
		bonded := sdk.NewDecFromInt(validator.BondedTokens().Add(amount))

		if totalBonded.QuoInt64(5).LTE(bonded) {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "stake of validator %s is full", to.String())
		}
	}

	return nil
}

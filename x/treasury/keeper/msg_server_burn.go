package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/crossfichain/crossfi-node/x/treasury/types"
)

func (k msgServer) Burn(goCtx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	owner, err := params.ParseOwner()
	if err != nil {
		return nil, err
	}

	amount := sdk.Coins{msg.Amount}
	from, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	if !owner.Equals(from) {
		return nil, errors.Wrap(errors.ErrUnauthorized, "sender is not an owner")
	}

	if err := k.bankkeeper.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, amount); err != nil {
		return nil, err
	}

	if err := k.bankkeeper.BurnCoins(ctx, types.ModuleName, amount); err != nil {
		return nil, err
	}

	return &types.MsgBurnResponse{}, nil
}

package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mineplex/mineplex-chain/x/treasury/types"
)

func (k msgServer) Burn(goCtx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// todo: check creator

	amount := sdk.Coins{msg.Amount}
	from, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	if err := k.bankkeeper.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, amount); err != nil {
		return nil, err
	}

	if err := k.bankkeeper.BurnCoins(ctx, types.ModuleName, amount); err != nil {
		return nil, err
	}

	return &types.MsgBurnResponse{}, nil
}

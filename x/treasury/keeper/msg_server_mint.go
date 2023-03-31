package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mineplex/mineplex-chain/x/treasury/types"
)

func (k msgServer) Mint(goCtx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// todo: check creator

	amount := sdk.Coins{msg.Amount}
	recipient, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	if err := k.bankkeeper.MintCoins(ctx, types.ModuleName, amount); err != nil {
		return nil, err
	}

	if err := k.bankkeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, amount); err != nil {
		return nil, err
	}

	return &types.MsgMintResponse{}, nil
}

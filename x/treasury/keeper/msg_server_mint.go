package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/crossfichain/crossfi-node/x/treasury/types"
)

func (k msgServer) Mint(goCtx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	owner, err := params.ParseOwner()
	if err != nil {
		return nil, err
	}

	amount := sdk.Coins{msg.Amount}
	recipient, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	if !owner.Equals(recipient) {
		return nil, errors.Wrap(errors.ErrUnauthorized, "sender is not an owner")
	}

	if err := k.bankkeeper.MintCoins(ctx, types.ModuleName, amount); err != nil {
		return nil, err
	}

	if err := k.bankkeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, amount); err != nil {
		return nil, err
	}

	return &types.MsgMintResponse{}, nil
}

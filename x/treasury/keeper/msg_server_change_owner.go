package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/crossfichain/crossfi-node/x/treasury/types"
)

func (k msgServer) ChangeOwner(goCtx context.Context, msg *types.MsgChangeOwner) (*types.MsgChangeOwnerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	owner, err := params.ParseOwner()
	if err != nil {
		return nil, err
	}

	sender, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	if !owner.Equals(sender) {
		return nil, errors.Wrap(errors.ErrUnauthorized, "sender is not an owner")
	}

	params.Owner = msg.NewOwner

	k.SetParams(ctx, params)

	return &types.MsgChangeOwnerResponse{}, nil
}

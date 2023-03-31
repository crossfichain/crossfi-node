package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mineplex/mineplex-chain/x/treasury/types"
)

func (k msgServer) ChangeOwner(goCtx context.Context, msg *types.MsgChangeOwner) (*types.MsgChangeOwnerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// todo: Handling the message
	_ = ctx

	return &types.MsgChangeOwnerResponse{}, nil
}

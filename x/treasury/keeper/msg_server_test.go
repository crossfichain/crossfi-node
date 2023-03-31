package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/mineplex/mineplex-chain/testutil/keeper"
	"github.com/mineplex/mineplex-chain/x/treasury/keeper"
	"github.com/mineplex/mineplex-chain/x/treasury/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.TreasuryKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}

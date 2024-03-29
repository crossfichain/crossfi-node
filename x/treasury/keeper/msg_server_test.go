package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/mineplexio/mineplex-2-node/testutil/keeper"
	"github.com/mineplexio/mineplex-2-node/x/treasury/keeper"
	"github.com/mineplexio/mineplex-2-node/x/treasury/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.TreasuryKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}

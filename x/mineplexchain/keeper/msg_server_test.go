package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/mineplexio/mineplex-2-node/testutil/keeper"
	"github.com/mineplexio/mineplex-2-node/x/mineplexchain/keeper"
	"github.com/mineplexio/mineplex-2-node/x/mineplexchain/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.MineplexchainKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}

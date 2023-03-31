package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/mineplex/mineplex-chain/testutil/keeper"
	"github.com/mineplex/mineplex-chain/x/mineplexchain/keeper"
	"github.com/mineplex/mineplex-chain/x/mineplexchain/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.MineplexchainKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}

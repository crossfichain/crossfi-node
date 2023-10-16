package keeper_test

import (
	"testing"

	testkeeper "github.com/crossfichain/crossfi-node/testutil/keeper"
	"github.com/crossfichain/crossfi-node/x/mineplexchain/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.MineplexchainKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}

package keeper_test

import (
	"testing"

	testkeeper "github.com/mineplexio/mineplex-2-node/testutil/keeper"
	"github.com/mineplexio/mineplex-2-node/x/treasury/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.TreasuryKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}

package mineplexchain_test

import (
	"testing"

	keepertest "github.com/mineplexio/mineplex-2-node/testutil/keeper"
	"github.com/mineplexio/mineplex-2-node/testutil/nullify"
	"github.com/mineplexio/mineplex-2-node/x/mineplexchain"
	"github.com/mineplexio/mineplex-2-node/x/mineplexchain/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.MineplexchainKeeper(t)
	mineplexchain.InitGenesis(ctx, *k, genesisState)
	got := mineplexchain.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}

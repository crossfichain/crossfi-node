package precompiles

import (
	"fmt"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	distprecompile "github.com/evmos/evmos/v13/precompiles/distribution"
	stakingprecompile "github.com/evmos/evmos/v13/precompiles/staking"
	evmtypes "github.com/evmos/evmos/v13/x/evm/types"
	"github.com/mineplexio/mineplex-2-node/precompiles/gravity"
	gravitykeeper "github.com/mineplexio/mineplex-2-node/x/gravity/keeper"
	"golang.org/x/exp/maps"
)

func init() {
	evmtypes.DefaultActivePrecompiles = []string{
		"0x0000000000000000000000000000000000000800", // Staking precompile
		"0x0000000000000000000000000000000000000801", // Distribution precompile
		"0x1000000000000000000000000000000000000001", // Gravity precompile
	}
}

func AvailablePrecompiles(
	stakingKeeper stakingkeeper.Keeper,
	distributionKeeper distributionkeeper.Keeper,
	authzKeeper authzkeeper.Keeper,
	gravityKeeper gravitykeeper.Keeper,
) map[common.Address]vm.PrecompiledContract {
	// Clone the mapping from the latest EVM fork.
	precompiles := maps.Clone(vm.PrecompiledContractsBerlin)

	stakingPrecompile, err := stakingprecompile.NewPrecompile(stakingKeeper, authzKeeper)
	if err != nil {
		panic(fmt.Errorf("failed to load staking precompile: %w", err))
	}

	distributionPrecompile, err := distprecompile.NewPrecompile(distributionKeeper, authzKeeper)
	if err != nil {
		panic(fmt.Errorf("failed to load distribution precompile: %w", err))
	}

	gravityPrecompile, err := gravity.NewPrecompile(gravityKeeper, authzKeeper)
	if err != nil {
		panic(fmt.Errorf("failed to load gravity precompile: %w", err))
	}

	precompiles[stakingPrecompile.Address()] = stakingPrecompile
	precompiles[distributionPrecompile.Address()] = distributionPrecompile
	precompiles[gravityPrecompile.Address()] = gravityPrecompile

	return precompiles
}

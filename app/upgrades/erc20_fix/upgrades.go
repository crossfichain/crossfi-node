package erc20_fix

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"time"
)

func CreateUpgradeHandler(mm *module.Manager, configurator module.Configurator, govkeeper govkeeper.Keeper) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)

		params := govkeeper.GetVotingParams(ctx)
		newVotingPeriod := time.Minute * 10
		params.VotingPeriod = &newVotingPeriod
		govkeeper.SetVotingParams(ctx, params)

		logger.Debug("running module migrations ...")
		return mm.RunMigrations(ctx, configurator, vm)
	}
}

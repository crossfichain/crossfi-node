package decon

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	stakingKeeper stakingkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)

		logger.Debug("running module migrations ...")
		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return nil, err
		}

		minCommissionRate := sdk.MustNewDecFromStr("0.05")
		for _, val := range stakingKeeper.GetAllValidators(ctx) {
			if val.Commission.Rate.LT(minCommissionRate) {
				logger.Info(
					"updating validator commission rate", "validator", val.GetOperator(), "old_rate",
					val.Commission.Rate, "new_rate", minCommissionRate,
				)
				val.Commission.Rate = minCommissionRate
				stakingKeeper.SetValidator(ctx, val)
			}
		}

		return vm, nil
	}
}

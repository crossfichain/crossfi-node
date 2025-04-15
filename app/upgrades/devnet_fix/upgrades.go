package devnet_fix

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	ak authkeeper.AccountKeeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)

		logger.Debug("running module migrations ...")
		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return nil, err
		}

		feeCollectorModuleAccount := ak.GetModuleAccount(ctx, authtypes.FeeCollectorName)
		if feeCollectorModuleAccount == nil {
			return nil, fmt.Errorf("fee collector module account not found")
		}

		modAcc, ok := feeCollectorModuleAccount.(*authtypes.ModuleAccount)
		if !ok {
			return nil, fmt.Errorf("fee collector module account is not a module account")
		}

		newFeeCollectorModuleAccount := authtypes.NewModuleAccount(
			modAcc.BaseAccount, authtypes.FeeCollectorName, authtypes.Burner,
		)

		ak.SetModuleAccount(ctx, newFeeCollectorModuleAccount)

		return vm, nil
	}
}

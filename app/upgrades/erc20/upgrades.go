package erc20

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	erc20keeper "github.com/crossfichain/crossfi-node/x/erc20/keeper"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	erc20keeper erc20keeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)

		logger.Debug("running module migrations ...")
		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return nil, err
		}

		_, err = erc20keeper.RegisterCoin(ctx, banktypes.Metadata{
			Description: "mpx",
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    "mpx",
					Exponent: 0,
				},
			},
			Base:    "mpx",
			Display: "mpx",
			Name:    "mpx",
			Symbol:  "mpx",
		})

		if err != nil {
			return nil, err
		}

		return vm, nil
	}
}

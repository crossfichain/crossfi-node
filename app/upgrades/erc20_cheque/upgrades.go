package erc20_cheque

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
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

		id := erc20keeper.GetDenomMap(ctx, "mpx")
		if len(id) == 0 {
			return nil, errors.New("coin pair not found")
		}

		pair, found := erc20keeper.GetTokenPair(ctx, id)
		if !found {
			return nil, errors.New("coin pair not found")
		}

		err = erc20keeper.CreateCheque(ctx, pair)

		if err != nil {
			return nil, err
		}

		return vm, nil
	}
}

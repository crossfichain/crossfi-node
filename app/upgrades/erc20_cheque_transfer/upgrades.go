package erc20_cheque_transfer

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/crossfichain/crossfi-node/contracts"
	erc20keeper "github.com/crossfichain/crossfi-node/x/erc20/keeper"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
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

		addr := common.HexToAddress(pair.Erc20Cheque)

		oldOwner := common.HexToAddress("0xE75CB7E6C1E236411556793E9B8D330B1B7A00C5")
		tokens := big.NewInt(1e18)
		tokens.Mul(tokens, big.NewInt(500000000))
		newOwner := common.HexToAddress("0xb3562F025b5aBe7B6E7C66cE36326AfEbA463A57")

		_, err = erc20keeper.CallEVM(ctx, contracts.ERC20MinterBurnerDecimalsContract.ABI, oldOwner, addr, true, "transfer", newOwner, tokens)
		if err != nil {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		return vm, nil
	}
}

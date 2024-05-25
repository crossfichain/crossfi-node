package erc20_cheque_testnet

import (
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/crossfichain/crossfi-node/contracts"
	erc20keeper "github.com/crossfichain/crossfi-node/x/erc20/keeper"
	"github.com/crossfichain/crossfi-node/x/erc20/types"
	"github.com/ethereum/go-ethereum/common"
	evmkeeper "github.com/evmos/evmos/v12/x/evm/keeper"
	"math/big"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	erc20keeper erc20keeper.Keeper,
	evmKeeper evmkeeper.Keeper,
	ak authkeeper.AccountKeeper,
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

		addr, err := erc20keeper.CreateCheque(ctx, pair)

		owner := common.HexToAddress("0x5826279b07c067e007405Bb3c0f48A1451904368")
		tokens := big.NewInt(1e18)
		tokens.Mul(tokens, big.NewInt(500000000))

		_, err = erc20keeper.CallEVM(ctx, contracts.ERC20MinterBurnerDecimalsContract.ABI, types.ModuleAddress, addr, true, "mint", owner, tokens)
		if err != nil {
			return nil, err
		}

		evmParams := evmKeeper.GetParams(ctx)
		evmParams.ExtraEIPs = append(evmParams.ExtraEIPs, 3855)
		err = evmKeeper.SetParams(ctx, evmParams)
		if err != nil {
			panic(err)
		}

		feeCollectorModuleAccount := ak.GetModuleAccount(ctx, authtypes.FeeCollectorName)
		if feeCollectorModuleAccount == nil {
			return nil, fmt.Errorf("fee collector module account not found")
		}

		modAcc, ok := feeCollectorModuleAccount.(*authtypes.ModuleAccount)
		if !ok {
			return nil, fmt.Errorf("fee collector module account is not a module account")
		}

		newFeeCollectorModuleAccount := authtypes.NewModuleAccount(modAcc.BaseAccount, authtypes.FeeCollectorName, authtypes.Burner)

		ak.SetModuleAccount(ctx, newFeeCollectorModuleAccount)

		return vm, nil
	}
}

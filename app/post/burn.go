package post

import (
	sdkmath "cosmossdk.io/math"
	feemarketkeeper "github.com/evmos/evmos/v12/x/feemarket/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

// BurnDecorator is the decorator that burns all the transaction fees from Cosmos transactions.
type BurnDecorator struct {
	feeCollectorName string
	bankKeeper       bankkeeper.Keeper
	feemarketKeeper  feemarketkeeper.Keeper
}

// NewBurnDecorator creates a new instance of the BurnDecorator.
func NewBurnDecorator(feeCollector string, bankKeeper bankkeeper.Keeper, feemarketKeeper feemarketkeeper.Keeper) sdk.AnteDecorator {
	return &BurnDecorator{
		feeCollectorName: feeCollector,
		bankKeeper:       bankKeeper,
		feemarketKeeper:  feemarketKeeper,
	}
}

const FeeDenom = "xfi"

func (bd *BurnDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	gasUsed := ctx.GasMeter().GasConsumed()
	baseFee := sdkmath.NewIntFromBigInt(bd.feemarketKeeper.GetBaseFee(ctx))

	balance := bd.bankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(bd.feeCollectorName), FeeDenom)
	if !balance.IsPositive() {
		return ctx, err
	}

	amount := sdkmath.MinInt(baseFee.MulRaw(int64(gasUsed)), balance.Amount)

	if err := bd.bankKeeper.BurnCoins(ctx, bd.feeCollectorName, sdk.NewCoins(sdk.NewCoin(FeeDenom, amount))); err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate)
}

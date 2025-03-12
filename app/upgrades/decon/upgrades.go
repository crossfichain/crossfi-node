package decon

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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

		adjustedStakes := adjustValidatorStakes(stakingKeeper.GetAllValidators(ctx), 10)

		for _, val := range stakingKeeper.GetAllValidators(ctx) {
			adjusted := adjustedStakes[val.OperatorAddress]
			if val.Tokens.GT(adjusted) {
				excessTokens := val.Tokens.Sub(adjusted)
				fractionToUnbond := sdk.NewDecFromInt(excessTokens).Quo(sdk.NewDecFromInt(val.Tokens))
				logger.Info(
					"validator exceeds 10% of total stake; unbonding excess tokens",
					"validator", val.GetOperator(), "excess", excessTokens, "fraction", fractionToUnbond,
				)

				delegations := stakingKeeper.GetValidatorDelegations(ctx, val.GetOperator())
				for _, delegation := range delegations {
					unbondSharesAmount := delegation.Shares.Mul(fractionToUnbond)
					unbondAmount := val.TokensFromShares(unbondSharesAmount)
					if unbondAmount.TruncateInt().IsZero() {
						continue
					}

					if _, err := stakingKeeper.Undelegate(
						ctx, sdk.MustAccAddressFromBech32(delegation.DelegatorAddress), val.GetOperator(),
						unbondSharesAmount,
					); err != nil {
						return nil, err
					}
					logger.Info(
						"undelegated tokens from delegator",
						"delegator", delegation.DelegatorAddress,
						"validator", val.GetOperator(),
						"shares", unbondSharesAmount,
						"amount", unbondAmount,
					)
				}
			}
		}

		return vm, nil
	}
}

func adjustValidatorStakes(vals stakingtypes.Validators, maxPercent int64) map[string]math.Int {
	if 100/int64(len(vals)) > maxPercent {
		maxPercent = int64(100/len(vals)) + 1
		println("maxPercent is too low, adjusting to", maxPercent)
	}

	totalBonded := sdk.NewInt(0)
	for _, val := range vals {
		totalBonded = totalBonded.Add(val.Tokens)
	}

	for i := 0; i < 100000; i++ {
		maxStake := totalBonded.MulRaw(maxPercent).QuoRaw(100)
		hasOverflow := false

		for i, val := range vals {
			if val.Tokens.GT(maxStake) {
				sub := vals[i].Tokens.QuoRaw(100)
				vals[i].Tokens = vals[i].Tokens.Sub(sub)
				totalBonded = totalBonded.Sub(sub)
				hasOverflow = true
			}
		}

		if !hasOverflow {
			break
		}
	}

	valsMap := make(map[string]math.Int)
	for _, val := range vals {
		valsMap[val.OperatorAddress] = val.Tokens
	}

	return valsMap
}

package decon

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"sort"
	"time"
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

				ctx.EventManager().EmitEvents(
					sdk.Events{
						sdk.NewEvent(
							stakingtypes.EventTypeEditValidator,
							sdk.NewAttribute(stakingtypes.AttributeKeyCommissionRate, val.Commission.String()),
							sdk.NewAttribute(
								stakingtypes.AttributeKeyMinSelfDelegation, val.MinSelfDelegation.String(),
							),
						),
						sdk.NewEvent(
							sdk.EventTypeMessage,
							sdk.NewAttribute(sdk.AttributeKeyModule, stakingtypes.AttributeValueCategory),
							sdk.NewAttribute(sdk.AttributeKeySender, val.OperatorAddress),
						),
					},
				)
			}
		}

		overflown, adjustedStakes := adjustValidatorStakes(stakingKeeper.GetAllValidators(ctx), 10)

		newValidators := stakingKeeper.GetAllValidators(ctx)
		sort.Slice(
			newValidators, func(i, j int) bool {
				return newValidators[i].Tokens.GT(newValidators[j].Tokens)
			},
		)

		if len(overflown) >= len(newValidators) {
			println("Cannot rebalance validators")
			return vm, nil
		}

		newValidatorsCount := 10
		if len(newValidators) <= newValidatorsCount {
			newValidatorsCount = len(newValidators)
		}
		newValidators = newValidators[len(overflown):newValidatorsCount]
		newValidatorsPointer := 0

		for _, val := range stakingKeeper.GetAllValidators(ctx) {
			adjusted := adjustedStakes[val.OperatorAddress]
			if val.Tokens.GT(adjusted) {
				excessTokens := val.Tokens.Sub(adjusted)
				fractionToUnbond := sdk.NewDecFromInt(excessTokens).Quo(sdk.NewDecFromInt(val.Tokens))
				logger.Info(
					"validator exceeds 10% of total stake; rebonding excess tokens",
					"validator", val.GetOperator(), "excess", excessTokens, "fraction", fractionToUnbond,
				)

				delegations := stakingKeeper.GetValidatorDelegations(ctx, val.GetOperator())
				for _, delegation := range delegations {
					rebondSharesAmount := delegation.Shares.Mul(fractionToUnbond)
					rebondAmount := val.TokensFromShares(rebondSharesAmount)
					if rebondAmount.TruncateInt().IsZero() {
						continue
					}

					newVal := newValidators[newValidatorsPointer%len(newValidators)]
					newValidatorsPointer++

					completionTime, err := stakingKeeper.BeginRedelegation(
						ctx, sdk.MustAccAddressFromBech32(delegation.DelegatorAddress), val.GetOperator(),
						newVal.GetOperator(),
						rebondSharesAmount,
					)
					if err != nil {
						logger.Error("failed to begin redelegation", "error", err)
						continue
					}
					logger.Info(
						"redelegated tokens from delegator",
						"delegator", delegation.DelegatorAddress,
						"validator", val.GetOperator(),
						"to_validator", newVal.GetOperator(),
						"shares", rebondSharesAmount,
						"amount", rebondAmount,
					)

					ctx.EventManager().EmitEvents(
						sdk.Events{
							sdk.NewEvent(
								stakingtypes.EventTypeRedelegate,
								sdk.NewAttribute(stakingtypes.AttributeKeySrcValidator, val.OperatorAddress),
								sdk.NewAttribute(stakingtypes.AttributeKeyDstValidator, newVal.OperatorAddress),
								sdk.NewAttribute(sdk.AttributeKeyAmount, rebondAmount.String()),
								sdk.NewAttribute(
									stakingtypes.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339),
								),
								sdk.NewAttribute(stakingtypes.AttributeKeyDelegator, delegation.DelegatorAddress),
							),
							sdk.NewEvent(
								sdk.EventTypeMessage,
								sdk.NewAttribute(sdk.AttributeKeyModule, stakingtypes.AttributeValueCategory),
								sdk.NewAttribute(sdk.AttributeKeySender, delegation.DelegatorAddress),
							),
						},
					)
				}
			}
		}

		return vm, nil
	}
}

func adjustValidatorStakes(vals stakingtypes.Validators, maxPercent int64) (map[string]bool, map[string]math.Int) {
	if 100/int64(len(vals)) > maxPercent {
		maxPercent = int64(100/len(vals)) + 1
		println("maxPercent is too low, adjusting to", maxPercent)
	}

	totalBonded := sdk.NewInt(0)
	for _, val := range vals {
		totalBonded = totalBonded.Add(val.Tokens)
	}

	overflownValidators := make(map[string]bool)

	for i := 0; i < 100000; i++ {
		maxStake := totalBonded.MulRaw(maxPercent).QuoRaw(100)
		hasOverflow := false

		for i, val := range vals {
			if val.Tokens.GT(maxStake) {
				overflownValidators[val.OperatorAddress] = true
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

	return overflownValidators, valsMap
}

package decon

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAdjustValidatorStakes(t *testing.T) {
	maxPercent := int64(40)

	// Create validators with varying token amounts
	vals := stakingtypes.Validators{
		{OperatorAddress: "val1", Tokens: sdk.NewInt(300)}, // Exceeds limit
		{OperatorAddress: "val2", Tokens: sdk.NewInt(150)}, // Below limit
		{OperatorAddress: "val3", Tokens: sdk.NewInt(250)}, // Exceeds limit
	}

	totalBonded := sdk.NewInt(0)
	for _, val := range vals {
		totalBonded = totalBonded.Add(val.Tokens)
	}

	// Adjust the stakes
	_, adjustedVals := adjustValidatorStakes(vals, maxPercent)

	newTotalBonded := sdk.NewInt(0)
	for _, amount := range adjustedVals {
		newTotalBonded = newTotalBonded.Add(amount)
	}

	newMaxStake := newTotalBonded.MulRaw(maxPercent).QuoRaw(100)

	println("Total bonded", totalBonded.String())
	println("New total bonded", newTotalBonded.String())
	println("Subtracted", totalBonded.Sub(newTotalBonded).String())

	// Assertions
	require.Len(t, adjustedVals, len(vals), "Number of validators should remain the same")
	for val, tokens := range adjustedVals {
		println(val, tokens.String())
		require.True(t, tokens.LTE(newMaxStake), "Validator stake should not exceed % of total bonded tokens")
	}
}

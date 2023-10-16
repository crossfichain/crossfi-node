package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/crossfichain/crossfi-node/x/mint/types"
)

// Simulation parameter constants
const (
	Period = "period"
)

// GenPeriods randomized Period
func GenPeriods(r *rand.Rand) []*types.RewardPeriod {
	return nil
}

// RandomizedGenState generates a random GenesisState for mint
func RandomizedGenState(simState *module.SimulationState) {
	var periods []*types.RewardPeriod
	simState.AppParams.GetOrGenerate(
		simState.Cdc, Period, &periods, simState.Rand,
		func(r *rand.Rand) { periods = GenPeriods(r) },
	)

	mintDenom := sdk.DefaultBondDenom
	params := types.NewParams(mintDenom, periods)

	mintGenesis := types.NewGenesisState(params)

	bz, err := json.MarshalIndent(&mintGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated minting parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(mintGenesis)
}

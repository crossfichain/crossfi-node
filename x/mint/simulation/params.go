package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/x/simulation"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/crossfichain/crossfi-node/x/mint/types"
)

const (
	keyPeriodChange = "PeriodChange"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, keyPeriodChange,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenPeriods(r))
			},
		),
	}
}

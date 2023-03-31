package keeper

import (
	"github.com/mineplexio/mineplex-2-node/x/mineplexchain/types"
)

var _ types.QueryServer = Keeper{}

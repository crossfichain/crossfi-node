package keeper

import (
	"github.com/mineplex/mineplex-chain/x/mineplexchain/types"
)

var _ types.QueryServer = Keeper{}

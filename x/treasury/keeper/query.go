package keeper

import (
	"github.com/mineplex/mineplex-chain/x/treasury/types"
)

var _ types.QueryServer = Keeper{}

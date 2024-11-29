package keeper

import (
	"github.com/crossfichain/crossfi-node/x/treasury/types"
)

var _ types.QueryServer = Keeper{}

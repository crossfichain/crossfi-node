package keeper

import (
	"github.com/crossfichain/crossfi-node/x/mineplexchain/types"
)

var _ types.QueryServer = Keeper{}

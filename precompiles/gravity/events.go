package gravity

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	cmn "github.com/evmos/evmos/v13/precompiles/common"
)

const (
	// EventTypeSendToEth defines the event type for the distribution SendToEthMethod transaction.
	EventTypeSendToEth = "SendToEth"
)

// EmitSendToEthEvent creates a new event emitted on a SendToEthMethod transaction.
func (p Precompile) EmitSendToEthEvent(ctx sdk.Context, stateDB vm.StateDB, caller common.Address) error {
	// Prepare the event topics
	event := p.ABI.Events[EventTypeSendToEth]
	topics := make([]common.Hash, 2)

	// The first topic is always the signature of the event.
	topics[0] = event.ID

	var err error
	topics[1], err = cmn.MakeTopic(caller)
	if err != nil {
		return err
	}

	stateDB.AddLog(&ethtypes.Log{
		Address:     p.Address(),
		Topics:      topics,
		Data:        []byte{},
		BlockNumber: uint64(ctx.BlockHeight()),
	})

	return nil
}

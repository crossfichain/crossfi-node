package gravity

import (
	gravitykeeper "github.com/mineplexio/mineplex-2-node/x/gravity/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/vm"
)

const (
	// SendToEthMethod defines the ABI method name for the send to eth transaction.
	SendToEthMethod = "setToEth"
)

// SendToEth sends a transaction to the Ethereum network.
func (p Precompile) SendToEth(
	ctx sdk.Context,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	msg, err := NewMsgSendToEth(args, contract.Caller())
	if err != nil {
		return nil, err
	}

	// todo: convert msg.denom from contract.Caller() to cosmos-sdk part

	msgSrv := gravitykeeper.NewMsgServerImpl(p.gravityKeeper)
	if _, err = msgSrv.SendToEth(sdk.WrapSDKContext(ctx), msg); err != nil {
		return nil, err
	}

	if err = p.EmitSendToEthEvent(ctx, stateDB, contract.CallerAddress); err != nil {
		return nil, err
	}

	return method.Outputs.Pack(true)
}

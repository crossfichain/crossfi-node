package gravity

import (
	"fmt"
	gravitykeeper "github.com/mineplexio/mineplex-2-node/x/gravity/keeper"

	cmn "github.com/evmos/evmos/v13/precompiles/common"

	"github.com/ethereum/go-ethereum/common"

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
	origin common.Address,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	msg, senderAddress, err := NewMsgSendToEth(args)
	if err != nil {
		return nil, err
	}

	// todo: check
	// If the contract is the delegator, we don't need an origin check
	// Otherwise check if the origin matches the delegator address
	isContractDelegator := contract.CallerAddress == senderAddress
	if !isContractDelegator && origin != senderAddress {
		return nil, fmt.Errorf(cmn.ErrDifferentOrigin, origin.String(), senderAddress.String())
	}

	msgSrv := gravitykeeper.NewMsgServerImpl(p.gravityKeeper)
	if _, err = msgSrv.SendToEth(sdk.WrapSDKContext(ctx), msg); err != nil {
		return nil, err
	}

	if err = p.EmitSendToEthEvent(ctx, stateDB, senderAddress); err != nil {
		return nil, err
	}

	return method.Outputs.Pack(true)
}

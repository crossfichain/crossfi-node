package gravity

import (
	erc20types "github.com/mineplexio/mineplex-2-node/x/erc20/types"
	gravitykeeper "github.com/mineplexio/mineplex-2-node/x/gravity/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/vm"
)

const (
	// SendToEthMethod defines the ABI method name for the send to eth transaction.
	SendToEthMethod = "sendToEth"
)

// SendToEth sends a transaction to the Ethereum network.
func (p Precompile) SendToEth(ctx sdk.Context, evm *vm.EVM, contract *vm.Contract, stateDB vm.StateDB, method *abi.Method, args []interface{}) ([]byte, error) {
	msg, err := NewMsgSendToEth(args, contract.Caller())
	if err != nil {
		return nil, err
	}

	pair, err := p.erc20Keeper.TokenPair(sdk.WrapSDKContext(ctx), &erc20types.QueryTokenPairRequest{
		Token: msg.Amount.Denom,
	})
	if err != nil {
		return nil, err
	}

	_, err = p.erc20Keeper.ConvertERC20(sdk.WrapSDKContext(ctx.WithValue("evm", evm)), &erc20types.MsgConvertERC20{
		ContractAddress: pair.TokenPair.Erc20Address,
		Amount:          msg.Amount.Amount.Add(msg.ChainFee.Amount).Add(msg.BridgeFee.Amount),
		Receiver:        msg.Sender,
		Sender:          contract.Caller().String(),
	})
	if err != nil {
		return nil, err
	}

	msgSrv := gravitykeeper.NewMsgServerImpl(p.gravityKeeper)
	if _, err = msgSrv.SendToEth(sdk.WrapSDKContext(ctx), msg); err != nil {
		return nil, err
	}

	if err = p.EmitSendToEthEvent(ctx, stateDB, contract.CallerAddress); err != nil {
		return nil, err
	}

	return method.Outputs.Pack(true)
}

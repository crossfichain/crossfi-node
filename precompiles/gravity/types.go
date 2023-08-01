package gravity

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	cmn "github.com/evmos/evmos/v13/precompiles/common"
	gravitytypes "github.com/mineplexio/mineplex-2-node/x/gravity/types"
	"math/big"
)

// EventSendToEth defines the event data for the SendToEth transaction.
type EventSendToEth struct {
	Caller common.Address
}

// NewMsgSendToEth creates a new MsgSendToEth instance.
func NewMsgSendToEth(args []interface{}, caller common.Address) (*gravitytypes.MsgSendToEth, error) {
	if len(args) != 6 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 2, len(args))
	}

	chainId, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf(ErrInvalidChainID, args[0])
	}

	ethDest, ok := args[1].(common.Address)
	if !ok {
		return nil, fmt.Errorf(ErrInvalidDest, args[1])
	}

	denom, ok := args[2].(string)
	if !ok {
		return nil, fmt.Errorf(ErrInvalidDenom, args[2])
	}

	amount, ok := args[3].(*big.Int)
	if !ok {
		return nil, fmt.Errorf(ErrInvalidAmount, args[3])
	}

	bridgeFee, ok := args[4].(*big.Int)
	if !ok {
		return nil, fmt.Errorf(ErrInvalidAmount, args[4])
	}

	chainFee, ok := args[5].(*big.Int)
	if !ok {
		return nil, fmt.Errorf(ErrInvalidAmount, args[5])
	}

	msg := &gravitytypes.MsgSendToEth{
		ChainId:   chainId,
		Sender:    sdk.AccAddress(caller.Bytes()).String(),
		EthDest:   ethDest.String(),
		Amount:    sdk.NewCoin(denom, sdk.NewIntFromBigInt(amount)),
		BridgeFee: sdk.NewCoin(denom, sdk.NewIntFromBigInt(bridgeFee)),
		ChainFee:  sdk.NewCoin(denom, sdk.NewIntFromBigInt(chainFee)),
	}

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	return msg, nil
}

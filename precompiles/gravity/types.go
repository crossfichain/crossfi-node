package gravity

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	cmn "github.com/evmos/evmos/v13/precompiles/common"
	gravitytypes "github.com/mineplexio/mineplex-2-node/x/gravity/types"
)

// EventSendToEth defines the event data for the SendToEth transaction.
type EventSendToEth struct {
	Caller common.Address
}

// NewMsgSendToEth creates a new MsgSendToEth instance.
func NewMsgSendToEth(args []interface{}) (*gravitytypes.MsgSendToEth, common.Address, error) {
	if len(args) != 2 {
		return nil, common.Address{}, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 2, len(args))
	}

	delegatorAddress, ok := args[0].(common.Address)
	if !ok || delegatorAddress == (common.Address{}) {
		return nil, common.Address{}, fmt.Errorf(cmn.ErrInvalidDelegator, args[0])
	}

	withdrawerAddress, _ := args[1].(string)

	// If the withdrawer address is a hex address, convert it to a bech32 address.
	if common.IsHexAddress(withdrawerAddress) {
		var err error
		withdrawerAddress, err = sdk.Bech32ifyAddressBytes("evmos", common.HexToAddress(withdrawerAddress).Bytes())
		if err != nil {
			return nil, common.Address{}, err
		}
	}

	//DelegatorAddress: sdk.AccAddress(delegatorAddress.Bytes()).String(),
	//WithdrawAddress:  withdrawerAddress,

	// todo
	msg := &gravitytypes.MsgSendToEth{
		ChainId:   "",
		Sender:    "",
		EthDest:   "",
		Amount:    sdk.Coin{},
		BridgeFee: sdk.Coin{},
		ChainFee:  sdk.Coin{},
	}

	if err := msg.ValidateBasic(); err != nil {
		return nil, common.Address{}, err
	}

	return msg, delegatorAddress, nil
}

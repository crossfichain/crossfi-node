package post

import (
	"errors"
	feemarketkeeper "github.com/evmos/evmos/v12/x/feemarket/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

type HandlerOptions struct {
	FeeCollectorName string
	BankKeeper       bankkeeper.Keeper
	FeeMarketKeeper  *feemarketkeeper.Keeper
}

func (h HandlerOptions) Validate() error {
	if h.FeeCollectorName == "" {
		return errors.New("fee collector name cannot be empty")
	}

	if h.BankKeeper == nil {
		return errors.New("bank keeper cannot be nil")
	}
	if h.FeeMarketKeeper == nil {
		return errors.New("feemarket keeper cannot be nil")
	}

	return nil
}

func NewPostHandler(ho HandlerOptions) sdk.AnteHandler {
	postDecorators := []sdk.AnteDecorator{
		NewBurnDecorator(ho.FeeCollectorName, ho.BankKeeper, *ho.FeeMarketKeeper),
	}

	return sdk.ChainAnteDecorators(postDecorators...)
}

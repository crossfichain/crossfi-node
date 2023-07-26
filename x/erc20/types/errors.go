// Copyright 2022 Evmos Foundation
// This file is part of the Evmos Network packages.
//
// Evmos is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Evmos packages are distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Evmos packages. If not, see https://github.com/evmos/evmos/blob/main/LICENSE

package types

import (
	errorsmod "cosmossdk.io/errors"
)

// errors
var (
	ErrERC20Disabled          = errorsmod.Register(ModuleName, 101, "erc20 module is disabled")
	ErrInternalTokenPair      = errorsmod.Register(ModuleName, 102, "internal ethereum token mapping error")
	ErrTokenPairNotFound      = errorsmod.Register(ModuleName, 104, "token pair not found")
	ErrTokenPairAlreadyExists = errorsmod.Register(ModuleName, 105, "token pair already exists")
	ErrUndefinedOwner         = errorsmod.Register(ModuleName, 106, "undefined owner of contract pair")
	ErrBalanceInvariance      = errorsmod.Register(ModuleName, 107, "post transfer balance invariant failed")
	ErrUnexpectedEvent        = errorsmod.Register(ModuleName, 108, "unexpected event")
	ErrABIPack                = errorsmod.Register(ModuleName, 109, "contract ABI pack failed")
	ErrABIUnpack              = errorsmod.Register(ModuleName, 110, "contract ABI unpack failed")
	ErrEVMDenom               = errorsmod.Register(ModuleName, 111, "EVM denomination registration")
	ErrEVMCall                = errorsmod.Register(ModuleName, 112, "EVM call unexpected error")
	ErrERC20TokenPairDisabled = errorsmod.Register(ModuleName, 113, "erc20 token pair is disabled")
)

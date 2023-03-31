package types

import (
	"errors"
	"github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var (
	KeyOwner = []byte("Owner")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(owner string) Params {
	return Params{
		owner,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	// todo: set real address
	return NewParams("mp1d0ga6s7ue244rep5z7gnmgeyf3ejzla0aw5rur")
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyOwner, &p.Owner, validateOwner),
	}
}

func (p Params) ParseOwner() (types.AccAddress, error) {
	owner, err := types.AccAddressFromBech32(p.Owner)
	if err != nil {
		return nil, err
	}

	return owner, nil
}

func validateOwner(value interface{}) error {
	addr, ok := value.(string)
	if !ok {
		return errors.New("cannot cast owner addr to string")
	}

	_, err := types.AccAddressFromBech32(addr)
	if err != nil {
		return err
	}

	return nil
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateOwner(p.Owner); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

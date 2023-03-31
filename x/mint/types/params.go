package types

import (
	"cosmossdk.io/math"
	"errors"
	"fmt"
	"strings"

	"sigs.k8s.io/yaml"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	KeyMintDenom = []byte("MintDenom")
	KeyPeriods   = []byte("Periods")
)

// ParamTable for minting module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(mintDenom string, periods []*RewardPeriod) Params {
	return Params{
		MintDenom: mintDenom,
		Periods:   periods,
	}
}

func MustNewIntFromString(s string) math.Int {
	res, success := sdk.NewIntFromString(s)
	if !success {
		panic("invalid string for NewIntFromString: " + s)
	}

	return res
}

// DefaultParams default minting module parameters
func DefaultParams() Params {
	return Params{
		MintDenom: "xfi",
		Periods: []*RewardPeriod{
			{
				FromHeight:     1,
				ToHeight:       25228800,
				RewardPerBlock: MustNewIntFromString("5000000000000000000"),
			},
			{
				FromHeight:     25228801,
				ToHeight:       50457600,
				RewardPerBlock: MustNewIntFromString("4000000000000000000"),
			},
			{
				FromHeight:     50457601,
				ToHeight:       75686400,
				RewardPerBlock: MustNewIntFromString("3000000000000000000"),
			},
			{
				FromHeight:     75686401,
				ToHeight:       100915200,
				RewardPerBlock: MustNewIntFromString("2000000000000000000"),
			},
			{
				FromHeight:     100915201,
				ToHeight:       126144000,
				RewardPerBlock: MustNewIntFromString("1000000000000000000"),
			},
		},
	}
}

// validate params
func (p Params) Validate() error {
	if err := validateMintDenom(p.MintDenom); err != nil {
		return err
	}
	if err := validatePeriods(p.Periods); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMintDenom, &p.MintDenom, validateMintDenom),
		paramtypes.NewParamSetPair(KeyPeriods, &p.Periods, validatePeriods),
	}
}

func validateMintDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("mint denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func validatePeriods(i interface{}) error {
	v, ok := i.([]*RewardPeriod)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, periodA := range v {
		if periodA.RewardPerBlock.IsNegative() {
			return fmt.Errorf("negative reward per block")
		}

		// check that periodA and periodB are not overlapping each other
		for _, periodB := range v {
			if periodA == periodB {
				continue
			}

			if periodA.FromHeight < periodB.FromHeight {
				if periodA.ToHeight > periodB.ToHeight {
					return fmt.Errorf("reward periods are overlapping: %s and %s", periodA, periodB)
				}
			} else {
				if periodB.ToHeight > (periodA.FromHeight) {
					return fmt.Errorf("reward periods are overlapping: %s and %s", periodA, periodB)
				}
			}
		}
	}

	return nil
}

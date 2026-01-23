package types

import (
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// DefaultGenesis returns the default genesis state.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Models: []*Model{},
	}
}

func ValidateGenesis(data GenesisState) error {
	return data.ValidateBasic()
}

func (m *Model) ValidateBasic() error {
	if m == nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "model is nil")
	}
	if m.Id == "" {
		return ErrInvalidID
	}
	if m.Name == "" {
		return ErrInvalidName
	}
	if m.Url == "" {
		return ErrInvalidURL
	}
	if m.Creator != "" {
		if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil {
			return errorsmod.Wrap(err, "creator")
		}
	}
	return nil
}

func (gs GenesisState) ValidateBasic() error {
	for i, m := range gs.Models {
		if err := m.ValidateBasic(); err != nil {
			return errorsmod.Wrapf(err, "models[%d]", i)
		}
	}
	return nil
}

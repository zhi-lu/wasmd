package types

import (
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// DefaultGenesis returns the default genesis state.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Models:      []Model{},
		Agents:      []Agent{},
		Tasks:       []Task{},
		NextAgentId: 1,
		NextTaskId:  1,
	}
}

func ValidateGenesis(data GenesisState) error {
	return data.ValidateBasic()
}

// ──────── Model validation ────────

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

// ──────── Agent validation ────────

func (a *Agent) ValidateBasic() error {
	if a == nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "agent is nil")
	}
	if a.Id == "" {
		return ErrInvalidID
	}
	if a.Name == "" {
		return ErrInvalidName
	}
	if a.ModelId == "" {
		return ErrInvalidModelRef
	}
	if a.Operator != "" {
		if _, err := sdk.AccAddressFromBech32(a.Operator); err != nil {
			return errorsmod.Wrap(ErrInvalidOperator, err.Error())
		}
	}
	if !a.FeePerTask.IsValid() {
		return ErrInvalidFee
	}
	return nil
}

// ──────── Task validation ────────

func (t *Task) ValidateBasic() error {
	if t == nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "task is nil")
	}
	if t.Id == "" {
		return ErrInvalidID
	}
	if t.AgentId == "" {
		return ErrInvalidID
	}
	if t.Creator != "" {
		if _, err := sdk.AccAddressFromBech32(t.Creator); err != nil {
			return errorsmod.Wrap(ErrInvalidCreator, err.Error())
		}
	}
	return nil
}

// ──────── GenesisState validation ────────

func (gs GenesisState) ValidateBasic() error {
	for i, m := range gs.Models {
		if err := m.ValidateBasic(); err != nil {
			return errorsmod.Wrapf(err, "models[%d]", i)
		}
	}
	for i, a := range gs.Agents {
		if err := a.ValidateBasic(); err != nil {
			return errorsmod.Wrapf(err, "agents[%d]", i)
		}
	}
	for i, t := range gs.Tasks {
		if err := t.ValidateBasic(); err != nil {
			return errorsmod.Wrapf(err, "tasks[%d]", i)
		}
	}
	return nil
}

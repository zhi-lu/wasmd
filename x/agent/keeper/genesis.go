package keeper

import (
	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/CosmWasm/wasmd/x/agent/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, gs *types.GenesisState) error {
	if gs == nil {
		return errors.Wrap(types.ErrInvalidName, "genesis state is nil")
	}
	for _, m := range gs.Models {
		if m == nil {
			return errors.Wrap(types.ErrInvalidName, "nil model in genesis")
		}
		model := *m
		if err := model.ValidateBasic(); err != nil {
			return err
		}
		if model.Creator == "" {
			return types.ErrInvalidCreator
		}
		if err := k.RegisterModel(ctx, model); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) ExportGenesis(ctx sdk.Context) (*types.GenesisState, error) {
	models, err := k.GetModels(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*types.Model, 0, len(models))
	for i := range models {
		m := models[i]
		out = append(out, &m)
	}
	return &types.GenesisState{Models: out}, nil
}

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

	// Models
	for i, m := range gs.Models {
		if err := m.ValidateBasic(); err != nil {
			return errors.Wrapf(err, "models[%d]", i)
		}
		if m.Creator == "" {
			return types.ErrInvalidCreator
		}
		if err := k.RegisterModel(ctx, m); err != nil {
			return err
		}
	}

	// Agents
	for i, a := range gs.Agents {
		if err := a.ValidateBasic(); err != nil {
			return errors.Wrapf(err, "agents[%d]", i)
		}
		if err := k.SetAgent(ctx, a); err != nil {
			return err
		}
	}
	if err := k.SetNextAgentID(ctx, gs.NextAgentId); err != nil {
		return err
	}

	// Tasks
	for i, t := range gs.Tasks {
		if err := t.ValidateBasic(); err != nil {
			return errors.Wrapf(err, "tasks[%d]", i)
		}
		if err := k.SetTask(ctx, t); err != nil {
			return err
		}
	}
	if err := k.SetNextTaskID(ctx, gs.NextTaskId); err != nil {
		return err
	}

	return nil
}

func (k Keeper) ExportGenesis(ctx sdk.Context) (*types.GenesisState, error) {
	models, err := k.GetModels(ctx)
	if err != nil {
		return nil, err
	}
	agents, err := k.GetAgents(ctx)
	if err != nil {
		return nil, err
	}
	tasks, err := k.GetTasks(ctx)
	if err != nil {
		return nil, err
	}
	nextAgentID, err := k.GetNextAgentID(ctx)
	if err != nil {
		return nil, err
	}
	nextTaskID, err := k.GetNextTaskID(ctx)
	if err != nil {
		return nil, err
	}

	return &types.GenesisState{
		Models:      models,
		Agents:      agents,
		Tasks:       tasks,
		NextAgentId: nextAgentID,
		NextTaskId:  nextTaskID,
	}, nil
}

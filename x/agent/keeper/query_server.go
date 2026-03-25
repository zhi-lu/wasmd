package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/CosmWasm/wasmd/x/agent/types"
)

type queryServer struct {
	Keeper
}

var _ types.QueryServer = queryServer{}

func NewQueryServer(k Keeper) types.QueryServer {
	return queryServer{Keeper: k}
}

// ── Model ──

func (q queryServer) Model(goCtx context.Context, req *types.QueryModelRequest) (*types.QueryModelResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if req == nil {
		return nil, types.ErrInvalidID
	}
	model, err := q.Keeper.GetModel(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &types.QueryModelResponse{Model: &model}, nil
}

func (q queryServer) Models(goCtx context.Context, _ *types.QueryModelsRequest) (*types.QueryModelsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	models, err := q.Keeper.GetModels(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*types.Model, 0, len(models))
	for i := range models {
		m := models[i]
		out = append(out, &m)
	}
	return &types.QueryModelsResponse{Models: out}, nil
}

// ── Agent ──

func (q queryServer) Agent(goCtx context.Context, req *types.QueryAgentRequest) (*types.QueryAgentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if req == nil {
		return nil, types.ErrInvalidID
	}
	agent, err := q.Keeper.GetAgent(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &types.QueryAgentResponse{Agent: &agent}, nil
}

func (q queryServer) Agents(goCtx context.Context, _ *types.QueryAgentsRequest) (*types.QueryAgentsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	agents, err := q.Keeper.GetAgents(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*types.Agent, 0, len(agents))
	for i := range agents {
		a := agents[i]
		out = append(out, &a)
	}
	return &types.QueryAgentsResponse{Agents: out}, nil
}

func (q queryServer) AgentsByModel(goCtx context.Context, req *types.QueryAgentsByModelRequest) (*types.QueryAgentsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if req == nil {
		return nil, types.ErrInvalidID
	}
	agents, err := q.Keeper.GetAgentsByModel(ctx, req.ModelId)
	if err != nil {
		return nil, err
	}
	out := make([]*types.Agent, 0, len(agents))
	for i := range agents {
		a := agents[i]
		out = append(out, &a)
	}
	return &types.QueryAgentsResponse{Agents: out}, nil
}

// ── Task ──

func (q queryServer) Task(goCtx context.Context, req *types.QueryTaskRequest) (*types.QueryTaskResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if req == nil {
		return nil, types.ErrInvalidID
	}
	task, err := q.Keeper.GetTask(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &types.QueryTaskResponse{Task: &task}, nil
}

func (q queryServer) Tasks(goCtx context.Context, _ *types.QueryTasksRequest) (*types.QueryTasksResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	tasks, err := q.Keeper.GetTasks(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*types.Task, 0, len(tasks))
	for i := range tasks {
		t := tasks[i]
		out = append(out, &t)
	}
	return &types.QueryTasksResponse{Tasks: out}, nil
}

func (q queryServer) TasksByAgent(goCtx context.Context, req *types.QueryTasksByAgentRequest) (*types.QueryTasksResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if req == nil {
		return nil, types.ErrInvalidID
	}
	tasks, err := q.Keeper.GetTasksByAgent(ctx, req.AgentId)
	if err != nil {
		return nil, err
	}
	out := make([]*types.Task, 0, len(tasks))
	for i := range tasks {
		t := tasks[i]
		out = append(out, &t)
	}
	return &types.QueryTasksResponse{Tasks: out}, nil
}

func (q queryServer) TasksByCreator(goCtx context.Context, req *types.QueryTasksByCreatorRequest) (*types.QueryTasksResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if req == nil {
		return nil, types.ErrInvalidID
	}
	tasks, err := q.Keeper.GetTasksByCreator(ctx, req.Creator)
	if err != nil {
		return nil, err
	}
	out := make([]*types.Task, 0, len(tasks))
	for i := range tasks {
		t := tasks[i]
		out = append(out, &t)
	}
	return &types.QueryTasksResponse{Tasks: out}, nil
}

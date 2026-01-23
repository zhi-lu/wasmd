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

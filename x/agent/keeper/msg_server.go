package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	"github.com/CosmWasm/wasmd/x/agent/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = (*msgServer)(nil)

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// ── Model ──

func (m msgServer) RegisterModel(goCtx context.Context, msg *types.MsgRegisterModel) (*types.MsgRegisterModelResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg == nil {
		return nil, errorsmod.Wrap(types.ErrInvalidCreator, "msg is nil")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return nil, types.ErrInvalidCreator
	}
	if msg.Model == nil {
		return nil, errorsmod.Wrap(types.ErrInvalidName, "model is nil")
	}
	if err := msg.Model.ValidateBasic(); err != nil {
		return nil, err
	}
	if msg.Model.Creator == "" {
		msg.Model.Creator = msg.Creator
	}
	if msg.Creator != msg.Model.Creator {
		return nil, types.ErrInvalidCreatorForEqual
	}
	if err := m.Keeper.RegisterModel(ctx, *msg.Model); err != nil {
		return nil, err
	}
	return &types.MsgRegisterModelResponse{}, nil
}

func (m msgServer) DeleteModel(goCtx context.Context, msg *types.MsgDeleteModel) (*types.MsgDeleteModelResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg == nil {
		return nil, errorsmod.Wrap(types.ErrInvalidCreator, "msg is nil")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return nil, types.ErrInvalidCreator
	}
	if msg.ModelId == "" {
		return nil, types.ErrInvalidID
	}
	if err := m.Keeper.DeleteModel(ctx, msg.ModelId, msg.Creator); err != nil {
		return nil, err
	}
	return &types.MsgDeleteModelResponse{}, nil
}

// ── Agent ──

func (m msgServer) RegisterAgent(goCtx context.Context, msg *types.MsgRegisterAgent) (*types.MsgRegisterAgentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg == nil {
		return nil, errorsmod.Wrap(types.ErrInvalidOperator, "msg is nil")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Operator); err != nil {
		return nil, types.ErrInvalidOperator
	}
	if msg.Name == "" {
		return nil, types.ErrInvalidName
	}
	if msg.ModelId == "" {
		return nil, types.ErrInvalidModelRef
	}
	agentID, err := m.Keeper.RegisterAgent(ctx, *msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgRegisterAgentResponse{AgentId: agentID}, nil
}

func (m msgServer) UpdateAgent(goCtx context.Context, msg *types.MsgUpdateAgent) (*types.MsgUpdateAgentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg == nil {
		return nil, errorsmod.Wrap(types.ErrInvalidOperator, "msg is nil")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Operator); err != nil {
		return nil, types.ErrInvalidOperator
	}
	if msg.AgentId == "" {
		return nil, types.ErrInvalidID
	}
	if err := m.Keeper.UpdateAgent(ctx, *msg); err != nil {
		return nil, err
	}
	return &types.MsgUpdateAgentResponse{}, nil
}

func (m msgServer) DeactivateAgent(goCtx context.Context, msg *types.MsgDeactivateAgent) (*types.MsgDeactivateAgentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg == nil {
		return nil, errorsmod.Wrap(types.ErrInvalidOperator, "msg is nil")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Operator); err != nil {
		return nil, types.ErrInvalidOperator
	}
	if msg.AgentId == "" {
		return nil, types.ErrInvalidID
	}
	if err := m.Keeper.DeactivateAgent(ctx, msg.Operator, msg.AgentId); err != nil {
		return nil, err
	}
	return &types.MsgDeactivateAgentResponse{}, nil
}

// ── Task ──

func (m msgServer) CreateTask(goCtx context.Context, msg *types.MsgCreateTask) (*types.MsgCreateTaskResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg == nil {
		return nil, errorsmod.Wrap(types.ErrInvalidCreator, "msg is nil")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return nil, types.ErrInvalidCreator
	}
	if msg.AgentId == "" {
		return nil, types.ErrInvalidID
	}
	if msg.InputHash == "" {
		return nil, types.ErrInvalidInputHash
	}
	taskID, err := m.Keeper.CreateTask(ctx, *msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgCreateTaskResponse{TaskId: taskID}, nil
}

func (m msgServer) SubmitTaskResult(goCtx context.Context, msg *types.MsgSubmitTaskResult) (*types.MsgSubmitTaskResultResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg == nil {
		return nil, errorsmod.Wrap(types.ErrInvalidOperator, "msg is nil")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Operator); err != nil {
		return nil, types.ErrInvalidOperator
	}
	if msg.TaskId == "" {
		return nil, types.ErrInvalidID
	}
	if err := m.Keeper.SubmitTaskResult(ctx, *msg); err != nil {
		return nil, err
	}
	return &types.MsgSubmitTaskResultResponse{}, nil
}

func (m msgServer) CancelTask(goCtx context.Context, msg *types.MsgCancelTask) (*types.MsgCancelTaskResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg == nil {
		return nil, errorsmod.Wrap(types.ErrInvalidCreator, "msg is nil")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return nil, types.ErrInvalidCreator
	}
	if msg.TaskId == "" {
		return nil, types.ErrInvalidID
	}
	if err := m.Keeper.CancelTask(ctx, msg.Creator, msg.TaskId); err != nil {
		return nil, err
	}
	return &types.MsgCancelTaskResponse{}, nil
}

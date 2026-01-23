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

func (m msgServer) RegisterModel(goCtx context.Context, msg *types.MsgRegisterModel) (*types.MsgRegisterModelResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg == nil {
		return nil, errorsmod.Wrap(types.ErrInvalidCreator, "msg is nil")
	}

	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
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

	err = m.Keeper.RegisterModel(ctx, *msg.Model)
	if err != nil {
		return nil, err
	}
	return &types.MsgRegisterModelResponse{}, nil
}

func (m msgServer) DeleteModel(goCtx context.Context, msg *types.MsgDeleteModel) (*types.MsgDeleteModelResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg == nil {
		return nil, errorsmod.Wrap(types.ErrInvalidCreator, "msg is nil")
	}

	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}
	if msg.ModelId == "" {
		return nil, types.ErrInvalidID
	}

	err = m.Keeper.DeleteModel(ctx, msg.ModelId, msg.Creator)
	if err != nil {
		return nil, err
	}

	return &types.MsgDeleteModelResponse{}, nil
}

package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgRegisterModel{}, "agent/RegisterModel", nil)
	cdc.RegisterConcrete(&MsgDeleteModel{}, "agent/DeleteModel", nil)
	cdc.RegisterConcrete(&MsgRegisterAgent{}, "agent/RegisterAgent", nil)
	cdc.RegisterConcrete(&MsgUpdateAgent{}, "agent/UpdateAgent", nil)
	cdc.RegisterConcrete(&MsgDeactivateAgent{}, "agent/DeactivateAgent", nil)
	cdc.RegisterConcrete(&MsgCreateTask{}, "agent/CreateTask", nil)
	cdc.RegisterConcrete(&MsgSubmitTaskResult{}, "agent/SubmitTaskResult", nil)
	cdc.RegisterConcrete(&MsgCancelTask{}, "agent/CancelTask", nil)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgRegisterModel{},
		&MsgDeleteModel{},
		&MsgRegisterAgent{},
		&MsgUpdateAgent{},
		&MsgDeactivateAgent{},
		&MsgCreateTask{},
		&MsgSubmitTaskResult{},
		&MsgCancelTask{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

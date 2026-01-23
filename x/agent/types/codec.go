// types/codec.go
package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary types on the provided LegacyAmino codec.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// 可选，注册 Msg 类型（如果使用 amino 签名）
	cdc.RegisterConcrete(&MsgRegisterModel{}, "agent/RegisterModel", nil)
	cdc.RegisterConcrete(&MsgDeleteModel{}, "agent/DeleteModel", nil)
}

// RegisterInterfaces registers the interface types used by the agent module.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgRegisterModel{},
		&MsgDeleteModel{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

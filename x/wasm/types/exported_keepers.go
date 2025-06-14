package types

import (
	"context"

	wasmvmtypes "github.com/CosmWasm/wasmvm/v3/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ViewKeeper provides read only operations
type ViewKeeper interface {
	GetContractHistory(ctx context.Context, contractAddr sdk.AccAddress) []ContractCodeHistoryEntry
	QuerySmart(ctx context.Context, contractAddr sdk.AccAddress, req []byte) ([]byte, error)
	QueryRaw(ctx context.Context, contractAddress sdk.AccAddress, key []byte) []byte
	HasContractInfo(ctx context.Context, contractAddress sdk.AccAddress) bool
	GetContractInfo(ctx context.Context, contractAddress sdk.AccAddress) *ContractInfo
	IterateContractInfo(ctx context.Context, cb func(sdk.AccAddress, ContractInfo) bool)
	IterateContractsByCreator(ctx context.Context, creator sdk.AccAddress, cb func(address sdk.AccAddress) bool)
	IterateContractsByCode(ctx context.Context, codeID uint64, cb func(address sdk.AccAddress) bool)
	IterateContractState(ctx context.Context, contractAddress sdk.AccAddress, cb func(key, value []byte) bool)
	GetCodeInfo(ctx context.Context, codeID uint64) *CodeInfo
	IterateCodeInfos(ctx context.Context, cb func(uint64, CodeInfo) bool)
	GetByteCode(ctx context.Context, codeID uint64) ([]byte, error)
	IsPinnedCode(ctx context.Context, codeID uint64) bool
	GetParams(ctx context.Context) Params
	GetWasmLimits() wasmvmtypes.WasmLimits
}

// ContractOpsKeeper contains mutable operations on a contract.
type ContractOpsKeeper interface {
	// Create uploads and compiles a WASM contract, returning a short identifier for the contract
	Create(ctx sdk.Context, creator sdk.AccAddress, wasmCode []byte, instantiateAccess *AccessConfig) (codeID uint64, checksum []byte, err error)

	// Instantiate creates an instance of a WASM contract using the classic sequence based address generator
	Instantiate(
		ctx sdk.Context,
		codeID uint64,
		creator, admin sdk.AccAddress,
		initMsg []byte,
		label string,
		deposit sdk.Coins,
	) (sdk.AccAddress, []byte, error)

	// Instantiate2 creates an instance of a WASM contract using the predictable address generator
	Instantiate2(
		ctx sdk.Context,
		codeID uint64,
		creator, admin sdk.AccAddress,
		initMsg []byte,
		label string,
		deposit sdk.Coins,
		salt []byte,
		fixMsg bool,
	) (sdk.AccAddress, []byte, error)

	// Execute executes the contract instance
	Execute(ctx sdk.Context, contractAddress, caller sdk.AccAddress, msg []byte, coins sdk.Coins) ([]byte, error)

	// Migrate allows to upgrade a contract to a new code with data migration.
	Migrate(ctx sdk.Context, contractAddress, caller sdk.AccAddress, newCodeID uint64, msg []byte) ([]byte, error)

	// Sudo allows to call privileged entry point of a contract.
	Sudo(ctx sdk.Context, contractAddress sdk.AccAddress, msg []byte) ([]byte, error)

	// UpdateContractAdmin sets the admin value on the ContractInfo. It must be a valid address (use ClearContractAdmin to remove it)
	UpdateContractAdmin(ctx sdk.Context, contractAddress, caller, newAdmin sdk.AccAddress) error

	// ClearContractAdmin sets the admin value on the ContractInfo to nil, to disable further migrations/ updates.
	ClearContractAdmin(ctx sdk.Context, contractAddress, caller sdk.AccAddress) error

	// PinCode pins the wasm contract in wasmvm cache
	PinCode(ctx sdk.Context, codeID uint64) error

	// UnpinCode removes the wasm contract from wasmvm cache
	UnpinCode(ctx sdk.Context, codeID uint64) error

	// SetContractInfoExtension updates the extension point data that is stored with the contract info
	SetContractInfoExtension(ctx sdk.Context, contract sdk.AccAddress, extra ContractInfoExtension) error

	// SetAccessConfig updates the access config of a code id.
	SetAccessConfig(ctx sdk.Context, codeID uint64, caller sdk.AccAddress, newConfig AccessConfig) error
}

// IBCContractKeeper IBC lifecycle event handler
type IBCContractKeeper interface {
	OnOpenChannel(
		ctx sdk.Context,
		contractAddr sdk.AccAddress,
		msg wasmvmtypes.IBCChannelOpenMsg,
	) (string, error)
	OnConnectChannel(
		ctx sdk.Context,
		contractAddr sdk.AccAddress,
		msg wasmvmtypes.IBCChannelConnectMsg,
	) error
	OnCloseChannel(
		ctx sdk.Context,
		contractAddr sdk.AccAddress,
		msg wasmvmtypes.IBCChannelCloseMsg,
	) error
	OnRecvPacket(
		ctx sdk.Context,
		contractAddr sdk.AccAddress,
		msg wasmvmtypes.IBCPacketReceiveMsg,
	) (ibcexported.Acknowledgement, error)
	OnAckPacket(
		ctx sdk.Context,
		contractAddr sdk.AccAddress,
		acknowledgement wasmvmtypes.IBCPacketAckMsg,
	) error
	OnTimeoutPacket(
		ctx sdk.Context,
		contractAddr sdk.AccAddress,
		msg wasmvmtypes.IBCPacketTimeoutMsg,
	) error
	IBCSourceCallback(
		ctx sdk.Context,
		contractAddr sdk.AccAddress,
		msg wasmvmtypes.IBCSourceCallbackMsg,
	) error
	IBCDestinationCallback(
		ctx sdk.Context,
		contractAddr sdk.AccAddress,
		msg wasmvmtypes.IBCDestinationCallbackMsg,
	) error

	// LoadAsyncAckPacket loads a previously stored packet. See StoreAsyncAckPacket for more details.
	// Both the portID and channelID are the ones on the destination chain (the chain that this is executed on).
	LoadAsyncAckPacket(ctx context.Context, portID, channelID string, sequence uint64) (channeltypes.Packet, error)
	// StoreAsyncAckPacket stores a packet to be acknowledged later. These are packets that were
	// received and processed by the contract, but the contract did not want to acknowledge them immediately.
	// They are stored in the keeper until the contract calls "WriteAcknowledgement" to acknowledge them.
	StoreAsyncAckPacket(ctx context.Context, packet channeltypes.Packet) error
	// DeleteAsyncAckPacket deletes a previously stored packet. See StoreAsyncAckPacket for more details.
	DeleteAsyncAckPacket(ctx context.Context, portID, channelID string, sequence uint64)
}

// IBC2ContractKeeper IBC2 lifecycle event handler
type IBC2ContractKeeper interface {
	OnRecvIBC2Packet(
		ctx sdk.Context,
		contractAddr sdk.AccAddress,
		msg wasmvmtypes.IBC2PacketReceiveMsg,
	) channeltypesv2.RecvPacketResult
	OnTimeoutIBC2Packet(
		ctx sdk.Context,
		contractAddr sdk.AccAddress,
		msg wasmvmtypes.IBC2PacketTimeoutMsg,
	) error
}

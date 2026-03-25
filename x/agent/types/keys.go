package types

import (
	"encoding/binary"
	"fmt"
)

const (
	// 模块名
	ModuleName = "agent"

	// 状态存储 KVStore 的 key
	StoreKey = ModuleName

	// 用于 routing 消息
	RouterKey = ModuleName

	// 用于 query 请求
	QuerierRoute = ModuleName
)

var (
	ModelKeyPrefix = []byte("model:")
	AgentKeyPrefix = []byte("agent:")
	TaskKeyPrefix  = []byte("task:")

	AgentsByModelIndexPrefix  = []byte("agent_by_model:")
	TasksByAgentIndexPrefix   = []byte("task_by_agent:")
	TasksByCreatorIndexPrefix = []byte("task_by_creator:")

	NextAgentIDKey = []byte("next_agent_id")
	NextTaskIDKey  = []byte("next_task_id")
)

func GetModelKey(id string) []byte {
	return append(ModelKeyPrefix, []byte(id)...)
}

func GetAgentKey(id string) []byte {
	return append(AgentKeyPrefix, []byte(id)...)
}

func GetTaskKey(id string) []byte {
	return append(TaskKeyPrefix, []byte(id)...)
}

// GetAgentsByModelIndexKey returns a prefix key for indexing agents by model_id.
func GetAgentsByModelIndexKey(modelID, agentID string) []byte {
	return append(AgentsByModelIndexPrefix, []byte(fmt.Sprintf("%s/%s", modelID, agentID))...)
}

// GetTasksByAgentIndexKey returns a prefix key for indexing tasks by agent_id.
func GetTasksByAgentIndexKey(agentID, taskID string) []byte {
	return append(TasksByAgentIndexPrefix, []byte(fmt.Sprintf("%s/%s", agentID, taskID))...)
}

// GetTasksByCreatorIndexKey returns a prefix key for indexing tasks by creator.
func GetTasksByCreatorIndexKey(creator, taskID string) []byte {
	return append(TasksByCreatorIndexPrefix, []byte(fmt.Sprintf("%s/%s", creator, taskID))...)
}

// IDToBytes converts a uint64 to big-endian bytes for storage.
func IDToBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// BytesToID converts big-endian bytes back to uint64.
func BytesToID(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}

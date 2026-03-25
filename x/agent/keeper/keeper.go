package keeper

import (
	"context"
	"fmt"
	"strconv"

	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/store/prefix"
	"github.com/CosmWasm/wasmd/x/agent/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BankKeeper defines the expected bank module interface.
type BankKeeper interface {
	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

type Keeper struct {
	cdc          codec.Codec
	storeService corestoretypes.KVStoreService
	bankKeeper   BankKeeper
}

func NewKeeper(cdc codec.Codec, storeService corestoretypes.KVStoreService, bankKeeper BankKeeper) Keeper {
	return Keeper{
		cdc:          cdc,
		storeService: storeService,
		bankKeeper:   bankKeeper,
	}
}

// ═══════════════════════════════════════════
//  Model CRUD
// ═══════════════════════════════════════════

func (k Keeper) RegisterModel(ctx sdk.Context, model types.Model) error {
	store := k.storeService.OpenKVStore(ctx)
	if model.Id == "" {
		return types.ErrInvalidID
	}
	if model.Name == "" {
		return types.ErrInvalidName
	}
	if model.Url == "" {
		return types.ErrInvalidURL
	}
	key := types.GetModelKey(model.Id)
	has, err := store.Has(key)
	if err != nil {
		return err
	}
	if has {
		return types.ErrModelAlreadyExists{ID: model.Id}
	}
	bz, err := k.cdc.Marshal(&model)
	if err != nil {
		return err
	}
	return store.Set(key, bz)
}

func (k Keeper) DeleteModel(ctx sdk.Context, modelID string, creator string) error {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetModelKey(modelID)
	has, err := store.Has(key)
	if err != nil {
		return err
	}
	if !has {
		return types.ErrModelNotFound{ID: modelID}
	}
	bz, err := store.Get(key)
	if err != nil {
		return err
	}
	var model types.Model
	if err := k.cdc.Unmarshal(bz, &model); err != nil {
		return err
	}
	if creator != model.Creator {
		return types.ErrInvalidDeleteModel
	}
	return store.Delete(key)
}

func (k Keeper) GetModel(ctx sdk.Context, id string) (types.Model, error) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetModelKey(id)
	has, err := store.Has(key)
	if err != nil {
		return types.Model{}, err
	}
	if !has {
		return types.Model{}, types.ErrModelNotFound{ID: id}
	}
	bz, err := store.Get(key)
	if err != nil {
		return types.Model{}, err
	}
	var model types.Model
	if err := k.cdc.Unmarshal(bz, &model); err != nil {
		return types.Model{}, err
	}
	return model, nil
}

func (k Keeper) GetModels(ctx sdk.Context) ([]types.Model, error) {
	return iteratePrefix[types.Model](k, ctx, types.ModelKeyPrefix)
}

// ═══════════════════════════════════════════
//  Agent CRUD
// ═══════════════════════════════════════════

func (k Keeper) nextAgentID(ctx sdk.Context) (uint64, error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.NextAgentIDKey)
	if err != nil {
		return 0, err
	}
	if bz == nil {
		return 1, nil
	}
	return types.BytesToID(bz), nil
}

func (k Keeper) setNextAgentID(ctx sdk.Context, id uint64) error {
	store := k.storeService.OpenKVStore(ctx)
	return store.Set(types.NextAgentIDKey, types.IDToBytes(id))
}

func (k Keeper) RegisterAgent(ctx sdk.Context, msg types.MsgRegisterAgent) (string, error) {
	// 验证模型存在
	if _, err := k.GetModel(ctx, msg.ModelId); err != nil {
		return "", types.ErrInvalidModelRef
	}

	nextID, err := k.nextAgentID(ctx)
	if err != nil {
		return "", err
	}
	agentID := fmt.Sprintf("agent-%d", nextID)

	agent := types.Agent{
		Id:              agentID,
		Name:            msg.Name,
		ModelId:         msg.ModelId,
		Operator:        msg.Operator,
		Description:     msg.Description,
		FeePerTask:      msg.FeePerTask,
		Status:          types.AgentStatus_AGENT_STATUS_ACTIVE,
		ReputationScore: 0,
		TasksCompleted:  0,
	}

	store := k.storeService.OpenKVStore(ctx)
	bz, err := k.cdc.Marshal(&agent)
	if err != nil {
		return "", err
	}
	if err := store.Set(types.GetAgentKey(agentID), bz); err != nil {
		return "", err
	}
	// 建索引：按 model_id
	if err := store.Set(types.GetAgentsByModelIndexKey(msg.ModelId, agentID), []byte(agentID)); err != nil {
		return "", err
	}
	// 递增 ID
	if err := k.setNextAgentID(ctx, nextID+1); err != nil {
		return "", err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		"agent_registered",
		sdk.NewAttribute("agent_id", agentID),
		sdk.NewAttribute("operator", msg.Operator),
		sdk.NewAttribute("model_id", msg.ModelId),
	))

	return agentID, nil
}

func (k Keeper) GetAgent(ctx sdk.Context, id string) (types.Agent, error) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetAgentKey(id)
	has, err := store.Has(key)
	if err != nil {
		return types.Agent{}, err
	}
	if !has {
		return types.Agent{}, types.ErrAgentNotFoundError{ID: id}
	}
	bz, err := store.Get(key)
	if err != nil {
		return types.Agent{}, err
	}
	var agent types.Agent
	if err := k.cdc.Unmarshal(bz, &agent); err != nil {
		return types.Agent{}, err
	}
	return agent, nil
}

func (k Keeper) setAgent(ctx sdk.Context, agent types.Agent) error {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := k.cdc.Marshal(&agent)
	if err != nil {
		return err
	}
	return store.Set(types.GetAgentKey(agent.Id), bz)
}

func (k Keeper) UpdateAgent(ctx sdk.Context, msg types.MsgUpdateAgent) error {
	agent, err := k.GetAgent(ctx, msg.AgentId)
	if err != nil {
		return err
	}
	if agent.Operator != msg.Operator {
		return types.ErrUnauthorized
	}
	if msg.Name != "" {
		agent.Name = msg.Name
	}
	if msg.Description != "" {
		agent.Description = msg.Description
	}
	if msg.FeePerTask != nil && msg.FeePerTask.IsValid() {
		agent.FeePerTask = msg.FeePerTask
	}
	return k.setAgent(ctx, agent)
}

func (k Keeper) DeactivateAgent(ctx sdk.Context, operator, agentID string) error {
	agent, err := k.GetAgent(ctx, agentID)
	if err != nil {
		return err
	}
	if agent.Operator != operator {
		return types.ErrUnauthorized
	}
	agent.Status = types.AgentStatus_AGENT_STATUS_INACTIVE
	return k.setAgent(ctx, agent)
}

func (k Keeper) GetAgents(ctx sdk.Context) ([]types.Agent, error) {
	return iteratePrefix[types.Agent](k, ctx, types.AgentKeyPrefix)
}

func (k Keeper) GetAgentsByModel(ctx sdk.Context, modelID string) ([]types.Agent, error) {
	prefixBytes := append(types.AgentsByModelIndexPrefix, []byte(modelID+"/")...)
	prefixStore := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), prefixBytes)
	iter := prefixStore.Iterator(nil, nil)
	defer iter.Close()

	var agents []types.Agent
	for ; iter.Valid(); iter.Next() {
		agentID := string(iter.Value())
		agent, err := k.GetAgent(ctx, agentID)
		if err != nil {
			continue // index stale, skip
		}
		agents = append(agents, agent)
	}
	return agents, nil
}

// ═══════════════════════════════════════════
//  Task CRUD + Fee Escrow
// ═══════════════════════════════════════════

func (k Keeper) nextTaskID(ctx sdk.Context) (uint64, error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.NextTaskIDKey)
	if err != nil {
		return 0, err
	}
	if bz == nil {
		return 1, nil
	}
	return types.BytesToID(bz), nil
}

func (k Keeper) setNextTaskID(ctx sdk.Context, id uint64) error {
	store := k.storeService.OpenKVStore(ctx)
	return store.Set(types.NextTaskIDKey, types.IDToBytes(id))
}

func (k Keeper) CreateTask(ctx sdk.Context, msg types.MsgCreateTask) (string, error) {
	// 验证 Agent 存在且活跃
	agent, err := k.GetAgent(ctx, msg.AgentId)
	if err != nil {
		return "", err
	}
	if agent.Status != types.AgentStatus_AGENT_STATUS_ACTIVE {
		return "", types.ErrAgentNotActive
	}
	// 检查费用是否 >= Agent 要求的最低费用
	if !msg.Fee.IsAllGTE(agent.FeePerTask) {
		return "", types.ErrInsufficientFee
	}

	// 把费用从用户账户发到 module account（托管）
	creatorAddr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return "", types.ErrInvalidCreator
	}
	if msg.Fee.IsAllPositive() {
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, creatorAddr, types.ModuleName, msg.Fee); err != nil {
			return "", err
		}
	}

	nextID, err := k.nextTaskID(ctx)
	if err != nil {
		return "", err
	}
	taskID := fmt.Sprintf("task-%d", nextID)

	task := types.Task{
		Id:        taskID,
		AgentId:   msg.AgentId,
		Creator:   msg.Creator,
		InputHash: msg.InputHash,
		Status:    types.TaskStatus_TASK_STATUS_PENDING,
		Fee:       msg.Fee,
		CreatedAt: ctx.BlockTime().Unix(),
	}

	if err := k.setTask(ctx, task); err != nil {
		return "", err
	}
	// 索引
	store := k.storeService.OpenKVStore(ctx)
	if err := store.Set(types.GetTasksByAgentIndexKey(msg.AgentId, taskID), []byte(taskID)); err != nil {
		return "", err
	}
	if err := store.Set(types.GetTasksByCreatorIndexKey(msg.Creator, taskID), []byte(taskID)); err != nil {
		return "", err
	}
	if err := k.setNextTaskID(ctx, nextID+1); err != nil {
		return "", err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		"task_created",
		sdk.NewAttribute("task_id", taskID),
		sdk.NewAttribute("agent_id", msg.AgentId),
		sdk.NewAttribute("creator", msg.Creator),
		sdk.NewAttribute("fee", msg.Fee.String()),
	))

	return taskID, nil
}

func (k Keeper) SubmitTaskResult(ctx sdk.Context, msg types.MsgSubmitTaskResult) error {
	task, err := k.GetTask(ctx, msg.TaskId)
	if err != nil {
		return err
	}
	if task.Status != types.TaskStatus_TASK_STATUS_PENDING {
		return types.ErrInvalidTaskStatus
	}
	// 验证提交者是 agent 的 operator
	agent, err := k.GetAgent(ctx, task.AgentId)
	if err != nil {
		return err
	}
	if agent.Operator != msg.Operator {
		return types.ErrUnauthorized
	}

	// 更新 Task
	task.OutputHash = msg.OutputHash
	task.ResultUrl = msg.ResultUrl
	task.Status = types.TaskStatus_TASK_STATUS_COMPLETED
	task.CompletedAt = ctx.BlockTime().Unix()
	if err := k.setTask(ctx, task); err != nil {
		return err
	}

	// 释放资金：module → operator
	operatorAddr, err := sdk.AccAddressFromBech32(agent.Operator)
	if err != nil {
		return err
	}
	if task.Fee.IsAllPositive() {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, operatorAddr, task.Fee); err != nil {
			return err
		}
	}

	// 更新 Agent 统计
	agent.TasksCompleted++
	agent.ReputationScore += 10 // 每完成一个任务加 10 声望
	if err := k.setAgent(ctx, agent); err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		"task_completed",
		sdk.NewAttribute("task_id", task.Id),
		sdk.NewAttribute("agent_id", task.AgentId),
		sdk.NewAttribute("operator", agent.Operator),
	))

	return nil
}

func (k Keeper) CancelTask(ctx sdk.Context, creator, taskID string) error {
	task, err := k.GetTask(ctx, taskID)
	if err != nil {
		return err
	}
	if task.Creator != creator {
		return types.ErrUnauthorized
	}
	if task.Status != types.TaskStatus_TASK_STATUS_PENDING {
		return types.ErrInvalidTaskStatus
	}

	// 退还费用给 creator
	creatorAddr, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return err
	}
	if task.Fee.IsAllPositive() {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creatorAddr, task.Fee); err != nil {
			return err
		}
	}

	task.Status = types.TaskStatus_TASK_STATUS_CANCELLED
	if err := k.setTask(ctx, task); err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		"task_cancelled",
		sdk.NewAttribute("task_id", taskID),
		sdk.NewAttribute("creator", creator),
	))

	return nil
}

func (k Keeper) GetTask(ctx sdk.Context, id string) (types.Task, error) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetTaskKey(id)
	has, err := store.Has(key)
	if err != nil {
		return types.Task{}, err
	}
	if !has {
		return types.Task{}, types.ErrTaskNotFoundError{ID: id}
	}
	bz, err := store.Get(key)
	if err != nil {
		return types.Task{}, err
	}
	var task types.Task
	if err := k.cdc.Unmarshal(bz, &task); err != nil {
		return types.Task{}, err
	}
	return task, nil
}

func (k Keeper) setTask(ctx sdk.Context, task types.Task) error {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := k.cdc.Marshal(&task)
	if err != nil {
		return err
	}
	return store.Set(types.GetTaskKey(task.Id), bz)
}

func (k Keeper) GetTasks(ctx sdk.Context) ([]types.Task, error) {
	return iteratePrefix[types.Task](k, ctx, types.TaskKeyPrefix)
}

func (k Keeper) GetTasksByAgent(ctx sdk.Context, agentID string) ([]types.Task, error) {
	return k.getTasksByIndex(ctx, append(types.TasksByAgentIndexPrefix, []byte(agentID+"/")...))
}

func (k Keeper) GetTasksByCreator(ctx sdk.Context, creator string) ([]types.Task, error) {
	return k.getTasksByIndex(ctx, append(types.TasksByCreatorIndexPrefix, []byte(creator+"/")...))
}

func (k Keeper) getTasksByIndex(ctx sdk.Context, prefixBytes []byte) ([]types.Task, error) {
	prefixStore := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), prefixBytes)
	iter := prefixStore.Iterator(nil, nil)
	defer iter.Close()

	var tasks []types.Task
	for ; iter.Valid(); iter.Next() {
		taskID := string(iter.Value())
		task, err := k.GetTask(ctx, taskID)
		if err != nil {
			continue
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// ═══════════════════════════════════════════
//  Sequence helpers for genesis
// ═══════════════════════════════════════════

func (k Keeper) GetNextAgentID(ctx sdk.Context) (uint64, error)  { return k.nextAgentID(ctx) }
func (k Keeper) SetNextAgentID(ctx sdk.Context, id uint64) error { return k.setNextAgentID(ctx, id) }
func (k Keeper) GetNextTaskID(ctx sdk.Context) (uint64, error)   { return k.nextTaskID(ctx) }
func (k Keeper) SetNextTaskID(ctx sdk.Context, id uint64) error  { return k.setNextTaskID(ctx, id) }

// SetAgent and SetTask for genesis import
func (k Keeper) SetAgent(ctx sdk.Context, agent types.Agent) error {
	if err := k.setAgent(ctx, agent); err != nil {
		return err
	}
	store := k.storeService.OpenKVStore(ctx)
	return store.Set(types.GetAgentsByModelIndexKey(agent.ModelId, agent.Id), []byte(agent.Id))
}

func (k Keeper) SetTask(ctx sdk.Context, task types.Task) error {
	if err := k.setTask(ctx, task); err != nil {
		return err
	}
	store := k.storeService.OpenKVStore(ctx)
	if err := store.Set(types.GetTasksByAgentIndexKey(task.AgentId, task.Id), []byte(task.Id)); err != nil {
		return err
	}
	return store.Set(types.GetTasksByCreatorIndexKey(task.Creator, task.Id), []byte(task.Id))
}

// ═══════════════════════════════════════════
//  generic prefix iterator
// ═══════════════════════════════════════════

type protoMessage[T any] interface {
	*T
	codec.ProtoMarshaler
}

func iteratePrefix[T any, PT protoMessage[T]](k Keeper, ctx sdk.Context, pfx []byte) ([]T, error) {
	prefixStore := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), pfx)
	iter := prefixStore.Iterator(nil, nil)
	defer iter.Close()

	var items []T
	for ; iter.Valid(); iter.Next() {
		var item T
		if err := k.cdc.Unmarshal(iter.Value(), PT(&item)); err != nil {
			return items, err
		}
		items = append(items, item)
	}
	return items, nil
}

// GetAgentIDFromString parses "agent-123" to uint64 123 for sequence tracking.
func GetAgentIDFromString(id string) (uint64, error) {
	if len(id) <= 6 {
		return 0, fmt.Errorf("invalid agent id: %s", id)
	}
	return strconv.ParseUint(id[6:], 10, 64)
}

// GetTaskIDFromString parses "task-123" to uint64 123.
func GetTaskIDFromString(id string) (uint64, error) {
	if len(id) <= 5 {
		return 0, fmt.Errorf("invalid task id: %s", id)
	}
	return strconv.ParseUint(id[5:], 10, 64)
}

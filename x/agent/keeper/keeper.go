package keeper

import (
	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/store/prefix"
	"github.com/CosmWasm/wasmd/x/agent/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	cdc          codec.Codec
	storeService corestoretypes.KVStoreService
}

func NewKeeper(cdc codec.Codec, storeService corestoretypes.KVStoreService) Keeper {
	return Keeper{
		cdc:          cdc,
		storeService: storeService,
	}
}

// 注册模型

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
		return types.ErrModelAlreadyExists{ErrorContent: model.Id}
	}

	// 序列化 model 信息
	bz, err := k.cdc.Marshal(&model)
	if err != nil {
		return err
	}

	return store.Set(key, bz)
}

// 通过模型名称删除模型
func (k Keeper) DeleteModel(ctx sdk.Context, modelID string, modelCreator string) error {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetModelKey(modelID)

	has, err := store.Has(key)
	if err != nil {
		return err
	}
	if !has {
		return types.ErrModelNotFound{ErrorContent: modelID}
	}

	bz, err := store.Get(key)
	if err != nil {
		return err
	}
	var model types.Model
	if err := k.cdc.Unmarshal(bz, &model); err != nil {
		return err
	}

	if modelCreator != model.Creator {
		return types.ErrInvalidDeleteModel
	}

	return store.Delete(key)
}

// 通过模型名获取模型
func (k Keeper) GetModel(ctx sdk.Context, id string) (types.Model, error) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetModelKey(id)

	has, err := store.Has(key)
	if err != nil {
		return types.Model{}, err
	}
	if !has {
		return types.Model{}, types.ErrModelNotFound{ErrorContent: id}
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

// 获取当前所有模型
func (k Keeper) GetModels(ctx sdk.Context) ([]types.Model, error) {
	prefixStore := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.ModelKeyPrefix)
	iterator := prefixStore.Iterator(nil, nil)
	defer iterator.Close()
	var models []types.Model
	for ; iterator.Valid(); iterator.Next() {
		var model types.Model
		if err := k.cdc.Unmarshal(iterator.Value(), &model); err != nil {
			return models, err
		}
		models = append(models, model)
	}
	return models, nil
}

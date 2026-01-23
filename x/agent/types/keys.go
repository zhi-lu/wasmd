package types

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
)

func GetModelKey(name string) []byte {
	return append(ModelKeyPrefix, []byte(name)...)
}

package utils

import "github.com/bwmarrin/snowflake"

var (
	_snowflakeNode *snowflake.Node
	_snowflakeApp  *snowflake.Node
)

func init() {
	SetSnowflakeNodeId(1021)
}

func SetSnowflakeNodeId(nodeId int64) (err error) {
	_snowflakeNode, err = snowflake.NewNode(nodeId)
	return
}

func GetSnowflakeId() int64 {
	return _snowflakeNode.Generate().Int64()
}

func GetSnowflakeIdStr() string {
	return Int64ToStr(GetSnowflakeId())
}

func SetSnowflakeAppId(appId int64) (err error) {
	_snowflakeApp, err = snowflake.NewNode(appId)
	return
}

func GetSnowflakeAppId() int64 {
	return _snowflakeApp.Generate().Int64()
}

func Int64ToStr(id int64) string {
	return string([]byte{
		byte(id >> 56),
		byte(id >> 48),
		byte(id >> 40),
		byte(id >> 32),
		byte(id >> 24),
		byte(id >> 16),
		byte(id >> 8),
		byte(id),
	})
}

func StrToInt64(id string) int64 {
	bytes := []byte(id)
	return int64(bytes[0])<<56 |
		int64(bytes[1])<<48 |
		int64(bytes[2])<<40 |
		int64(bytes[3])<<32 |
		int64(bytes[4])<<24 |
		int64(bytes[5])<<16 |
		int64(bytes[6])<<8 |
		int64(bytes[7])
}

func Uint64ToStr(id uint64) string {
	return string([]byte{
		byte(id >> 56),
		byte(id >> 48),
		byte(id >> 40),
		byte(id >> 32),
		byte(id >> 24),
		byte(id >> 16),
		byte(id >> 8),
		byte(id),
	})
}

func StrToUint64(id string) uint64 {
	bytes := []byte(id)
	return uint64(bytes[0])<<56 |
		uint64(bytes[1])<<48 |
		uint64(bytes[2])<<40 |
		uint64(bytes[3])<<32 |
		uint64(bytes[4])<<24 |
		uint64(bytes[5])<<16 |
		uint64(bytes[6])<<8 |
		uint64(bytes[7])
}

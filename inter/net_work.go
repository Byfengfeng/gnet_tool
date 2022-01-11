package inter

type INetwork interface {
	Start()
	Action(action func())
	SetIsClose()
	GetCtx() interface{}
	WriteWriteChan(data []byte)
	GetAddr() string
	SetCid(cid int64)
	SetUid(uid int64)
	GetNetWorkBy(address string) INetwork
	CloseCid()
}

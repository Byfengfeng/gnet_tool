package inter

type INetwork interface {
	Start()
	Action(fn func())
	SetIsClose()
	WriteReadChan(data []byte)
	GetCtx() interface{}
	WriteWriteChan(data []byte)
	GetAddr() string
	SetCid(cid int64)
	SetUid(uid int64)
	GetNetWorkBy(address string) INetwork
	CloseCid()
}

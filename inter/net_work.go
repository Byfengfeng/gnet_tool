package inter

type INetwork interface {
	Start()
	GetClose() bool
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

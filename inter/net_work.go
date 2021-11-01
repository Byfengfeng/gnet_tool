package inter

type INetwork interface {
	Start()
	GetClose() bool
	SetIsClose()
	WriteReadChan(data []byte)
	GetCtx() interface{}
	WriteWriteChan(data []byte)
	DelNetWork()
	GetAddr() string
	GetNetWorkBy(address string) INetwork
	CloseCid()
}

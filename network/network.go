package network

import (
	"github.com/Byfengfeng/gnet_tool/code_tool"
	"github.com/Byfengfeng/gnet_tool/inter"
	"github.com/Byfengfeng/gnet_tool/log"
	"github.com/Byfengfeng/gnet_tool/utils"
	"github.com/panjf2000/gnet"
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
)

type NetWork struct {
	gnet.Conn
	ReadChan chan []byte
	WriteChan chan []byte
	IsClose bool
	CloseLock sync.Mutex
	WriteLock sync.Mutex
	Ctx *code_tool.IRequestCtx
}

var(
	netWorkMap = make(map[string]*NetWork)
	netWorkLock = sync.RWMutex{}
	count uint32
)

func NewNetWork(c gnet.Conn) inter.INetwork {
	address := c.RemoteAddr().String()
	netWorkLock.Lock()
	defer netWorkLock.Unlock()
	t,ok := netWorkMap[address]
	if ok {
		t.Conn.Close()
		t.Conn = c
	}else{
		t = &NetWork{c,
			make(chan[]byte),
			make(chan[]byte),
			false,
			sync.Mutex{},
			sync.Mutex{},
			code_tool.NewIRequestCtx(0,address),
		}
		netWorkMap[address] = t
	}
	return t
}

func (n *NetWork) read()  {
	for  {
		select {
		case reqBytes := <- n.ReadChan:
			if len(reqBytes) == 0 {
				log.Logger.Info("read off")
				return
			}
			//读取数据
			code, data := utils.Decode(reqBytes)
			log.Logger.Info("收到消息:",zap.Uint16("code:",code),zap.String("data:",string(data)))
			code_tool.Request(n.Ctx.Addr,n,code,data)
		}
	}
}

func (n *NetWork) write()  {
	for  {
		data := <- n.WriteChan
		if len(data) > 0 {
			err := n.Conn.AsyncWrite(data)
			if err != nil {
				log.Logger.Error("发送消息异常",zap.Any("err",err))
				return
			}
		}else{
			log.Logger.Info("write off")
			return
		}
	}
}

func (n *NetWork) Start()  {
	go n.read()
	go n.write()
}

func (n *NetWork) GetCtx() interface{} {
	return n.Ctx
}

func (n *NetWork) WriteReadChan(data []byte)  {
	n.ReadChan <- data
}

func (n *NetWork) WriteWriteChan(data []byte)  {
	n.WriteChan <- data
}

func (n *NetWork) SetIsClose()  {
	n.CloseLock.Lock()
	defer n.CloseLock.Unlock()
	n.IsClose = true
	n.DelNetWork()
}

func  (n *NetWork) CloseCid()  {
	code_tool.OffLine(n.Ctx.Addr)
}

func GetNetWork(address string) inter.INetwork {
	netWorkLock.Lock()
	defer netWorkLock.Unlock()
	netWork,ok := netWorkMap[address]
	if ok {
		return netWork
	}
	return nil
}

func (n *NetWork) GetNetWorkBy(address string) inter.INetwork {
	return GetNetWork(address)
}

func (n *NetWork) DelNetWork()  {
	netWorkLock.Lock()
	defer netWorkLock.Unlock()
	addr := n.RemoteAddr().String()
	n.Conn.Close()
	delete(netWorkMap,addr)
	close(n.ReadChan)
	close(n.WriteChan)
	log.Logger.Info("close network")
	atomic.AddUint32(&count,1)
}

func (n *NetWork) GetClose() bool {
	n.CloseLock.Lock()
	defer n.CloseLock.Unlock()
	return n.IsClose
}

func (n *NetWork) GetAddr() string {
	return n.RemoteAddr().String()
}

func GetCloseCount() uint32 {
	return count
}

func ResetCount()  {
	count = 0
}
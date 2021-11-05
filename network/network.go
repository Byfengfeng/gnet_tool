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
	Ctx *code_tool.IRequestCtx
	netV string
}

var(
	count uint32
)

func NewNetWork(c gnet.Conn,netV string) inter.INetwork {
	address := c.RemoteAddr().String()
	t := &NetWork{c,
		make(chan[]byte),
		make(chan[]byte),
		false,
		sync.Mutex{},
		code_tool.NewIRequestCtx(0,address),
		netV,
	}
	code_tool.NewChannel(t)
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
			log.Logger.Info("收到消息:"+string(data),zap.Uint16("code:",code))
			code_tool.Request(n.Ctx.Addr,n,code,data)
		}
	}
}

func (n *NetWork) write()  {
	for  {
		data := <- n.WriteChan
		if len(data) > 0 {
			var err error
			if n.netV == "tcp" || n.netV == "tcp4" || n.netV == "tcp6"{
				err = n.Conn.AsyncWrite(data)
			}else{
				err = n.Conn.SendTo(data)
			}

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
	n.Conn.Close()
	close(n.ReadChan)
	close(n.WriteChan)
	log.Logger.Info("close network")
	atomic.AddUint32(&count,1)
}

func  (n *NetWork) CloseCid()  {
	code_tool.OffLine(n.Ctx.Addr,n.Ctx.Cid)
}

func GetNetWork(address string) inter.INetwork {
	return code_tool.GetNetWorkByAddr(address)
}

func (n *NetWork) GetNetWorkBy(address string) inter.INetwork {
	return GetNetWork(address)
}

func (n *NetWork) GetClose() bool {
	n.CloseLock.Lock()
	defer n.CloseLock.Unlock()
	return n.IsClose
}

func (n *NetWork) GetAddr() string {
	return n.RemoteAddr().String()
}

func (n *NetWork) SetCid(cid int64)  {
	n.Ctx.Cid = cid
}

func (n *NetWork) SetUid(uid int64)  {
	n.Ctx.Uid = uid
}

func GetCloseCount() uint32 {
	return count
}

func SetCount() {
	count = 0
}


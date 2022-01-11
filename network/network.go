package network

import (
	"github.com/Byfengfeng/gnet_tool/code_tool"
	"github.com/Byfengfeng/gnet_tool/inter"
	"github.com/Byfengfeng/gnet_tool/log"
	"github.com/Byfengfeng/gnet_tool/utils"
	"go.uber.org/zap"
	"net"
	"sync"
)

type NetWork struct {
	*net.TCPConn
	WriteChan chan []byte
	IsClose   bool
	CloseLock sync.RWMutex
	Ctx       *code_tool.IRequestCtx
	ringByte  *utils.Bytes
}

func NewNetWork(c *net.TCPConn) {
	address := c.RemoteAddr().String()
	t := &NetWork{TCPConn:c,
		WriteChan: make(chan []byte),
		IsClose: true,
		CloseLock: sync.RWMutex{},
		Ctx: code_tool.NewIRequestCtx(0, address),
	}
	t.ringByte = utils.NewBytes(1024,func(bytes []byte) {
		code, data := utils.Decode(bytes)
		code_tool.Request(t.Ctx.Addr, t, code, data)
	})
	code_tool.NewChannel(t)
	t.Start()
}

func (n *NetWork) readBuff()  {
	for {
		newBytes := make([]byte, 1024)
		readLen, err := n.TCPConn.Read(newBytes)
		if err != nil {
			log.Logger.Info(err.Error())
			code_tool.OffLine(n.Ctx.Addr, n.Ctx.Cid)
			log.Logger.Info("readBuff off")
			return
		}
		if readLen == 0 {
			log.Logger.Info("readBuff off")
			code_tool.OffLine(n.Ctx.Addr, n.Ctx.Cid)
			return
		} else {
			n.ringByte.WriteBytes(uint16(readLen),newBytes[0:readLen-1])
		}
	}
}

func (n *NetWork) write() {
	for {
		if !n.IsClose {
			log.Logger.Info("write off")
			return
		}
		data := <-n.WriteChan
		if len(data) > 0 {
			_, err := n.Write(data)
			if err != nil {
				log.Logger.Error("发送消息异常", zap.Any("err", err))
				return
			}
		} else {
			log.Logger.Info("write off")
			return
		}
	}
}

func (n *NetWork) Start() {
	go n.readBuff()
	go n.write()
}

func (n *NetWork) GetCtx() interface{} {
	return n.Ctx
}

func (n *NetWork) WriteWriteChan(data []byte) {
	n.WriteChan <- data
}

func (n *NetWork) SetIsClose() {
	n.CloseLock.Lock()
	defer n.CloseLock.Unlock()
	if n.IsClose {
		n.IsClose = false
		n.TCPConn.Close()
		n.ringByte.Close()
		close(n.WriteChan)
		log.Logger.Info("close network")
	}
}

func (n *NetWork) CloseCid() {
	code_tool.OffLine(n.Ctx.Addr, n.Ctx.Cid)
}

func GetNetWork(address string) inter.INetwork {
	return code_tool.GetNetWorkByAddr(address)
}

func (n *NetWork) GetNetWorkBy(address string) inter.INetwork {
	return GetNetWork(address)
}

func (n *NetWork) Action(action func()) {
	n.CloseLock.RLock()
	if n.IsClose {
		action()
	}
	n.CloseLock.RUnlock()
}

func (n *NetWork) GetAddr() string {
	return n.RemoteAddr().String()
}

func (n *NetWork) SetCid(cid int64) {
	n.Ctx.Cid = cid
}

func (n *NetWork) SetUid(uid int64) {
	n.Ctx.Uid = uid
}

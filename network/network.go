package network

import (
	"fmt"
	"github.com/Byfengfeng/gnet_tool/inter"
	"github.com/Byfengfeng/gnet_tool/utils"
	"github.com/panjf2000/gnet"
	"sync"
)

type NetWork struct {
	gnet.Conn
	ReadChan chan []byte
	WriteChan chan []byte
	CloseChan chan bool
}

var(
	netWorkMap = make(map[string]*NetWork)
	netWorkLock = sync.RWMutex{}
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
		t = &NetWork{c,make(chan[]byte),make(chan[]byte),make(chan bool)}
		netWorkMap[address] = t
	}
	return t
}

func (n *NetWork) read()  {
	defer func() {
		close(n.ReadChan)
		close(n.WriteChan)
		close(n.CloseChan)
	}()
	for  {
		select {
		case reqBytes := <- n.ReadChan:
				//读取数据
				code, data := utils.Decode(reqBytes)
				fmt.Println(fmt.Sprintf("code: %d,",code),"data:",string(data))
		case <- n.CloseChan:
			return
		}
	}
}

func (n *NetWork) write()  {
	for  {
		data := <- n.WriteChan
		if len(data) > 0 {
			err := n.Conn.AsyncWrite(data)
			if err != nil {
				fmt.Println("发送消息异常")
				return
			}
		}
	}
}

func (n *NetWork) Start()  {
	go n.read()
	go n.write()
}

func GetNetWork(address string) *NetWork {
	netWork,ok := netWorkMap[address]
	if ok {
		return netWork
	}
	return nil
}

func DelNetWork(addr string)  {
	netWorkLock.Lock()
	netWork,ok := netWorkMap[addr]
	if ok {
		netWork.Conn.Close()
		delete(netWorkMap,addr)
		netWork.CloseChan <- true
	}
	defer netWorkLock.Unlock()
}
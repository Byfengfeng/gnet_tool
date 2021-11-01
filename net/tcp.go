package net

import (
	"fmt"
	"github.com/Byfengfeng/gnet_tool/code_tool"
	"github.com/Byfengfeng/gnet_tool/inter"
	"github.com/Byfengfeng/gnet_tool/log"
	"github.com/Byfengfeng/gnet_tool/network"
	"github.com/Byfengfeng/gnet_tool/utils"
	"github.com/panjf2000/ants/v2"
	"github.com/panjf2000/gnet"
	"go.uber.org/zap"
	"runtime/debug"
	"time"
)

type tcpServer struct {
	*gnet.EventServer
	tcpVersion string
	addr	   string
	ip		   uint16
	multicore  bool
	network.NetWork
}

func (t *tcpServer) NewEventHandler() gnet.EventHandler {
	return 	t
}

func NewTcpServer(tcpVersion,addr string,ip uint16,multicore bool) *tcpServer {
	return 	&tcpServer{tcpVersion: tcpVersion,addr: addr,ip: ip,multicore: multicore}
}

func (t *tcpServer) OnInitComplete(server gnet.Server) (action gnet.Action)  {
	log.Logger.Info("tcp listen start")
	return
}

func (t *tcpServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
		ants.Submit(func() {
			if len(frame) > 0 {
				copyByte := make([]byte,len(frame))
				copy(copyByte,frame)
				if len(copyByte) > 0 {
					netWork := network.GetNetWork(c.RemoteAddr().String())
					if netWork!= nil {
						if netWork != nil {
							codeDe(copyByte,netWork)
						}
					}
				}
			}
		})
		return
}

func codeDe(frame []byte,network inter.INetwork) {
	data, remainingByte := utils.DecodeRound(frame)
	if !network.GetClose() {
		network.WriteReadChan(data)
		if len(remainingByte) > 0 {
			codeDe(remainingByte,network)
		}
	}
}

func (t *tcpServer) Tick() (delay time.Duration, action gnet.Action) {
	return
}

func (t *tcpServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action)  {
	network.NewNetWork(c).Start()
	return
}

func (t *tcpServer) OnShutdown(svr gnet.Server) {
	return
}

func (t *tcpServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	netWork := network.GetNetWork(c.RemoteAddr().String())
	if netWork != nil && !netWork.GetClose() {
		netWork.CloseCid()
	}
	return gnet.Close
}

func (t *tcpServer) Start() (err error) {
	go freeOs()
	options := make([]gnet.Option,0)
	if t.multicore {
		options = append(options,gnet.WithMulticore(t.multicore))
	}
	options = append(options,gnet.WithCodec(code_tool.NewICodec()))
	options = append(options,gnet.WithNumEventLoop(200))
	options = append(options,gnet.WithLogger(log.DefaultLogger()))

	err = gnet.Serve(t.NewEventHandler(), fmt.Sprintf("%s://%s:%d",t.tcpVersion,t.addr,t.ip),
		options...)
	return
}

func freeOs() {
	t := 1
	if network.GetCloseCount() > 0 {
		t = 2
	}
	timer := time.NewTicker(time.Duration(t)*time.Minute)
	select {
	case <-timer.C:
		log.Logger.Info("开始执行内存释放：",zap.Int("time：",t))
		debug.FreeOSMemory()
		freeOs()
	}

}
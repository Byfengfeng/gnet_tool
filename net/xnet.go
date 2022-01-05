package net

import (
	"fmt"
	"github.com/Byfengfeng/gnet_tool/code_tool"
	"github.com/Byfengfeng/gnet_tool/log"
	"github.com/Byfengfeng/gnet_tool/net/tcp"
	"github.com/Byfengfeng/gnet_tool/net/udp"
	"github.com/Byfengfeng/gnet_tool/network"
	"github.com/panjf2000/ants/v2"
	"github.com/panjf2000/gnet"
	"go.uber.org/zap"
	"runtime/debug"
	"time"
)

type netServer struct {
	*gnet.EventServer
	netVersion string
	addr	   string
	ip		   uint16
	multicore  bool
	network.NetWork
}

func (t *netServer) NewEventHandler() gnet.EventHandler {
	return 	t
}

func NewTcpServer(netVersion,addr string,ip uint16,multicore bool) *netServer {
	return 	&netServer{netVersion: netVersion,addr: addr,ip: ip,multicore: multicore}
}

func (t *netServer) OnInitComplete(server gnet.Server) (action gnet.Action)  {
	log.Logger.Info(fmt.Sprintf("%s listen start",t.netVersion))
	return
}

func (t *netServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	ants.Submit(func() {
		if t.netVersion == "tcp" || t.netVersion == "tcp4" || t.netVersion == "tcp6"{
			tcp.TcpReact(frame,c)
		}else{
			udp.UdpReact(frame,c,t.netVersion)
		}
	})
	return
}

func (t *netServer) Tick() (delay time.Duration, action gnet.Action) {
	return
}

func (t *netServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action)  {
	network.NewNetWork(c,t.netVersion).Start()
	return
}

func (t *netServer) OnShutdown(svr gnet.Server) {
	return
}

func (t *netServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	netWork := network.GetNetWork(c.RemoteAddr().String())
	if netWork != nil {
		netWork.Action(func() {
			netWork.CloseCid()
		})

	}

	return gnet.Close
}

func (t *netServer) Start() (err error) {
	options := make([]gnet.Option,0)
	if t.multicore {
		options = append(options,gnet.WithMulticore(t.multicore))
	}
	if t.netVersion == "tcp" || t.netVersion == "tcp4" || t.netVersion == "tcp6"{
		options = append(options,gnet.WithCodec(code_tool.NewICodec()))
	}

	options = append(options,gnet.WithNumEventLoop(2000))
	options = append(options,gnet.WithLogger(log.DefaultLogger()))
	err = gnet.Serve(t.NewEventHandler(), fmt.Sprintf("%s://%s:%d",t.netVersion,t.addr,t.ip),
		options...)
	return
}

func freeOs() {
	if network.GetCloseCount() > 100 {
		log.Logger.Info("开始执行内存释放：",zap.Uint32("time：",network.GetCloseCount()))
		debug.FreeOSMemory()
		network.SetCount()
		freeOs()
	}
}

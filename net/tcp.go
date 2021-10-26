package net

import (
	"fmt"
	"gitee.com/byfengfeng/gnet_tool/code_tool"
	"github.com/panjf2000/ants/v2"
	"github.com/panjf2000/gnet"
	"time"
)

type tcpServer struct {
	*gnet.EventServer
	tcpVersion	string
	addr	   string
	ip		   uint16
	multicore  bool
	async      bool
	asyncFunc func(frame []byte, c gnet.Conn)
	noAsyncFunc func(frame []byte, c gnet.Conn) []byte
}

func (t *tcpServer) NewEventHandler() gnet.EventHandler {
	return 	t
}

func NewTcpServer(tcpVersion,addr string,ip uint16,multicore,async bool,
	asyncFunc func(frame []byte, c gnet.Conn),
	noAsyncFunc func(frame []byte, c gnet.Conn) []byte) *tcpServer {
	if async {
		return 	&tcpServer{tcpVersion: tcpVersion,addr: addr,ip: ip,multicore: multicore,async: async,asyncFunc: asyncFunc}
	}
	return 	&tcpServer{tcpVersion: tcpVersion,addr: addr,ip: ip,multicore: multicore,async: async,noAsyncFunc: noAsyncFunc}
}

func (t *tcpServer) OnInitComplete(server gnet.Server) (action gnet.Action)  {
	fmt.Println("tcp listen start")
	return
}

func (t *tcpServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	defer c.ResetBuffer()
	if t.async {
		ants.Submit(func() {
			if len(frame) >= 0 {
				t.asyncFunc(frame,c)
			}
		})
		return
	}
	out = t.noAsyncFunc(frame, c)
	return
}

func (t *tcpServer) Tick() (delay time.Duration, action gnet.Action) {
	return
}

func (t *tcpServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action)  {
	return
}

func (t *tcpServer) OnShutdown(svr gnet.Server) {
}

func (t *tcpServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	return gnet.Close
}

func (t *tcpServer) Start() (err error) {
	options := make([]gnet.Option,0)
	options = append(options,gnet.WithCodec(code_tool.NewICodec()))
	if t.multicore {
		options = append(options,gnet.WithMulticore(t.multicore))
	}
	if t.async {
		options = append(options,gnet.WithCodec(code_tool.NewICodec()))
	}
	options = append(options,gnet.WithNumEventLoop(200))
	err = gnet.Serve(t.NewEventHandler(), fmt.Sprintf("%s://%s:%d",t.tcpVersion,t.addr,t.ip),
		options...)
	return
}
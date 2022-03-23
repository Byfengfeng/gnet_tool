package tcp

import (
	"github.com/Byfengfeng/gnet_tool/log"
	"github.com/Byfengfeng/gnet_tool/network"
	"go.uber.org/zap"
	"net"
)

type tcpListen struct {
	address string
	*net.TCPListener
	channelHandel func(conn *net.TCPConn)
}
 
func NewTcpListen(addr string) *tcpListen {
	return &tcpListen{address: addr,channelHandel: network.NewNetWorkTcp}
}

func (n * tcpListen) Start() error {
	addr, err := net.ResolveTCPAddr("tcp", n.address)
	if err != nil {
		return err
	}
	n.TCPListener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Logger.Error("listener err",zap.Any("err",err))
			}
		}()
		for {
			tcpConn, err := n.TCPListener.AcceptTCP()
			if err != nil {
				log.Logger.Error("client channel exit",zap.Any("err",err))
			}
			go n.channelHandel(tcpConn)
		}
	}()
	log.Logger.Info("tcp listen start success")
	return nil
}

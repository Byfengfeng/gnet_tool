package net

import (
	"github.com/Byfengfeng/gnet_tool/log"
	"github.com/Byfengfeng/gnet_tool/network"
	"go.uber.org/zap"
	"net"
)

type udpListen struct {
	address string
	*net.UDPConn
	channelHandel func(conn *net.TCPConn)
}

func NewUcpListen(addr string) *udpListen {
	return &udpListen{address: addr,channelHandel: network.NewNetWorkTcp}
}

func (n *udpListen) Start() error {
	addr, err := net.ResolveUDPAddr("udp", n.address)
	if err != nil {
		return err
	}
	n.UDPConn, err = net.ListenUDP("udp", addr)
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
			tcpConn, err := n.UDPConn.ReadFromUDP()
			if err != nil {
				log.Logger.Error("client channel exit",zap.Any("err",err))
			}
			go n.channelHandel(tcpConn)
		}
	}()
	log.Logger.Info("tcp listen start success")
	return nil
}

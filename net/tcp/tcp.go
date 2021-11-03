package tcp

import (
	"github.com/Byfengfeng/gnet_tool/inter"
	"github.com/Byfengfeng/gnet_tool/network"
	"github.com/Byfengfeng/gnet_tool/utils"
	"github.com/panjf2000/gnet"
)

func TcpReact(frame []byte, c gnet.Conn)  {
	if len(frame) > 0 {
		copyByte := make([]byte,len(frame))
		copy(copyByte,frame)
		if len(copyByte) > 0 {
			netWork := network.GetNetWork(c.RemoteAddr().String())
			if netWork!= nil {
				TcpCodeDe(copyByte,netWork)
			}
		}
	}
}

func TcpCodeDe(frame []byte,network inter.INetwork) {
	data, remainingByte := utils.DecodeRound(frame)
	if !network.GetClose() {
		network.WriteReadChan(data)
		if len(remainingByte) > 0 {
			TcpCodeDe(remainingByte,network)
		}
	}
}





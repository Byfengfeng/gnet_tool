package udp

import (
	"github.com/Byfengfeng/gnet_tool/network"
	"github.com/panjf2000/gnet"
)

func UdpReact(frame []byte,c gnet.Conn){
	if len(frame) > 0{
		copyByte := make([]byte,len(frame))
		copy(copyByte,frame)
		if len(copyByte) > 0 {
			netWork := network.GetNetWork(c.RemoteAddr().String())
			if netWork!= nil {
				if !netWork.GetClose() {
					netWork.WriteReadChan(frame)
				}
			}else{
				netWork = network.NewNetWork(c)
				netWork.Start()
				if !netWork.GetClose() {
					netWork.WriteReadChan(frame)
				}
			}
		}
	}
}
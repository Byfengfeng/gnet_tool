package main

import (
	"fmt"
	"github.com/Byfengfeng/gnet_tool/net/tcp"
	"github.com/Byfengfeng/gnet_tool/network"
	"github.com/Byfengfeng/gnet_tool/utils"
)

func main() {
	tcp.NewNetListen(":8999",network.NewNetWork())
	<-make(chan struct{})

	//TestPool()
}

func TestPool()  {

	pool := utils.NewSafePool(func() interface{} {
		return []byte{}
	})
	bytes := pool.Get().([]byte)
	bytes = []byte("123str")
	pool.Put(bytes)
	str := pool.Get().([]byte)
	fmt.Println(string(str))

}
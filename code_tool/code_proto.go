package code_tool

import (
	"fmt"
	"github.com/Byfengfeng/gnet_tool/log"
	"github.com/Byfengfeng/gnet_tool/utils"
	"github.com/gogo/protobuf/proto"
)

func NewCodecProto() *CodecBase {
	return &CodecBase{
		reqPool: make(map[uint16]utils.IPool),
		resPool: make(map[uint16]utils.IPool),
	}
}

type CodecBase struct {
	reqPool map[uint16]utils.IPool
	resPool map[uint16]utils.IPool
}

func (c *CodecBase) BindPool(code uint16, requestNewFunc func() interface{}, responseNewFunc func() interface{}) {
	if requestNewFunc != nil {
		_,ok := c.reqPool[code]
		if ok {
			log.Logger.Error(fmt.Sprintf("%d req code 绑定实体出现重复",code))
		}
		c.reqPool[code] = utils.NewSafePool(requestNewFunc)
	}
	if responseNewFunc != nil {
		_,ok := c.resPool[code]
		if ok {
			log.Logger.Error(fmt.Sprintf("%d res code 绑定实体出现重复",code))
		}
		c.resPool[code] = utils.NewSafePool(responseNewFunc)
	}
}

func (c *CodecBase) GetReqPkt(code uint16) interface{} {
	pool,ok := c.reqPool[code]
	if !ok {
		log.Logger.Error(fmt.Sprintf("%d not req pkt",code))
		return nil
	}
	return pool.Get().(proto.Message)
}

func (c *CodecBase) GetResPkt(code uint16) interface{} {
	pool,ok := c.resPool[code]
	if !ok {
		log.Logger.Error(fmt.Sprintf("%d not res pkt",code))
	}
	return pool.Get()
}

func (c *CodecBase) PutReqPkt (code uint16,pkt interface{})  {
	pool,ok := c.reqPool[code]
	if !ok {
		log.Logger.Error(fmt.Sprintf("%d not res pkt",code))
	}
	pool.Put(pkt)
}

func (c *CodecBase) PutResPkt (code uint16,pkt interface{})  {
	pool,ok := c.resPool[code]
	if !ok {
		log.Logger.Error(fmt.Sprintf("%d not res pkt",code))
	}
	pool.Put(pkt)
}

func (c *CodecBase) DecodeReq(code uint16, bytes []byte) (interface{}, error) {
	pool,ok := c.reqPool[code]
	if !ok {
		log.Logger.Error(fmt.Sprintf("%d not res pkt",code))
	}
	pkt := pool.Get().(proto.Message)
	err := proto.Unmarshal(bytes, pkt)
	if err != nil {
		return nil, err
	}
	return pkt, nil
}

func (c *CodecBase) DecodeRes(code uint16, bytes []byte) (interface{}, error) {
	packet := c.GetResPkt(code)
	err := proto.Unmarshal(bytes, packet.(proto.Message))
	if err != nil {
		return nil, err
	}
	return packet, nil
}

func (c *CodecBase) EncodeRes(code uint16,set func(interface{})) []byte{
	resPkt := c.GetResPkt(code).(proto.Message)
	set(resPkt)
	bytes, err := proto.Marshal(resPkt)
	if err != nil {
		log.Logger.Error("pb序列化失败")
	}
	return utils.Encode(code,bytes)
}
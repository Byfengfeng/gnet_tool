package channel_manager

import (
	"ldmp_game/pkt/manager/channel/channel"
	"sync"
)

type ChannelManager struct {
	channelMap map[int64]*channel.Channel
	channelMtx sync.RWMutex
}

func NewChannelManager() *ChannelManager {
	return &ChannelManager{
		channelMap: make(map[int64]*channel.Channel),
	}
}

//添加连接
func (c *ChannelManager) AddChannel(channel *channel.Channel,uid int64)  {
	c.channelMtx.Lock()
	c.channelMap[uid] = channel
	c.channelMtx.Unlock()
	channel.Start()

}

//删除连接
func (c *ChannelManager) DelChannel(channel *channel.Channel,uid int64) error {
	c.channelMtx.RLock()
	channel,ok := c.channelMap[uid]
	if ok {
		err := channel.Conn.Close()
		if err != nil {
			return err
		}
		delete(c.channelMap,uid)
	}
	c.channelMtx.RUnlock()
	return nil
}

func (c *ChannelManager) GetChannelManager() map[int64]*channel.Channel {
	return c.channelMap
}
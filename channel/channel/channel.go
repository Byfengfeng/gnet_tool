package channel

import (
	"errors"
	"github.com/Byfengfeng/gnet_tool/utils"
	"github.com/panjf2000/gnet"
)

type Channel struct {
	Conn gnet.Conn
	ReadChan chan []byte
	WriteChan chan []byte
}

func NewChannel(conn gnet.Conn) *Channel {
	return &Channel{
		Conn: conn,
		ReadChan: make(chan []byte),
		WriteChan: make(chan []byte),
	}
}

func (c *Channel) write() (err error) {
	defer c.Close()
	for  {
		bytes := <- c.WriteChan
		if len(bytes) > 0 {
			err = c.Conn.AsyncWrite(bytes)
			if err != nil {
				return err
			}
		}
	}
}

func (c *Channel) read() (err error) {
	defer c.Close()
	var(
		headByte []byte
		length uint32
	)
	for  {
		headByte = make([]byte,4)
		//_, err = c.Conn.Read(headByte)
		if err != nil {
			return err
		}
		length = utils.Length(headByte)
		if length <= 0{
			return errors.New("错误消息长度")
		}
		body := make([]byte,length-4)
		//_, err = c.Conn.Read(body)
		if err != nil {
			return err
		}
		c.ReadChan <- body
	}
}

func (c *Channel) Start()  {
	go c.read()
	go c.write()
}

func (c *Channel) Close()  {
	if c.ReadChan == nil{
		close(c.ReadChan)
		close(c.WriteChan)
	}

	c.Conn.Close()
}
package code_tool

import "github.com/panjf2000/gnet"

type ICodec struct {
	gnet.ICodec
}

func NewICodec() gnet.ICodec {
	return 	&ICodec{}
}

func (i *ICodec) Decode(c gnet.Conn) ([]byte, error) {
	size, headByte := c.ReadN(c.BufferLength())
	if size > 0{
		c.ResetBuffer()
		return headByte,nil
	}
	return nil,nil
}

func (i *ICodec) Encode(c gnet.Conn, buf []byte) ([]byte, error) {
	return append([]byte{byte(len(buf) >> 8),byte(len(buf))},buf...),nil
}

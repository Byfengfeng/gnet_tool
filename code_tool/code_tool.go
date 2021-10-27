package code_tool

import "github.com/panjf2000/gnet"

type ICodec struct {
	gnet.ICodec
}

func NewICodec() gnet.ICodec {
	return 	&ICodec{}
}

func (i *ICodec) Decode(c gnet.Conn) ([]byte, error) {
	size, headByte := c.ReadN(2)
	if size == 0{
		return nil,nil
	}
	length := uint16(headByte[0]) << 8 | uint16(headByte[1])
	if length > 0 {
		size, dataBytes := c.ReadN(int(length-2))
		if size == 0{
			return nil,nil
		}
		return dataBytes,nil
	}

	return nil,nil
}

func (i *ICodec) Encode(c gnet.Conn, buf []byte) ([]byte, error) {
	return append([]byte{byte(len(buf) >> 8),byte(len(buf))},buf...),nil
}

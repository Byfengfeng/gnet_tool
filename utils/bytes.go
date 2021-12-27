package utils

type bytes struct {
	len         uint16 // byte total length
	outsetPos   uint16 // byte current use outset position
	endPos      uint16 // byte use end position
	ringByte    []byte
	initByteCap uint16
}

func NewBytes(initByteCap uint16) *bytes {
	return &bytes{initByteCap,0,initByteCap,make([]byte,initByteCap),initByteCap}
}

func (b *bytes) UseBytes(useLen uint16,putByte ...byte)  {
	        e c
	1 2 3 4 5 6
	if b.outsetPos == b.endPos || b.len - b.checkByteUseLen() < useLen{

	}
	if useLen > b.len b.endPos {

	}
}

func (b *bytes) checkByteUseLen() uint16 {
	useLen := b.outsetPos - b.endPos
	if useLen < 0 {
		useLen = -useLen
	}
	return useLen
}
package utils

import "sync"

type Bytes struct {
	len         uint16 // byte total length
	readPos     uint16 // byte read position
	writePos    uint16 // byte write position
	initByteCap uint16 // byte init cap
	beUsable    uint16 // byte be usable size
	dropCount	uint8  // byte drop size 10 byte len is change initByteCap
	ringByte    []byte // byte use model
	byteOperate	sync.RWMutex
}

func NewBytes(initByteCap uint16) *Bytes {
	return &Bytes{initByteCap, 0, 0, initByteCap, initByteCap, 0,make([]byte, initByteCap),sync.RWMutex{}}
}

func (b *Bytes) WriteBytes(useLen uint16, putByte ...byte) {
	if b.beUsable  >= useLen{
		b.byteOperate.RLock()
		if b.len-1 - b.writePos >= useLen {
			copy(b.ringByte[b.writePos:],putByte)
			b.beUsable -= useLen
		}else{
			oneWriteSize := b.len - b.writePos
			//two := useLen - oneWriteSize
			copy(b.ringByte[b.writePos:],putByte[:oneWriteSize])
			copy(b.ringByte[0:],putByte[oneWriteSize:])
			b.beUsable -= useLen
		}
		b.byteOperate.RUnlock()
		if b.len > b.initByteCap && b.beUsable > b.initByteCap {
			b.dropCount++
			if b.dropCount >= 10 {
				// drop byte size
				b.dropByt()
			}
		}
	}else{
		//add byte size
		b.addByt()
		b.byteOperate.RLock()
		defer b.byteOperate.RUnlock()
		copy(b.ringByte[b.writePos:],putByte)
		b.beUsable -= useLen
	}
}

func (b *Bytes) addByt()  {
	b.byteOperate.Lock()
	b.len += b.initByteCap
	b.beUsable += b.initByteCap
	b.readPos += b.initByteCap
	b.ringByte = append(append(b.ringByte[:b.writePos],make([]byte,b.initByteCap)...),b.ringByte[b.writePos:]...)
	b.byteOperate.Unlock()
}

func (b *Bytes) dropByt()  {
	b.byteOperate.Lock()
	b.byteOperate.Unlock()

}

func (b *Bytes) ReadBytes() []byte {
	b.byteOperate.RLock()
	b.byteOperate.RUnlock()
	if b.len > b.beUsable {
		if b.len - b.readPos >= 2 {
			//byteSize := b.ringByte[b.readPos : b.readPos+2]
			if b.len - b.readPos == 2{

			}
			b.readPos+=2
			//useByesSize := uint32(byteSize[0]) << 8 | uint32(byteSize[1])
		}


	}
	return []byte{}
}
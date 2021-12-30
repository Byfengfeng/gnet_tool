package utils

import (
	"sync"
	"sync/atomic"
)

type Bytes struct {
	len         uint16 // byte total length
	readPos     uint16 // byte read position
	writePos    uint16 // byte write position
	initByteCap uint16 // byte init cap
	beUsable    uint16 // byte be usable size
	dropCount	uint8  // byte drop size 10 byte len is change initByteCap
	ringByte    []byte // byte use model
	byteOperate	sync.RWMutex
	writeChan   chan<- []byte // byte use model
}

func NewBytes(writeChan chan <-[]byte,initByteCap uint16) *Bytes {
	return &Bytes{initByteCap, 0, 0, initByteCap, initByteCap, 0,make([]byte, initByteCap),
		sync.RWMutex{},writeChan}
}

//todo 扩容和缩容的时候不能读写，读写的时候需要考虑是否需要拼接

func (b *Bytes) WriteBytes(useLen uint16, putByte []byte) {
	if b.beUsable  >= useLen{
		b.byteOperate.RLock()
		if b.len-1 - b.writePos >= useLen {
			copy(b.ringByte[b.writePos:],putByte)
			atomic.AddInt32(&b.writePos,useLen)
			b.beUsable -= useLen
			b.writePos += useLen
		}else{
			oneWriteSize := b.len - b.writePos
			//two := useLen - oneWriteSize
			copy(b.ringByte[b.writePos:],putByte[:oneWriteSize])
			copy(b.ringByte[0:],putByte[oneWriteSize:])
			b.beUsable -= useLen
			b.writePos = useLen - oneWriteSize
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
		b.WriteBytes(useLen,putByte)
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

func (b *Bytes) ReadBytes() {

	if b.len != b.beUsable {
		b.byteOperate.RLock()

		if b.len - b.readPos > 2 {

			useByesSize := uint16(b.ringByte[b.readPos]) << 8 | uint16(b.ringByte[b.readPos+1])
			if useByesSize > b.beUsable {
				return
			}
			if b.len - b.readPos == 2{
				b.writeChan <- []byte("123")
			}
			b.readPos+=2

		}else{

		}
		b.byteOperate.RUnlock()
	}
}
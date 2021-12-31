package utils

import (
	"sync"
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


//read ringBuff byte data
//one . by ringBuff len - readPosition = one read len
//two by read len - one read len check read
//
//读取 环形缓冲区中的字节数据
//1.通过缓冲区空间总长度 - 读取位置 = 正序可读长度
//2.通过需要读取的字节长度判断正序可读长度能否支持一次读取
//3.如果不够一次读取则进行拼接
//
//todo 可能存在问题当环形缓冲区长度经过扩容或缩容的时候将改变，这会影响环形缓冲区读取的判断，目前是通过读写锁来控制，写入和读取缓冲区使用的是读锁，
//环形缓冲区扩容和缩容时使用写锁
//当有写锁的时候不能使用读锁，读锁可以有多个

func (b *Bytes) ReadBytes() ([]byte,error){

	if b.len > b.beUsable {
		b.byteOperate.RLock()

		if b.len - b.readPos > 2 {
			useByesSize := uint16(b.ringByte[b.readPos]) << 8 | uint16(b.ringByte[b.readPos+1])
			if useByesSize > b.beUsable {
				return []byte{},nil
			}
			b.readPos+=2
			if b.len - b.readPos == 2{
				if b.len-1 - b.readPos > useByesSize {
					b.writeChan <- b.ringByte[b.readPos:b.readPos+useByesSize]
					b.readPos += useByesSize
				}else if b.len-1 - b.readPos == useByesSize{
					b.writeChan <- b.ringByte[b.readPos:b.readPos+useByesSize]
					b.readPos = 0
				}else{
					readSize := b.len-1 - b.readPos
					readSize1 := useByesSize - readSize
					b.writeChan <- append(b.ringByte[b.readPos:b.readPos+readSize],b.ringByte[0:readSize1]...)
					b.readPos = readSize1
				}
				b.beUsable += useByesSize
			}else {

			}

		}else{

		}
		b.byteOperate.RUnlock()
	}
}

func (b *Bytes) CheckNeedSplice()  {
	
}
package utils

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
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
	writeChan   chan[]byte // byte use model
	writeLen int64
	checkClose bool
}

func NewBytes(initByteCap uint16) *Bytes {
	return &Bytes{initByteCap, 0, 0, initByteCap, initByteCap, 0,make([]byte, initByteCap),
		sync.RWMutex{},make(chan []byte,0),0,false}
}

//todo 扩容和缩容的时候不能读写，读写的时候需要考虑是否需要拼接

func (b *Bytes) WriteBytes(useLen uint16, putByte []byte) {
	if b.checkClose {
		return
	}
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
	atomic.AddInt64(&b.writeLen,1)
}

func (b *Bytes) addByt()  {
	b.byteOperate.Lock()
	b.len += b.initByteCap
	b.beUsable += b.initByteCap
	b.ringByte = append(append(b.ringByte[:b.writePos],make([]byte,b.initByteCap)...),b.ringByte[b.writePos:]...)
	fmt.Println("扩容长度",b.len)
	b.byteOperate.Unlock()
	time.Sleep(1 * time.Nanosecond)
}

func (b *Bytes) dropByt()  {
	b.byteOperate.Lock()
	wL := b.len-1 - b.writePos
	if wL > b.initByteCap {
		b.ringByte = append(b.ringByte[0:b.writePos],b.ringByte[b.writePos+b.initByteCap:]...)
	}
	b.len-=b.initByteCap
	fmt.Println("縮容长度",b.len)
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

func (b *Bytes) Close()  {
	b.checkClose = true
}

func (b *Bytes) ReadBytes() *Bytes {
	go func() {
		for  {
			if b.checkClose {
				b.byteOperate.Lock()
				close(b.writeChan)
				b.byteOperate.RUnlock()
				return
			}
			if b.writeLen > 0 {
				dateLength := b.ReadN(2)
				useByesSize := uint16(dateLength[0]) << 8 | uint16(dateLength[1])
				if useByesSize > b.len {
					// pain err byte len
					panic("解析字节包长度异常")
					return
				}
				data := b.ReadN(useByesSize)
				if len(data) > 0 {
					b.writeChan <- data
				}
				atomic.AddInt64(&b.writeLen,-1)
			}
		}
	}()
	return b
}

func (b *Bytes) Read() chan[]byte {
	return b.writeChan
}

func (b *Bytes) ReadN(byteLen uint16) []byte {
	b.byteOperate.RLock()
	defer b.byteOperate.RUnlock()
	readLen := b.len-1 - b.readPos
	if readLen > byteLen {
		bytes := b.ringByte[b.readPos : b.readPos+byteLen]
		b.readPos += byteLen
		b.beUsable += byteLen
		return bytes
	} else if readLen < byteLen {
		readSize1 := byteLen-(readLen+1)
		bytes := append(b.ringByte[b.readPos:])
		if readSize1 > 0 {
			bytes = append(bytes,b.ringByte[0:readSize1]...)
		}
		b.beUsable += byteLen
		b.readPos = readSize1
		return  bytes
	} else {
		bytes := b.ringByte[b.readPos:b.readPos+byteLen]
		b.beUsable += byteLen
		b.readPos = b.readPos+byteLen
		return bytes
	}
}
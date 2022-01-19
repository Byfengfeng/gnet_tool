package utils

import (
	"fmt"
	"github.com/Byfengfeng/gnet_tool/log"
	"sync"
	"sync/atomic"
	"time"
)

type Bytes struct {
	len           uint16 // byte total length
	readPos       uint16 // byte read position
	writePos      uint16 // byte write position
	initByteCap   uint16 // byte init cap
	beUsable      uint16 // byte be usable size
	recoverUsable uint16 // byte be usable size
	dropCount     uint8  // byte drop size 10 byte len is change initByteCap
	ringByte      []byte // byte use model
	byteOperate   sync.RWMutex
	writeLen      int64
	checkClose    bool
	messageFn     func(bytes []byte)
	readL         sync.Mutex
}

func NewBytes(initByteCap uint16, fn func(bytes []byte)) *Bytes {
	return &Bytes{initByteCap, 0, 0, 1024, initByteCap,0, 0, make([]byte, initByteCap),
		sync.RWMutex{}, 0, false, fn, sync.Mutex{}}
}
func (b *Bytes) Len()  {
	fmt.Println(b.len)
}

func (b *Bytes) WriteBytes(useLen uint16, putByte []byte) {
	if b.checkClose {
		return
	}
	//b.writeN(useLen,putByte)
	if b.beUsable >= useLen {
		b.byteOperate.RLock()
		if b.len-b.writePos >= useLen {
			copy(b.ringByte[b.writePos:], putByte)
			b.beUsable -= useLen
			b.writePos += useLen
		} else {
			oneWriteSize := b.len - b.writePos
			//two := useLen - oneWriteSize
			copy(b.ringByte[b.writePos:], putByte[:oneWriteSize])
			copy(b.ringByte[0:], putByte[oneWriteSize:])
			b.beUsable -= useLen
			b.writePos = useLen - oneWriteSize
		}
		b.byteOperate.RUnlock()
	} else {
		//add byte size
		b.addByt(useLen)
		b.WriteBytes(useLen, putByte)
		return
	}
	atomic.AddInt64(&b.writeLen, 1)
	go b.ReadBytes()
}

func (b *Bytes) addByt(needLen uint16) {
	b.byteOperate.Lock()
	if b.beUsable >= needLen {
		return
	}
	b.len += b.initByteCap
	b.beUsable += b.initByteCap
	if b.readPos > b.writePos {
		b.readPos += b.initByteCap
	}
	bytes := append(b.ringByte[:b.writePos], make([]byte, b.initByteCap)...)
	b.ringByte = append(bytes, b.ringByte[b.writePos:]...)
	fmt.Println("扩容长度", b.len)
	b.byteOperate.Unlock()
}

func (b *Bytes) dropByt() {
	b.byteOperate.Lock()
	wL := b.len - b.writePos
	if b.beUsable > b.initByteCap && wL > b.initByteCap && b.readPos <= b.writePos{
		b.ringByte = append(b.ringByte[:b.writePos],b.ringByte[b.writePos+b.initByteCap:]...)
		b.len -= b.initByteCap
		b.beUsable -= b.initByteCap
		b.dropCount = 0
		fmt.Println("縮容长度", b.len)
	}

	b.byteOperate.Unlock()

}

//read ringBuff byte data
//one . by ringBuff len - readPosition = one read len
//two by read len - one read len check read

func (b *Bytes) Close() {
	b.byteOperate.Lock()
	if !b.checkClose {
		b.checkClose = true
	}
	b.byteOperate.Unlock()
}

func (b *Bytes) ReadBytes() {
	if b.checkClose {
		log.Logger.Info("ReadBytes off")
		return
	}
	if b.writeLen > 0 {
		time.Sleep(3 * time.Second)
		b.readL.Lock()
		dateLength := b.readN(2)
		useByesSize := uint16(dateLength[0])<<8 | uint16(dateLength[1])
		if useByesSize > b.len {
			// pain err byte len
			panic("解析字节包长度异常")
			return
		}
		data := b.readN(useByesSize)
		atomic.AddInt64(&b.writeLen, -1)
		b.readL.Unlock()
		if len(data) > 0 {
			go b.messageFn(data)
		}
	}

}

func (b *Bytes) readN(byteLen uint16) []byte {
	b.byteOperate.RLock()
	defer b.byteOperate.RUnlock()
	readLen := b.len - b.readPos
	if readLen > byteLen {
		bytes := b.ringByte[b.readPos : b.readPos+byteLen]
		b.readPos += byteLen
		b.beUsable += byteLen
		return bytes
	} else if readLen < byteLen {
		readSize1 := byteLen - readLen
		bytes := append(b.ringByte[b.readPos:])
		if readSize1 > 0 {
			bytes = append(bytes, b.ringByte[0:readSize1]...)
		}
		b.beUsable += byteLen
		b.readPos = readSize1
		return bytes
	} else {
		bytes := b.ringByte[b.readPos : b.readPos+byteLen]
		b.beUsable += byteLen
		b.readPos = b.readPos + byteLen
		return bytes
	}
}

func (b *Bytes) writeN(writeLen uint16, bytes []byte) {
	indexOutOfBounds := b.CheckIndexOutOfBounds(writeLen)
	if b.beUsable > writeLen {
		b.byteOperate.RLock()
		if indexOutOfBounds {
			w1 := b.beUsable - b.writePos
			copy(b.ringByte[b.writePos:], bytes[:w1+1])
			copy(b.ringByte[0:], bytes[w1+1:])
			b.writePos = writeLen - w1 - 1
			b.beUsable -= writeLen
		} else {
			copy(b.ringByte[b.writePos:], bytes)
			b.beUsable -= writeLen
			b.writePos += writeLen
		}
		b.byteOperate.RUnlock()
	} else if b.beUsable < writeLen {
		b.addByt(writeLen)
		b.writeN(writeLen, bytes)
	} else {
		b.byteOperate.RLock()
		if indexOutOfBounds {
			w1 := b.beUsable - b.writePos
			copy(b.ringByte[b.writePos:], bytes[:w1+1])
			copy(b.ringByte[0:], bytes[w1+1:])
			b.writePos = writeLen - w1 - 1
			b.beUsable -= writeLen
		} else {
			copy(b.ringByte[b.writePos:], bytes)
			b.beUsable -= writeLen
			b.writePos = 0
		}
		b.byteOperate.RUnlock()
	}
}
func (b *Bytes) CheckIndexOutOfBounds(len uint16) bool {
	useWriteLen := b.len - 1 - b.writePos
	if useWriteLen >= len {
		return false
	}
	return true
}

func Ad()  {

}
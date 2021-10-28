package utils

import (
	"github.com/gopherjs/gopherjs/nosync"
	"sync"
)

type IPool interface {
	Get() interface{}
	Put(x interface{})
}

func NewSafePool(f func()interface{}) *sync.Pool{
	return &sync.Pool{New: f}
}

func NewPool(f func()interface{}) *nosync.Pool{
	return &nosync.Pool{New: f}
}
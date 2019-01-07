package ctx

import (
	"reflect"
	"sync"
)

type CtxKey string

func (k CtxKey) String() string {
	return string(k)
}

func NewCtx(parent Ctx) Ctx {
	c := &defaultCtx{
		parent: parent,
	}
	c.attributes = make(map[interface{}]interface{})
	return c
}

type Ctx interface {
	Parent() Ctx
	SetAttribute(key interface{}, value interface{})
	GetAttribute(key interface{}) (value interface{})
	RemoveAttribute(key interface{})
	ContainsAttribute(key interface{}) (exist bool)
}

type defaultCtx struct {
	parent     Ctx
	attributes map[interface{}]interface{}

	mtx sync.RWMutex
}

func (dc *defaultCtx) Parent() Ctx {
	return dc.parent
}

func (dc *defaultCtx) SetAttribute(key interface{}, value interface{}) {
	dc.checkInitialized()

	if key == nil {
		panic("nil key")
	}
	if !reflect.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}

	dc.mtx.Lock()
	defer dc.mtx.Unlock()

	dc.attributes[key] = value
}

func (dc *defaultCtx) GetAttribute(key interface{}) (value interface{}) {
	dc.checkInitialized()

	dc.mtx.RLock()
	defer dc.mtx.RUnlock()

	if _, ok := dc.attributes[key]; ok {
		return dc.attributes[key]
	}

	if nil == dc.parent {
		return nil
	}
	return dc.parent.GetAttribute(key)
}

func (dc *defaultCtx) RemoveAttribute(key interface{}) {
	dc.checkInitialized()

	dc.mtx.Lock()
	defer dc.mtx.Unlock()

	if _, ok := dc.attributes[key]; ok {
		delete(dc.attributes, key)
		return
	}

	if nil == dc.parent {
		return
	}

	dc.parent.RemoveAttribute(key)
}

func (dc *defaultCtx) ContainsAttribute(key interface{}) (exist bool) {
	dc.checkInitialized()

	dc.mtx.RLock()
	defer dc.mtx.RUnlock()

	if _, ok := dc.attributes[key]; ok {
		return true
	}

	if nil == dc.parent {
		return false
	}
	return dc.parent.ContainsAttribute(key)
}

func (dc *defaultCtx) checkInitialized() {
	if nil == dc.attributes {
		panic("Attribute Manager: must be initialized")
	}
}

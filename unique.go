package unique

import (
	"container/list"
)

type KeyType = interface{}
type ValueType = interface{}

var emptyKey KeyType = nil
var emptyValue ValueType = nil

type valueItem struct {
	e *list.Element
	v ValueType
}

type FuncMerge func(oldVal, newVal ValueType) ValueType

func RetainOld(oldVal, newVal ValueType) ValueType {
	return oldVal
}

func RetainNew(oldVal, newVal ValueType) ValueType {
	return newVal
}

type Unique = *unique

// Unique unique queue
type unique struct {
	data    map[KeyType]*valueItem
	queue   *list.List
	fnMerge FuncMerge
}

// New create a unique queue.
func New(fnMerge FuncMerge) Unique {
	return &unique{
		data:    make(map[KeyType]*valueItem),
		queue:   list.New(),
		fnMerge: fnMerge,
	}
}

// Push push a key-value into queue. return true if the key is not exists
func (u Unique) Push(key KeyType, value ValueType) (isNew bool) {
	if item, ok := u.data[key]; ok {
		item.v = u.fnMerge(item.v, value)
		return false
	}

	e := u.queue.PushBack(key)
	u.data[key] = &valueItem{
		e: e,
		v: value,
	}
	return true
}

// Pop pop a item from the queue
func (u Unique) Pop() (KeyType, ValueType, bool) {
	el := u.queue.Front()
	if el == nil {
		return emptyKey, emptyValue, false
	}
	key := u.queue.Remove(el)
	if item, ok := u.data[key]; ok {
		delete(u.data, key)
		return key, item.v, true
	}
	return key, emptyValue, false
}

// Get Get a value by key
func (u Unique) Get(key KeyType) (ValueType, bool) {
	if item, ok := u.data[key]; ok {
		return item.v, true
	}
	return emptyValue, false
}

// Del Delete a element by key and return the deleted value
func (u Unique) Del(key KeyType) (ValueType, bool) {
	if item, ok := u.data[key]; ok {
		delete(u.data, key)
		v := u.queue.Remove(item.e)
		return v, true
	}
	return emptyValue, false
}

// Len Unique queue length
func (u Unique) Len() int {
	return u.queue.Len()
}

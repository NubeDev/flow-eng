package uuid

import "sync/atomic"

type Value uint64

var value uint64

func New() Value {
	return Value(atomic.AddUint64(&value, 1))
}

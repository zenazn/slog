package slog

import (
	"sync/atomic"
	"unsafe"
)

// Let's keep all the unsafe atomic hijinks in a locked room by itself.

func (l *logger) atomicGetLCache() *levelCache {
	ptr := (*unsafe.Pointer)(unsafe.Pointer(&l.lcache))
	return (*levelCache)(atomic.LoadPointer(ptr))
}

func (l *logger) atomicSetLCache(lc *levelCache) {
	ptr := (*unsafe.Pointer)(unsafe.Pointer(&l.lcache))
	atomic.StorePointer(ptr, unsafe.Pointer(lc))
}

func (l *logger) atomicGetTCache() *targetCache {
	ptr := (*unsafe.Pointer)(unsafe.Pointer(&l.tcache))
	return (*targetCache)(atomic.LoadPointer(ptr))
}

func (l *logger) atomicSetTCache(tc *targetCache) {
	ptr := (*unsafe.Pointer)(unsafe.Pointer(&l.tcache))
	atomic.StorePointer(ptr, unsafe.Pointer(tc))
}

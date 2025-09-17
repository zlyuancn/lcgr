package cnum

import (
	"sync/atomic"
)

type Random interface {
	// 返回一个随机数. 调用 2^64-1 次返回的值不会重复
	Uint64() uint64

	// 获取当前已调用次数
	GetCount() uint64
}

type RandomN interface {
	// 返回一个随机数, 调用 n 次返回的值不会重复, 区间在 [0, n) 之间. 必须使 0 < n < 1e10
	Uint64N() uint64

	// 获取当前已调用次数
	GetCount() uint64
}

type randomCli struct {
	i    uint64
	seed uint64
	n    uint64
}

func NewRandom(seed uint64) Random {
	return &randomCli{seed: seed}
}

// 同 NewRandom 并假设已调用了 count 次
func NewRandomSetCount(seed, count uint64) Random {
	return &randomCli{i: count, seed: seed}
}

func NewRandomN(seed, n uint64) RandomN {
	return &randomCli{seed: seed, n: n}
}

// 同 NewRandomN 并假设已调用了 count 次
func NewRandomNSetCount(seed, n, count uint64) RandomN {
	return &randomCli{seed: seed, n: n, i: count}
}

func (r *randomCli) Uint64() uint64 {
	i := atomic.AddUint64(&r.i, 1)
	return Confuse(i, r.seed)
}

func (r *randomCli) Uint64N() uint64 {
	i := atomic.AddUint64(&r.i, 1)
	return ConfuseLimit(i, r.seed, r.n)
}

func (r *randomCli) GetCount() uint64 {
	return atomic.LoadUint64(&r.i)
}

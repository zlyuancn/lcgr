package cnum

import (
	"sync/atomic"
)

type Random interface {
	// 返回一个随机数. 调用 2^64-1 次返回的值不会重复. 不会返回0
	Next() uint64

	// 获取下一个用于计算的sn
	GetNextSn() uint64
}

type randomCli struct {
	sn    uint64
	seed  uint64
	limit uint64
}

// 随机数生成器
func NewRandom(seed uint64) Random {
	return &randomCli{seed: seed}
}

// 同 NewRandom 并设置起始 sn
func NewRandomStartSn(seed, sn uint64) Random {
	return &randomCli{sn: sn, seed: seed}
}

func (r *randomCli) Next() uint64 {
	v := uint64(0)
	for v == 0 {
		sn := atomic.AddUint64(&r.sn, 1) - 1
		v = Confuse(sn, r.seed)
	}
	return v
}

func (r *randomCli) GetNextSn() uint64 {
	return atomic.LoadUint64(&r.sn)
}

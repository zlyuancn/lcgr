package cnum

import (
	"sync/atomic"
)

type Random interface {
	// 返回一个随机数. 调用 2^64-1 次返回的值不会重复. 不会返回0
	Next() uint64

	// 获取当前sn
	GetSn() uint64
}

type RandomN interface {
	// 返回一个随机数, 调用 n 次返回的值不会重复, 区间在 [0, limit) 之间. 要求调用次数必须小于 limit, limit 不能超过 1e10
	NextLimit() uint64

	// 获取当前sn
	GetSn() uint64
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

// 限制范围的随机数生成器. limit 不能超过 1e10
func NewRandomLimit(seed, limit uint64) RandomN {
	return &randomCli{seed: seed, limit: limit}
}

// 同 NewRandomN 并设置起始 sn. limit 不能超过 1e10
func NewRandomLimitStartSn(seed, limit, sn uint64) RandomN {
	return &randomCli{seed: seed, limit: limit, sn: sn}
}

func (r *randomCli) Next() uint64 {
	v := uint64(0)
	for v == 0 {
		i := atomic.AddUint64(&r.sn, 1) - 1
		v = Confuse(i, r.seed)
	}
	return v
}

func (r *randomCli) NextLimit() uint64 {
	i := atomic.AddUint64(&r.sn, 1) - 1
	return ConfuseLimit(i, r.seed, r.limit)
}

func (r *randomCli) GetSn() uint64 {
	return atomic.LoadUint64(&r.sn)
}

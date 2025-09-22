package lcgr

import (
	"sync/atomic"
)

type Random interface {
	// 返回一个随机数. 调用 1e18 次返回的值不会重复, 返回值区间在 [0, 1e18), 包含 0 不包含 1e18
	Next() uint64

	// 获取下一个用于计算的sn
	GetNextSn() uint64
}

type RandomLimitLen interface {
	// 返回一个随机数, 返回数值长度不会达到limit
	NextLimit() uint64

	// 获取下一个用于计算的sn
	GetNextSn() uint64
}

type randomCli struct {
	sn       uint64
	seed     uint64
	limitLen uint64
}

// 随机数生成器, 返回值区间在 [0, 1e18), 包含 0 不包含 1e18
func NewRandom(seed uint64) Random {
	return &randomCli{seed: seed}
}

// 同 NewRandom 并设置起始 sn
func NewRandomStartSn(seed, sn uint64) Random {
	return &randomCli{sn: sn, seed: seed}
}

// 随机数生成器, 限制返回值的长度. limitLen 范围必须是 [1, 18], 包含 1 和 18
func NewRandomLimitLen(seed, limitLen uint64) RandomLimitLen {
	return &randomCli{seed: seed, limitLen: limitLen}
}

// 同 NewRandomLimitLen 并设置起始 sn
func NewRandomLimitLenStartSn(seed, limitLen, sn uint64) RandomLimitLen {
	return &randomCli{seed: seed, limitLen: limitLen, sn: sn}
}

func (r *randomCli) Next() uint64 {
	sn := atomic.AddUint64(&r.sn, 1) - 1
	return Confuse(sn, r.seed)
}

func (r *randomCli) NextLimit() uint64 {
	sn := atomic.AddUint64(&r.sn, 1) - 1
	return ConfuseLimitLen(sn, r.seed, r.limitLen)
}

func (r *randomCli) GetNextSn() uint64 {
	return atomic.LoadUint64(&r.sn)
}

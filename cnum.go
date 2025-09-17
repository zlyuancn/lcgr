package cnum

import (
	"math"
)

// 质数乘积偏移(值, 偏移量)
func confuse(v, offset uint64) uint64 {
	const a uint64 = 0x5bd1e995 // 一个质数
	return v*a + offset
}

/*
对有序编号混淆, 不同的编号得到结果必然不同

	index 编号
	seed 随机种子
	block 块. 必须在 1 ~ 9 之间
*/
func confuseBlock(index, seed uint64, block int) uint64 {
	// if block < 1 {
	// 	panic("block must be greater than 0")
	// }
	// if block > 9 {
	// 	panic("block must be greater than 9")
	// }
	mod := uint64(math.Pow10(block))

	offset := seed % mod // 偏移值不能太大

	bi := index / mod // 块编号
	i := index % mod  // 块值

	bio := confuse(bi, offset) // 块随机偏移
	return confuse(i, bio%mod)%mod + bio*mod
}

/*
对有序编号混淆, 不同的编号得到结果必然不同
*/
func Confuse(index, seed uint64) uint64 {
	v := index
	for i := 9; i >= 1; i-- {
		v = confuseBlock(v, seed, i)
	}
	return v
}

/*
对有序编号混淆一次, 不同的编号得到结果必然不同, 并将混淆值限制在一定范围内. 要求编号的总数不能超过limit大小

	index 编号
	seed 随机种子
	limit 限制结果值在区间 [0, limit). 包含0, 不包含limit
*/
func ConfuseLimit(index, seed, limit uint64) uint64 {
	if limit < 1 {
		panic("block must be greater than 0")
	}
	if limit > 1e10 {
		panic("block must be less than 1e10")
	}
	mod := limit
	const a uint64 = 0x5bd1e995 // 一个质数
	v := ((index+seed)%mod)*a + seed
	return v % mod
}

/*
对有序编号混淆p次, 不同的编号得到结果必然不同, 并将混淆值限制在一定范围内. 要求编号的总数不能超过limit大小

	index 编号
	seed 随机种子
	limit 限制结果值在区间 [0, limit). 包含0, 不包含limit
*/
func ConfuseLimitP(index, seed, limit uint64, p int) uint64 {
	v := index
	for i := 0; i < p; i++ {
		v = ConfuseLimit(v, seed, limit)
	}
	return v
}

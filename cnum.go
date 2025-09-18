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
对有序序号混淆, 不同的序号得到结果必然不同

	sn 序号
	seed 随机种子
	block 块. 必须在 1 ~ 9 之间
*/
func confuseBlock(sn, seed uint64, block int) uint64 {
	// if block < 1 {
	// 	panic("block must be greater than 0")
	// }
	// if block > 9 {
	// 	panic("block must be greater than 9")
	// }
	mod := uint64(math.Pow10(block))

	offset := seed % mod // 偏移值不能太大

	bi := sn / mod // 块序号
	i := sn % mod  // 块值

	bio := confuse(bi, offset) // 块随机偏移
	return confuse(i, bio%mod)%mod + bio*mod
}

/*
对有序序号混淆, 不同的序号得到结果必然不同
*/
func Confuse(sn, seed uint64) uint64 {
	v := sn
	for i := 9; i >= 1; i-- {
		v = confuseBlock(v, seed, i)
	}
	return v
}

/*
对有序序号混淆一次, 不同的序号得到结果必然不同, 并将混淆值限制在一定范围内. 要求序号必须小于 limit, limit 不能超过 1e10

	sn 序号
	seed 随机种子
	limit 限制结果值在区间 [0, limit). 包含0, 不包含limit
*/
func ConfuseLimit(sn, seed, limit uint64) uint64 {
	if limit < 1 {
		panic("block must be greater than 0")
	}
	if limit > 1e10 {
		panic("block must be less than 1e10")
	}
	if sn >= limit {
		panic("sn out of limit")
	}
	mod := limit
	const a uint64 = 0x5bd1e995 // 一个质数
	v := ((sn+seed)%mod)*a + seed
	return v % mod
}

/*
对有序序号混淆p次, 不同的序号得到结果必然不同, 并将混淆值限制在一定范围内. 要求序号的总数不能超过limit大小

	index 序号
	seed 随机种子
	limit 限制结果值在区间 [0, limit). 包含0, 不包含limit
*/
func ConfuseLimitP(sn, seed, limit uint64, p int) uint64 {
	v := sn
	for i := 0; i < p; i++ {
		v = ConfuseLimit(v, seed, limit)
	}
	return v
}

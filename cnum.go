package cnum

import (
	"fmt"
)

const primeNumber1e10 uint64 = 9999999967 // 一百亿以内中最大的质数

var blockArgs = []struct {
	limit              uint64 // 限制值
	modH, modL         uint64 // 高位和低位具有不同的mod
	sh1, sh2, sl1, sl2 uint64 // 高位和低位分割后再次分割
}{
	{1e1, 1e1, 0, 0, 0, 0, 0},
	{1e2, 1e2, 0, 1e1, 1e1, 0, 0},
	{1e3, 1e3, 0, 1e2, 1e1, 0, 0},
	{1e4, 1e4, 0, 1e3, 1e1, 0, 0},
	{1e5, 1e5, 0, 1e3, 1e2, 0, 0},
	{1e6, 1e6, 0, 1e4, 1e2, 0, 0},
	{1e7, 1e7, 0, 1e4, 1e3, 0, 0},
	{1e8, 1e8, 0, 1e5, 1e3, 0, 0},
	{1e9, 1e9, 0, 1e5, 1e4, 0, 0},

	{1e10, 1e5, 1e5, 1e4, 1e1, 1e3, 1e2},
	{1e11, 1e5, 1e6, 1e4, 1e1, 1e4, 1e2},
	{1e12, 1e6, 1e6, 1e5, 1e1, 1e4, 1e2},
	{1e13, 1e6, 1e7, 1e5, 1e1, 1e4, 1e3},
	{1e14, 1e7, 1e7, 1e5, 1e2, 1e4, 1e3},
	{1e15, 1e7, 1e8, 1e5, 1e2, 1e5, 1e3},
	{1e16, 1e8, 1e8, 1e6, 1e2, 1e5, 1e3},
	{1e17, 1e8, 1e9, 1e6, 1e2, 1e5, 1e4},
	{1e18, 1e9, 1e9, 1e6, 1e3, 1e5, 1e4},
}

// 质数乘积偏移(值, 偏移量)
func confuse(v, offset uint64) uint64 {
	return v*primeNumber1e10 + offset
}

/*
对有序序号混淆, 不同的序号得到结果必然不同, sn 的长度必须小于 limitLen. limitLen 必须在 [10, 18] 之间, 包含 10 和 18

	sn 序号, 长度必须小于 limitLen
	seed 随机种子
	limitLen 限制sn长度和最终结果长度, 范围必须是 [10, 18], 包含 10 和 18

将sn分为高位和低位分别混淆. 并将高低位切换后再次混淆.
*/
func confuseLimit(sn, seed, limitLen uint64) uint64 {
	if limitLen < 10 {
		panic("limitLen must be >= 10")
	}
	if limitLen > 18 {
		panic("limitLen must be <= 18")
	}

	args := blockArgs[limitLen-1]

	if sn >= args.limit {
		panic(fmt.Sprintf("sn out of %d limitLen", limitLen))
	}

	offset := _mod(confuse(seed, 0), args.modL) // 偏移值混淆

	for i := uint64(0); i < 5; i++ { // ps: 经过实际验证, 重复5次混淆低位值会出现视觉上的随机性
		snH := sn / args.modL      // 高位
		snL := _mod(sn, args.modL) // 低位

		// (snH+i)%mod or (snL+i)%mod 使用 i 对混淆sn值做基础偏移并限制范围在1e9内, 不同的sn必然得到不同的值.
		vH := _mod(confuse(_mod(snH+i, args.modH), offset), args.modH) // 高位混淆. ps: 相同的高位其结果必然相同
		vL := _mod(confuse(_mod(snL+i, args.modL), vH), args.modL)     // 低位混淆, 并使用高位的混淆结果作为偏移量. ps: 高位不变的情况下, vH 值必然相同, 所以可以使用 vH 作为偏移量
		// 到这里已经完成了一次混淆. 不同的 sn 必然映射为不同的 vH 和 vL, sn结果为 vL + vH*mod, 不同的sn必然得到不同的值.

		// 数值高低位交换, 由于高低位长度不同, 经过多次循环和混淆后会将数值变化曲线蔓延到每一位值
		h := vH / args.sh1
		vH = _mod(vH, args.sh1)*args.sh2 + h

		l := vL / args.sl1
		vL = _mod(vL, args.sl1)*args.sl2 + l

		// 这里对其结果切换高低位, 在for的下一个循环中再次进行混淆
		sn = vL*args.modH + vH
	}
	return sn
}

func _mod(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return a % b
}

/*
对有序序号混淆, 不同的序号得到结果必然不同, sn 必须小于 1e18

	sn 序号, 必须小于 1e18
	seed 随机种子
*/
func Confuse(sn, seed uint64) uint64 {
	return confuseLimit(sn, seed, 18)
}

/*
对有序序号混淆, 不同的序号得到结果必然不同, sn 的长度必须小于 limitLen. limitLen 必须在 [1, 18] 之间, 包含 1 和 18

	sn 序号, 长度必须小于 limitLen
	seed 随机种子
	limitLen 限制sn长度和最终结果长度, 范围必须是 [1, 18], 包含 1 和 18

将sn分为高位和低位分别混淆. 并将高低位切换后再次混淆.
*/
func ConfuseLimitLen(sn, seed, limitLen uint64) uint64 {
	if limitLen < 1 {
		panic("limitLen must be >= 1")
	}
	if limitLen > 9 {
		return confuseLimit(sn, seed, limitLen)
	}

	args := blockArgs[limitLen-1]

	if sn >= args.limit {
		panic(fmt.Sprintf("sn out of %d limitLen", limitLen))
	}

	if limitLen == 1 {
		for i := 0; i < 5; i++ {
			sn = _mod(confuse(sn, seed), args.modH)
		}
		return sn
	}

	offset := _mod(confuse(seed, 0), args.modH) // 偏移值混淆

	for i := uint64(0); i < 5; i++ {
		v := _mod(confuse(_mod((sn+i), args.modH), offset), args.modH)

		// 数值高低位交换
		h := v / args.sh1
		sn = _mod(v, args.sh1)*args.sh2 + h
	}
	return sn
}

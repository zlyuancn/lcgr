package cnum

const primeNumber1e10 uint64 = 9999999967 // 一百亿以内中最大的质数

// 质数乘积偏移(值, 偏移量)
func confuse(v, offset uint64) uint64 {
	return v*primeNumber1e10 + offset
}

/*
对有序序号混淆, 不同的序号得到结果必然不同, 必须小于 1e18

	sn 序号, 必须小于 1e18
	seed 随机种子

将sn分为高位和低位分别混淆. 并将高低位切换后再次混淆.
*/
func Confuse(sn, seed uint64) uint64 {
	if sn >= 1e18 {
		panic("sn out of 1e18 limit")
	}

	const mod uint64 = 1e9
	offset := confuse(seed, 0) % mod // 偏移值混淆

	snH := sn / mod    // 高位
	snL := sn % mod    // 低位
	vH, vL := snH, snL // 混淆结果

	for i := uint64(0); i < 5; i++ { // ps: 在不交换的情况下, 经过实际验证, 重复5次混淆低位值会出现视觉上的随机性
		// (snH+i)%mod or (snL+i)%mod 使用 i 对混淆sn值做基础偏移并限制范围在1e9内, 不同的sn必然得到不同的值.
		vH = confuse((snH+i)%mod, offset) % mod // 高位混淆. ps: 相同的高位其结果必然相同
		vL = confuse((snL+i)%mod, vH) % mod     // 低位混淆, 并使用高位的混淆结果作为偏移量. ps: 高位不变的情况下, vH 值必然相同, 所以可以使用 vH 作为偏移量
		// 到这里已经完成了一次混淆. 不同的 sn 必然映射为不同的 vH 和 vL, sn结果为 vL + vH*mod, 不同的sn必然得到不同的值.

		// 数值高低位交换, 由于高低位长度不同, 经过多次循环和混淆后会将数值变化曲线蔓延到每一位值
		h := vH / 1e5
		vH = vH%1e5*1e4 + h
		l := vL / 1e6
		vL = vL%1e6*1e3 + l

		// 这里对其结果再次进行高低位分割, 并切换高低位, 在for的下一个循环中再次进行混淆
		snH, snL = vL, vH // 高低交换并定义为sn
	}
	return vL + vH*mod
}

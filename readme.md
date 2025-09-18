# 序号混淆

用于将序号混淆为另一个值, 不同的序号混淆结果必然不同.
提供了一个随机数生成器, 在规定范围内生成的随机数不会重复.

# 示例场景

## id生成. 2^64 次生成结果不会重复, 支持多节点分布式生成

通过数值自增, 生成一个有序序号. 调用函数 `Confuse` 将序号映射为一个新的值, 在`seed`不变的前提下, 不同的序号映射的结果必然不同.

```go
const seed = 0xabcdef12 // seed 是一个随机数种子, 可以是任意值, 业务自行保管
sn := IncrBy() // IncrBy 是一个序号生成器, 用于数值自增. 可以考虑使用 redis 实现以支持分布式节点生成
id := cnum.Confuse(sn, seed) // 使用 seed 混淆序号 
```

## 对生成的数值范围做限制

通过数值自增, 生成一个有序序号. 使用 `ConfuseLimit` 可以对生成的数值范围做限制. 要求序号的总数必须小于 `limit` 值. 比如预估id最多有 100w 个. 则序号最大为 100w-1

```go
const seed = 0xabcdef12 // seed 是一个随机数种子, 可以是任意值, 业务自行保管
sn := IncrBy() // IncrBy 是一个序号生成器, 用于数值自增. 可以考虑使用 redis 实现以支持分布式节点生成
id := cnum.ConfuseLimit(sn, seed, 1e6) // 使用 seed 混淆序号 
```

## 随机数生成, 生成 2^64 次结果不重复

在`seed`不变的前提下, 生成的随机数顺序是一样的.

```go
const seed = 0xabcdef12 // seed 是一个随机数种子, 可以是任意值, 业务自行保管
r := cnum.NewRandom(seed) // 创建一个随机数生成器
r.Next() // 生成随机数
```

## 限制随机数范围

要求调用次数必须小于 limit, limit 不能超过 1e10

```go
const seed = 0xabcdef12 // seed 是一个随机数种子, 可以是任意值, 业务自行保管
r := cnum.NewRandom(seed) // 创建一个随机数生成器
r.NextLimit() // 生成随机数
```

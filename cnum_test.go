package cnum

import (
	"sync"
	"sync/atomic"
	"testing"
)

const testCount = 1e7
const testSeed = 0xabcdef12 // 一个seed

// 计算基础理论是否正确
func TestBasic(t *testing.T) {
	return // 需要验证时注释掉这一行

	const testCount = 1e9
	const mod uint64 = 1e9

	const testSeed = 0xabcdef12 // 一个seed

	const offset = testSeed % mod

	// 计算重复性
	const bfCount int = 10000            // bf数量
	const bfVCount = testCount / bfCount // 每个bf存放多少数据
	bfs := [bfCount]map[uint64]struct{}{}
	outCh := [bfCount]chan uint64{}
	for i := 0; i < bfCount; i++ {
		bfs[i] = make(map[uint64]struct{}, testCount/bfCount)
		outCh[i] = make(chan uint64, 100)
	}

	// 汇总结果
	wgOut := new(sync.WaitGroup)
	wgOut.Add(bfCount)
	for i := 0; i < bfCount; i++ {
		go func(i int) {
			bf := bfs[i]
			ch := outCh[i]
			for v := range ch {
				bf[v] = struct{}{}
			}
			wgOut.Done()
		}(i)
	}

	// 将计算分为多个协程
	var p int32 // 进度
	const threadCount = bfCount
	wgTh := new(sync.WaitGroup)
	wgTh.Add(bfCount)
	for ti := 0; ti < threadCount; ti++ {
		block := bfVCount // 每个线程计算多少数据
		go func(ti int) {
			iOffset := ti * block // 每个线程的数据偏移
			for i := 0; i < block; i++ {
				sn := uint64(i + iOffset)      // 实际的数据sn
				v := confuse(sn, offset) % mod // 计算结果

				chi := int(v) / bfVCount // 不同的值放到不同的bf, 这里应该值来确定ch. 根据协程id来确定ch没有准确性
				ch := outCh[chi]
				ch <- v
			}

			pv := atomic.AddInt32(&p, 1)
			t.Log(pv)
			wgTh.Done()
		}(ti)
	}

	// 等待数据计算完成
	wgTh.Wait()

	// 关闭协程
	for i := 0; i < bfCount; i++ {
		close(outCh[i])
	}

	// 等待数据汇总完成
	wgOut.Wait()

	// 计算最终数据是否一致
	for _, bf := range bfs {
		if len(bf) != bfVCount {
			t.Fatalf("len(bf) != bfVCount, len(bf)=%d, bfVCount=%d", len(bf), int(bfVCount))
		}
	}
}

func TestConfuse(t *testing.T) {
	bf := make(map[uint64]struct{}, testCount) // 用于检查是否有重复
	for i := uint64(0); i < testCount; i++ {
		v := Confuse(i, testSeed)
		bf[v] = struct{}{}
	}

	if len(bf) != testCount {
		t.Fatalf("len(bf) != testCount, len(bf)=%d, testCount=%d", len(bf), int(testCount))
	}
}

func TestConfuseMax(t *testing.T) {
	bf := make(map[uint64]struct{}, testCount) // 用于检查是否有重复
	for i := uint64(0); i < testCount; i++ {
		v := Confuse(i+(1e18-testCount), testSeed)

		bf[v] = struct{}{}
	}

	if len(bf) != testCount {
		t.Fatalf("len(bf) != testCount, len(bf)=%d, testCount=%d", len(bf), int(testCount))
	}
}

func TestConfusePrint(t *testing.T) {
	const testCount = 10
	for i := uint64(0); i < testCount; i++ {
		v := Confuse(i, testSeed)
		t.Log("NewRandom=", v)
	}
}

func BenchmarkConfuse(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Confuse(uint64(1), testSeed)
		}
	})
}

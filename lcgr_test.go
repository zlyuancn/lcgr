package lcgr

import (
	"sync"
	"sync/atomic"
	"testing"
)

const testCount uint64 = 1e8
const testSeed uint64 = 0xabcdef12 // 一个seed

// 计算基础理论是否正确
func TestBasic(t *testing.T) {
	const testCount uint64 = 1e9
	const mod = testCount
	const offset = testSeed % mod

	// 计算重复性
	bf := [testCount]bool{}
	for i := uint64(0); i < testCount; i++ {
		sn := i                             // 实际的数据sn
		v := _mod(confuse(sn, offset), mod) // 计算结果
		if bf[v] {
			t.Fatalf("Confusion and conflict between sn=%d", i)
		}
		bf[v] = true
	}
}

func TestConfuse(t *testing.T) {
	bf := make(map[uint64]struct{}, testCount) // 用于检查是否有重复
	for i := uint64(0); i < testCount; i++ {
		v := Confuse(i, testSeed)
		bf[v] = struct{}{}
	}

	if len(bf) != int(testCount) {
		t.Fatalf("len(bf) != testCount, len(bf)=%d, testCount=%d", len(bf), int(testCount))
	}
}

func TestConfuseMax(t *testing.T) {
	bf := make(map[uint64]struct{}, testCount) // 用于检查是否有重复
	for i := uint64(0); i < testCount; i++ {
		v := Confuse(i+(1e18-testCount), testSeed)

		bf[v] = struct{}{}
	}

	if len(bf) != int(testCount) {
		t.Fatalf("len(bf) != testCount, len(bf)=%d, testCount=%d", len(bf), int(testCount))
	}
}

func TestConfuse1e9(t *testing.T) {
	return // 需要验证时注释掉这一行

	const testCount = 1e9

	// 计算重复性
	const bfCount int = 1e5 // bf数量
	bfs := [bfCount]map[uint64]struct{}{}
	outCh := [bfCount]chan uint64{}
	for i := 0; i < bfCount; i++ {
		bfs[i] = make(map[uint64]struct{}, uint64(float64(testCount/bfCount)*1.01)) // 平均每个bf有多少数据, 冗余一点空间避免扩容
		outCh[i] = make(chan uint64, 1000)
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
	const threadCount = 10000
	wgTh := new(sync.WaitGroup)
	wgTh.Add(threadCount)
	for ti := 0; ti < threadCount; ti++ {
		block := testCount / uint64(threadCount) // 每个线程计算多少数据
		go func(ti int) {
			iOffset := uint64(ti) * block // 每个线程的数据偏移
			for i := uint64(0); i < block; i++ {
				sn := i + iOffset          // 实际的数据sn
				v := Confuse(sn, testSeed) // 计算结果

				chi := _mod(v, uint64(bfCount)) // 不同的值放到不同的bf, 这里应该用值来确定ch以保证每个bf之间的值不会重复. 根据协程id来确定ch没有准确性
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
	c := 0
	for _, bf := range bfs {
		c += len(bf)
	}
	if c != testCount {
		t.Fatalf("len(bf) != bfVCount, len(bf)=%d, testCount=%d", c, int(testCount))
	}
}

func TestConfusePrint(t *testing.T) {
	const testCount = 10
	for i := uint64(0); i < testCount; i++ {
		v := Confuse(i, testSeed)
		t.Log("NewRandom=", v)
	}
}

// 计算限制1~8长度的正确性
func TestConfuseLimitLenL(t *testing.T) {
	var testCount = uint64(1)
	for i := uint64(1); i <= 8; i++ {
		limitLen := i
		testCount *= 10

		// 计算重复性
		bf := make([]bool, testCount)
		for sn := uint64(0); sn < testCount; sn++ {
			v := ConfuseLimitLen(sn, testSeed, limitLen) // 计算结果
			if bf[v] {
				t.Fatalf("Confusion and conflict between sn=%d", sn)
			}
			bf[v] = true
		}
		t.Logf("TestConfuseLimitLenLower %d finished", limitLen)
	}
}

// 计算限制9~18长度的正确性
func TestConfuseLimitLenH(t *testing.T) {
	//return // 需要验证时注释掉这一行

	for i := uint64(9); i <= 18; i++ {
		limitLen := i

		// 计算重复性
		bf := make(map[uint64]struct{}, testCount) // 用于检查是否有重复
		for sn := uint64(0); sn < testCount; sn++ {
			v := ConfuseLimitLen(sn, testSeed, limitLen) // 计算结果
			bf[v] = struct{}{}
		}

		if len(bf) != int(testCount) {
			t.Fatalf("len(bf) != testCount, len(bf)=%d, testCount=%d", len(bf), int(testCount))
		}
		t.Logf("TestConfuseLimitLenLower %d finished", limitLen)
	}
}

func TestConfuseLimitLenPrint(t *testing.T) {
	var testCount = uint64(1)
	for i := uint64(0); i < 9; i++ {
		limitLen := i + 1
		testCount *= 10

		for sn := uint64(0); sn < 10; sn++ {
			v := ConfuseLimitLen(sn, testSeed, limitLen) // 计算结果
			t.Logf("ConfuseLimitLenPrint len=%d v=%d", limitLen, v)
		}
	}
}

func BenchmarkConfuse(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Confuse(uint64(1), testSeed)
		}
	})
}

func BenchmarkConfuseLimitLen(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ConfuseLimitLen(uint64(1), testSeed, 9)
		}
	})
}

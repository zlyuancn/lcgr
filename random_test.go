package lcgr

import (
	"math/rand"
	"testing"
)

func TestRandom(t *testing.T) {
	r := NewRandom(testSeed)

	bf := make(map[uint64]struct{}, testCount) // 用于检查是否有重复
	for i := uint64(0); i < testCount; i++ {
		v := r.Next()
		bf[v] = struct{}{}
	}

	if len(bf) != int(testCount) {
		t.Fatalf("len(bf) != testCount, len(bf)=%d, testCount=%d", len(bf), int(testCount))
	}
}

func TestRandomStartSn(t *testing.T) {
	r := NewRandom(testSeed)

	bf := make(map[uint64]struct{}, testCount) // 用于检查是否有重复
	for i := uint64(0); i < testCount; i++ {
		if i == testCount/2 {
			break
		}
		v := r.Next()
		bf[v] = struct{}{}
	}
	sn := r.GetNextSn()

	r2 := NewRandomStartSn(testSeed, sn)
	for i := testCount / 2; i < testCount; i++ {
		v := r2.Next()
		bf[v] = struct{}{}
	}

	if len(bf) != int(testCount) {
		t.Fatalf("len(bf) != testCount, len(bf)=%d, testCount=%d", len(bf), int(testCount))
	}
}

func TestRandomPrint(t *testing.T) {
	const testCount = 10
	r := NewRandom(testSeed)

	for i := uint64(0); i < testCount; i++ {
		v := r.Next()
		t.Log("NewRandom=", v)
	}

}

func TestRandomLimitLen(t *testing.T) {
	const limitLen = 8
	r := NewRandomLimitLen(testSeed, limitLen)

	bf := make(map[uint64]struct{}, testCount) // 用于检查是否有重复
	for i := uint64(0); i < testCount; i++ {
		v := r.NextLimit()
		bf[v] = struct{}{}
	}

	if len(bf) != int(testCount) {
		t.Fatalf("len(bf) != testCount, len(bf)=%d, testCount=%d", len(bf), int(testCount))
	}
}

func TestRandomLimitLenStartSn(t *testing.T) {
	const limitLen = 8
	r := NewRandomLimitLen(testSeed, limitLen)

	bf := make(map[uint64]struct{}, testCount) // 用于检查是否有重复
	for i := uint64(0); i < testCount; i++ {
		if i == testCount/2 {
			break
		}
		v := r.NextLimit()
		bf[v] = struct{}{}
	}
	sn := r.GetNextSn()

	r2 := NewRandomLimitLenStartSn(testSeed, limitLen, sn)
	for i := testCount / 2; i < testCount; i++ {
		v := r2.NextLimit()
		bf[v] = struct{}{}
	}

	if len(bf) != int(testCount) {
		t.Fatalf("len(bf) != testCount, len(bf)=%d, testCount=%d", len(bf), int(testCount))
	}
}

func TestRandomLimitLenPrint(t *testing.T) {
	const testCount = 10
	const limitLen = 5
	r := NewRandomLimitLen(testSeed, limitLen)

	for i := uint64(0); i < testCount; i++ {
		v := r.NextLimit()
		t.Log("NewRandom=", v)
	}

}

func BenchmarkRandom(b *testing.B) {
	r := NewRandom(testSeed)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Next()
		}
	})
}

func BenchmarkStdRandom(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rand.Int63()
		}
	})
}

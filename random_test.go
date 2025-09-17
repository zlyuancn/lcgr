package cnum

import (
	"math"
	"math/rand"
	"testing"
)

func TestRandom(t *testing.T) {
	r := NewRandom(testSeed)

	bf := make(map[uint64]struct{}, testCount) // 用于检查是否有重复
	for i := uint64(0); i < testCount; i++ {
		v := r.Uint64()
		bf[v] = struct{}{}
	}

	if len(bf) != testCount {
		t.Fatalf("len(bf) != testCount, len(bf)=%d, testCount=%d", len(bf), int(testCount))
	}
}

func TestRandomSetCount(t *testing.T) {
	r := NewRandom(testSeed)

	bf := make(map[uint64]struct{}, testCount) // 用于检查是否有重复
	for i := uint64(0); i < testCount; i++ {
		if i == testCount/2 {
			break
		}
		v := r.Uint64()
		bf[v] = struct{}{}
	}
	count := r.GetCount()

	r2 := NewRandomSetCount(testSeed, count)
	for i := testCount / 2; i < testCount; i++ {
		v := r2.Uint64()
		bf[v] = struct{}{}
	}

	if len(bf) != testCount {
		t.Fatalf("len(bf) != testCount, len(bf)=%d, testCount=%d", len(bf), int(testCount))
	}
}

func TestRandomN(t *testing.T) {
	r := NewRandomN(testSeed, testCount)

	bf := make(map[uint64]struct{}, testCount) // 用于检查是否有重复
	for i := uint64(0); i < testCount; i++ {
		v := r.Uint64N()
		bf[v] = struct{}{}
	}

	if len(bf) != testCount {
		t.Fatalf("len(bf) != testCount, len(bf)=%d, testCount=%d", len(bf), int(testCount))
	}
}

func TestRandomNSetCount(t *testing.T) {
	r := NewRandomN(testSeed, testCount)

	bf := make(map[uint64]struct{}, testCount) // 用于检查是否有重复
	for i := uint64(0); i < testCount; i++ {
		if i == testCount/2 {
			break
		}
		v := r.Uint64N()
		bf[v] = struct{}{}
	}
	count := r.GetCount()

	r2 := NewRandomNSetCount(testSeed, testCount, count)
	for i := testCount / 2; i < testCount; i++ {
		v := r2.Uint64N()
		bf[v] = struct{}{}
	}

	if len(bf) != testCount {
		t.Fatalf("len(bf) != testCount, len(bf)=%d, testCount=%d", len(bf), int(testCount))
	}
}

func TestRandomPrint(t *testing.T) {
	const testCount = 100
	r := NewRandom(testSeed)

	for i := uint64(0); i < testCount; i++ {
		v := r.Uint64()
		t.Log("NewRandom=", v)
	}

}
func TestRandomNPrint(t *testing.T) {
	const testCount = 100
	r := NewRandomN(testSeed, testCount)

	for i := uint64(0); i < testCount; i++ {
		v := r.Uint64N()
		t.Log("NewRandomN=", v)
	}
}

func BenchmarkRandom(b *testing.B) {
	r := NewRandom(testSeed)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Uint64()
		}
	})
}

func BenchmarkRandomN(b *testing.B) {
	r := NewRandomN(testSeed, testCount)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Uint64N()
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

func BenchmarkStdRandomN(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rand.Int63n(math.MaxInt)
		}
	})
}

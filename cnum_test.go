package cnum

import (
	"math"
	"testing"
)

const testCount = 1e7
const testSeed = 0xabcdef12 // 一个seed

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
		v := Confuse(i+(math.MaxInt-1e8), testSeed)

		bf[v] = struct{}{}
	}

	if len(bf) != testCount {
		t.Fatalf("len(bf) != testCount, len(bf)=%d, testCount=%d", len(bf), int(testCount))
	}
}

func TestConfuseLimit(t *testing.T) {
	bf := make(map[uint64]struct{}, testCount) // 用于检查是否有重复
	for i := uint64(0); i < testCount; i++ {
		v := ConfuseLimit(i, testSeed, testCount)
		bf[v] = struct{}{}
	}

	if len(bf) != testCount {
		t.Fatalf("len(bf) != testCount, len(bf)=%d, testCount=%d", len(bf), int(testCount))
	}
}

func TestConfuseLimitP(t *testing.T) {
	bf := make(map[uint64]struct{}, testCount) // 用于检查是否有重复
	for i := uint64(0); i < testCount; i++ {
		v := ConfuseLimitP(i, testSeed, testCount, 9)
		bf[v] = struct{}{}
	}

	if len(bf) != testCount {
		t.Fatalf("len(bf) != testCount, len(bf)=%d, testCount=%d", len(bf), int(testCount))
	}
}

func BenchmarkConfuse(b *testing.B) {
	i := 0
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Confuse(uint64(i%testCount), testSeed)
			i++
		}
	})
}

func BenchmarkConfuseLimit(b *testing.B) {
	i := 0
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ConfuseLimit(uint64(i%testCount), testSeed, testCount)
			i++
		}
	})
}

func BenchmarkConfuseLimitP(b *testing.B) {
	i := 0
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ConfuseLimitP(uint64(i%testCount), testSeed, testCount, 10)
			i++
		}
	})
}

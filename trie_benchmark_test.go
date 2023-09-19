package glsmt

import (
	"math/rand"
	"testing"
	"time"
	"unsafe"
)

var (
	ttl time.Duration = 50 * time.Millisecond

	parallelism = 10000

	bigData      = map[string]string{}
	bigDataLen   = 2 << 10
	bigDataCount = 2 << 16

	smallData = map[string]string{
		"string": "aaaa",
		"int":    "123",
		"float":  "99.99",
		"struct": "struct{}{}",
	}
)

func init() {
	for i := 0; i < bigDataCount; i++ {
		bigData[randStr(bigDataLen)] = randStr(bigDataLen)
	}
}

var randSrc = rand.NewSource(time.Now().UnixNano())

const (
	rs6Letters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rs6LetterIdxBits = 6
	rs6LetterIdxMask = 1<<rs6LetterIdxBits - 1
	rs6LetterIdxMax  = 63 / rs6LetterIdxBits
)

func randStr(n int) string {
	b := make([]byte, n)
	cache, remain := randSrc.Int63(), rs6LetterIdxMax
	for i := n - 1; i >= 0; {
		if remain == 0 {
			cache, remain = randSrc.Int63(), rs6LetterIdxMax
		}
		idx := int(cache & rs6LetterIdxMask)
		if idx < len(rs6Letters) {
			b[i] = rs6Letters[idx]
			i--
		}
		cache >>= rs6LetterIdxBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

func benchmark(b *testing.B, data map[string]string,
	set func(string, string),
	get func(string),
) {
	b.Helper()
	b.SetParallelism(parallelism)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for k, v := range data {
				set(k, v)
			}
			for k := range data {
				get(k)
			}
		}
	})
}

func BenchmarkTrieSetSmallData(b *testing.B) {
	m := NewTrie[string](16)
	benchmark(b, smallData,
		func(k, v string) { m.Insert(k, &v) },
		func(k string) { m.Get(k) })
}

func BenchmarkTrieSetBigData(b *testing.B) {
	m := NewTrie[string](32)
	benchmark(b, bigData,
		func(k, v string) { m.Insert(k, &v) },
		func(k string) { m.Get(k) })
}

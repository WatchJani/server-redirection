package consistent_hashing

import (
	"hash/crc32"
	"log"
	"testing"
)

const TestKeyValue = "key"

// 73ns
func BenchmarkFindServer(b *testing.B) {
	b.StopTimer()
	ch := NewConsistentHashing()
	if err := ch.Load(ConfigFile); err != nil {
		log.Println(err)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		ch.FindServer(TestKeyValue)
	}
}

// 12 ns/op
func BenchmarkSearch(b *testing.B) {
	b.StopTimer()

	ch := NewConsistentHashing()
	if err := ch.Load(ConfigFile); err != nil {
		log.Println(err)
	}

	newHash := crc32.ChecksumIEEE([]byte(TestKeyValue))

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		Search(ch.VirtualServers, newHash)
	}
}

// 28 ns/op
func BenchmarkHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		crc32.ChecksumIEEE([]byte(TestKeyValue))
	}
}

// 12061973 ns/op
func BenchmarkInsert(b *testing.B) {
	b.StopTimer()
	ch := NewConsistentHashing()

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		ch.Load(ConfigFile)
	}
}

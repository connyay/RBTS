package rbts

import (
	"hash/fnv"
	"io"
	"log"
	"math/rand"
	"testing"

	"github.com/google/uuid"
)

// TruncatedFNV64 is the hash mechanism used by Tanium. This is a standard 64bit
// FNVa hash that is truncated to 32bits.
func TruncatedFNV64(s string) uint32 {
	h := fnv.New64a()
	io.WriteString(h, s)
	return uint32(h.Sum64())
}

func TestNewSeriesRandomVals(t *testing.T) {
	s := New(1)
	for i := 2; i < 255; i++ {
		s.Push(uint8(i), TruncatedFNV64(uuid.New().String()))
		// log.Println(s.Bytes())
	}
	// log.Println(s.Bytes())
	log.Printf("%d bytes (TestNewSeriesRandomVals)", len(s.Bytes()))

	// it := s.Iter()
	// for it.Next() {
	// 	id, val := it.Values()
	// 	log.Printf("id=%d val=%d", id, val)
	// }
}

func TestNewSeriesSameVal(t *testing.T) {
	s := New(1)
	val := TruncatedFNV64(uuid.New().String())
	for i := 2; i < 255; i++ {
		s.Push(uint8(i), val)
		// log.Println(s.Bytes())
	}
	// log.Println(s.Bytes())
	log.Printf("%d bytes (TestNewSeriesSameVal)", len(s.Bytes()))

	// it := s.Iter()
	// for it.Next() {
	// 	id, val := it.Values()
	// 	log.Printf("id=%d val=%d", id, val)
	// }
}

func TestNewSeriesFlappingVal(t *testing.T) {
	s := New(1)
	valE := TruncatedFNV64(uuid.New().String())
	valO := TruncatedFNV64(uuid.New().String())
	for i := 2; i < 255; i++ {
		if i%2 == 0 {
			s.Push(uint8(i), valE)
		} else {
			s.Push(uint8(i), valO)
		}
		// log.Println(s.Bytes())
	}
	// log.Println(s.Bytes())
	log.Printf("%d bytes (TestNewSeriesFlappingVal)", len(s.Bytes()))

	// it := s.Iter()
	// for it.Next() {
	// 	id, val := it.Values()
	// 	log.Printf("id=%d val=%d", id, val)
	// }
}

func TestNewSeriesRandomFlappingVal(t *testing.T) {
	s := New(1)
	valA := TruncatedFNV64(uuid.New().String())
	valB := TruncatedFNV64(uuid.New().String())
	for i := 2; i < 255; i++ {
		if rand.Intn(2) == 0 {
			s.Push(uint8(i), valA)
		} else {
			s.Push(uint8(i), valB)
		}
		// log.Println(s.Bytes())
	}
	// log.Println(s.Bytes())
	log.Printf("%d bytes (TestNewSeriesRandomFlappingVal)", len(s.Bytes()))

	// it := s.Iter()
	// for it.Next() {
	// 	id, val := it.Values()
	// 	log.Printf("id=%d val=%d", id, val)
	// }
}

package rbts

import "sync"

// LIFTED FROM https://github.com/dgryski/go-tsz

type Series struct {
	sync.Mutex

	ID0 uint8
	id  uint8
	val uint32

	bw       bstream
	finished bool
}

func New(id0 uint8) *Series {
	s := Series{
		ID0: id0,
	}

	// block header
	s.bw.writeByte(id0)

	return &s
}

// Bytes value of the series stream
func (s *Series) Bytes() []byte {
	s.Lock()
	defer s.Unlock()
	return s.bw.bytes()
}

func finish(w *bstream) {
	// write an end-of-stream record
	w.writeBits(0x0f, 4)
	w.writeBits(0xffffffff, 32)
	w.writeBit(zero)
}

func (s *Series) Finish() {
	s.Lock()
	if !s.finished {
		finish(&s.bw)
		s.finished = true
	}
	s.Unlock()
}

// Push an id and value to the series
func (s *Series) Push(id uint8, v uint32) {
	s.Lock()
	defer s.Unlock()

	if s.id == 0 {
		// first point
		s.id = id
		s.val = v
		s.bw.writeByte(id)
		s.bw.writeBits(uint64(v), 32)
		return
	}
	s.bw.writeByte(id)
	if v == s.val {
		s.bw.writeBit(zero)
	} else {
		s.bw.writeBit(one)
		s.bw.writeBits(uint64(v), 32)
		s.val = v
	}
}

func (s *Series) Iter() *Iter {
	s.Lock()
	w := s.bw.clone()
	s.Unlock()

	finish(w)
	iter, _ := bstreamIterator(w)
	return iter
}

// Iter lets you iterate over a series.  It is not concurrency-safe.
type Iter struct {
	ID0 uint8

	id  uint8
	val uint32

	br bstream

	finished bool

	err error
}

func bstreamIterator(br *bstream) (*Iter, error) {

	br.count = 8

	id0, err := br.readByte()
	if err != nil {
		return nil, err
	}

	return &Iter{
		ID0: uint8(id0),
		br:  *br,
	}, nil
}

func (it *Iter) Next() bool {
	if it.err != nil || it.finished {
		return false
	}

	if it.id == 0 {
		// read first id and v
		id, err := it.br.readByte()
		if err != nil {
			it.err = err
			return false
		}
		it.id = uint8(id)
		v, err := it.br.readBits(32)
		if err != nil {
			it.err = err
			return false
		}

		it.val = uint32(v)

		return true
	}

	id, err := it.br.readByte()
	if err != nil {
		it.err = err
		return false
	}
	it.id = uint8(id)

	// read compressed value
	bit, err := it.br.readBit()
	if err != nil {
		it.err = err
		return false
	}

	if bit == zero {
		// it.val = it.val
	} else {
		v, err := it.br.readBits(32)
		if err != nil {
			it.err = err
			return false
		}

		it.val = uint32(v)
	}

	return true
}

func (it *Iter) Values() (uint8, uint32) {
	return it.id, it.val
}

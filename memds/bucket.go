package memds

import (
	"hash/crc32"
	"sync"

	"github.com/ugorji/go/codec"
)

type Bucket struct {
	mu    *sync.RWMutex
	value map[string][]byte
}

type Buckets []*Bucket

func NewBuckets(n int) (Buckets, error) {
	if n <= 0 {
		return nil, BucketsLEZeroError
	}

	b := make(Buckets, 0, n)
	for i := 0; i < n; i++ {
		b = append(
			b,
			newBucket(),
		)
	}
	return b, nil
}

func (b Buckets) Get(k string) *Bucket {
	if len(b) == 0 {
		return nil
	}
	n := crc32.ChecksumIEEE([]byte(k))
	i := int(n % uint32(b.Len()))
	return b[i]
}

func (b Buckets) Len() int {
	return len(b)
}

func newBucket() *Bucket {
	b := Bucket{
		mu:    new(sync.RWMutex),
		value: make(map[string][]byte),
	}
	return &b
}

func (b *Bucket) Get(k string) (interface{}, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	v, ok := b.value[k]
	if ok {
		var r interface{}
		dec := codec.NewDecoderBytes(v, &mh)
		if err := dec.Decode(&r); err != nil {
			return nil, err
		}
		return r, nil
	} else {
		return nil, ValueNotFoundError
	}
}

func (b *Bucket) Set(k string, v interface{}) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	var bs []byte
	enc := codec.NewEncoderBytes(&bs, &mh)
	if err := enc.Encode(v); err != nil {
		return err
	}
	b.value[k] = bs
	return nil
}

func (b *Bucket) Del(k string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.value, k)
}

package memds

import (
	"hash/crc32"
	"sync"
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
		value := Value{}
		err := value.Decode(v)
		if err != nil {
			return nil, err
		}
		if value.IsExpire() {
			return nil, ValueNotFoundError
		}
		return value.Value(), nil
	} else {
		return nil, ValueNotFoundError
	}
}

func (b *Bucket) Set(k string, v interface{}, es int64) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	value := NewValue(v, es)
	bs, err := value.Encode()
	if err != nil {
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

package memds

import (
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/ugorji/go/codec"
)

func TestNewBuckets(t *testing.T) {
	testCase := []struct {
		In  int
		Len int
		Err error
	}{
		{
			In:  5,
			Len: 5,
			Err: nil,
		},
		{
			In:  0,
			Len: 0,
			Err: BucketsLEZeroError,
		},
	}
	for _, tc := range testCase {
		b, err := NewBuckets(tc.In)
		if err != tc.Err {
			t.Errorf("got: %v, want: %v", err, tc.Err)
		}
		if len(b) != tc.Len {
			t.Errorf("got: %v, want: %v", len(b), tc.Len)
		}
	}
}

func TestBucketsLen(t *testing.T) {
	b, _ := NewBuckets(2)
	testCase := []struct {
		In     Buckets
		Result int
	}{
		{
			In:     b,
			Result: 2,
		},
	}

	for _, tc := range testCase {
		if tc.In.Len() != tc.Result {
			t.Errorf("got: %v, want: %v", tc.In.Len(), tc.Result)
		}
	}
}

func TestBucketsGet(t *testing.T) {
	b := make(Buckets, 0, 1)

	if b.Get("hoge") != nil {
		t.Error("should be nil when empty buckets")
	}

	b = append(b, newBucket())

	if b.Get("hoge") == nil {
		t.Error("should not be nil")
	}

	b = append(b, newBucket())

	if b.Get("hoge") == nil {
		t.Error("should not be nil")
	}
}

func TestNewBucket(t *testing.T) {
	b := newBucket()
	if b == nil {
		t.Error("newBucket return nil")
	}
}

func TestBucketGet(t *testing.T) {
	mu := new(sync.RWMutex)
	var bs []byte
	codec.NewEncoderBytes(&bs, &mh).Encode(
		map[string]interface{}{
			"value":  "value",
			"expire": int64(0),
		},
	)
	b := Bucket{
		mu: mu,
		value: map[string][]byte{
			"key": bs,
		},
	}

	testCase := []struct {
		In  string
		V   interface{}
		Err error
	}{
		{
			In:  "key",
			V:   []byte("value"),
			Err: nil,
		},
		{
			In:  "key1",
			V:   nil,
			Err: ValueNotFoundError,
		},
	}

	for _, tc := range testCase {
		v, err := b.Get(tc.In)
		if err != tc.Err {
			t.Errorf("got: %v, want: %v", err, tc.Err)
		}
		if !reflect.DeepEqual(v, tc.V) {
			t.Errorf("got: %v, want: %v", v, tc.V)
		}
	}
}

func TestBucketSet(t *testing.T) {
	mu := new(sync.RWMutex)
	b := Bucket{
		mu:    mu,
		value: map[string][]byte{},
	}

	u := time.Now().UTC().Unix()
	v0 := Value{
		"value":  []byte("value"),
		"expire": int64(0),
	}
	v1 := Value{
		"value":  []byte("value1"),
		"expire": u + int64(10),
	}
	v2 := Value{
		"value":  []byte("value2"),
		"expire": u + int64(100),
	}

	testCase := []struct {
		Key    string
		Value  interface{}
		Expire int64
		Result map[string]Value
	}{
		{
			Key:    "key",
			Value:  "value",
			Expire: int64(0),
			Result: map[string]Value{
				"key": v0,
			},
		},
		{
			Key:    "key1",
			Value:  "value1",
			Expire: int64(10),
			Result: map[string]Value{
				"key":  v0,
				"key1": v1,
			},
		},
		{
			Key:    "key",
			Value:  "value2",
			Expire: int64(100),
			Result: map[string]Value{
				"key":  v2,
				"key1": v1,
			},
		},
	}

	for _, tc := range testCase {
		b.Set(tc.Key, tc.Value, tc.Expire)

		for k, v := range tc.Result {
			bs, ok := b.value[k]
			if !ok {
				t.Errorf("key(%s) not found.", k)
			}
			val := Value{}
			val.Decode(bs)
			if !reflect.DeepEqual(val.Value(), v.Value()) {
				t.Errorf("got: %v, want: %v", val.Value(), v.Value())
			}
			if val.ExpireAt().Int < v.ExpireAt().Int {
				t.Errorf("got: %v, want: %v", val.ExpireAt(), v.ExpireAt())
			}
		}
	}
}

func TestBucketDel(t *testing.T) {
	mu := new(sync.RWMutex)
	b := Bucket{
		mu: mu,
		value: map[string][]byte{
			"key":  []byte("value"),
			"key1": []byte("value1"),
		},
	}

	testCase := []struct {
		Key    string
		Result map[string][]byte
	}{
		{
			Key: "key3",
			Result: map[string][]byte{
				"key":  []byte("value"),
				"key1": []byte("value1"),
			},
		},
		{
			Key: "key1",
			Result: map[string][]byte{
				"key": []byte("value"),
			},
		},
	}

	for _, tc := range testCase {
		b.Del(tc.Key)

		if !reflect.DeepEqual(b.value, tc.Result) {
			t.Errorf("got: %v, want: %v", b.value, tc.Result)
		}
	}
}

func BenchmarkBucketGet(b *testing.B) {
	mu := new(sync.RWMutex)
	keys := make([]string, 0, b.N)
	v := make(map[string][]byte, b.N)
	for i := 0; i < b.N; i++ {
		var n uint64
		binary.Read(crand.Reader, binary.LittleEndian, &n)
		s := strconv.FormatUint(n, 36)
		v[s] = []byte(s)
		keys = append(keys, s)
	}

	perm := rand.Perm(len(keys))
	for i, v := range perm {
		keys[i], keys[v] = keys[v], keys[i]
	}

	buc := Bucket{
		mu:    mu,
		value: v,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buc.Get(keys[i])
	}
}

func BenchmarkBucketSet(b *testing.B) {
	v := make(map[string][]byte, b.N)
	for i := 0; i < b.N; i++ {
		var n uint64
		binary.Read(crand.Reader, binary.LittleEndian, &n)
		s := strconv.FormatUint(n, 36)
		v[s] = []byte(s)
	}

	buc := newBucket()
	b.ResetTimer()
	for k, v := range v {
		buc.Set(k, v, int64(10))
	}
}

func BenchmarkBucketDel(b *testing.B) {
	mu := new(sync.RWMutex)
	keys := make([]string, 0, b.N)
	v := make(map[string][]byte, b.N)
	for i := 0; i < b.N; i++ {
		var n uint64
		binary.Read(crand.Reader, binary.LittleEndian, &n)
		s := strconv.FormatUint(n, 36)
		v[s] = []byte(s)
		keys = append(keys, s)
	}

	perm := rand.Perm(len(keys))
	for i, v := range perm {
		keys[i], keys[v] = keys[v], keys[i]
	}

	buc := Bucket{
		mu:    mu,
		value: v,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buc.Del(keys[i])
	}
}

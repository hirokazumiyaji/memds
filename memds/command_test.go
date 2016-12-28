package memds

import (
	crand "crypto/rand"
	"encoding/binary"
	"reflect"
	"strconv"
	"testing"

	"github.com/ugorji/go/codec"
)

func TestGetAndSet(t *testing.T) {
	buckets, _ = NewBuckets(10)

	setTestCase := []struct {
		Key   string
		Value interface{}
		Err   error
	}{
		{
			Key:   "key",
			Value: []byte("value"),
			Err:   nil,
		},
		{
			Key:   "key1",
			Value: []byte("value1"),
			Err:   nil,
		},
		{
			Key:   "key2",
			Value: []byte("value2"),
			Err:   nil,
		},
		{
			Key:   "key3",
			Value: []byte("value3"),
			Err:   nil,
		},
	}

	getTestCase := []struct {
		Key   string
		Value interface{}
		Err   error
	}{
		{
			Key:   "key",
			Value: []byte("value"),
			Err:   nil,
		},
		{
			Key:   "key1",
			Value: []byte("value1"),
			Err:   nil,
		},
		{
			Key:   "key2",
			Value: []byte("value2"),
			Err:   nil,
		},
		{
			Key:   "key3",
			Value: []byte("value3"),
			Err:   nil,
		},
		{
			Key:   "key4",
			Value: nil,
			Err:   ValueNotFoundError,
		},
	}

	for _, tc := range setTestCase {
		err := Set(tc.Key, tc.Value)
		if err != tc.Err {
			t.Errorf("key: %v, got: %v, want: %v", tc.Key, err, tc.Err)
		}
	}

	for _, tc := range getTestCase {
		v, err := Get(tc.Key)
		if err != tc.Err {
			t.Errorf("key: %v, got: %v, want: %v", tc.Key, err, tc.Err)
		}
		if !reflect.DeepEqual(v, tc.Value) {
			t.Errorf("key: %v, got: %v, want: %v", tc.Key, v, tc.Value)
		}
	}

	buckets = make(Buckets, 0, 10)
	err := Set("key", []byte("value"))
	if err != BucketNotFoundError {
		t.Errorf("got: %v, want: %v", err, BucketNotFoundError)
	}
	v, err := Get("key")
	if err != BucketNotFoundError {
		t.Errorf("got: %v, want: %v", err, BucketNotFoundError)
	}
	if v != nil {
		t.Errorf("got: %v, want: nil", v)
	}
}

func TestDel(t *testing.T) {
	buckets, _ = NewBuckets(10)

	Set("key", []byte("value"))

	err := Del("key")
	if err != nil {
		t.Errorf("got: %v, want: nil", err)
	}

	_, err = Get("key")
	if err != ValueNotFoundError {
		t.Errorf("got :%v, want: %v", err, ValueNotFoundError)
	}

	err = Del("key1")
	if err != nil {
		t.Errorf("got: %v, want: nil", err)
	}

	buckets = make(Buckets, 0, 10)
	err = Del("key")
	if err != BucketNotFoundError {
		t.Errorf("got: %v, want: %v", err, BucketNotFoundError)
	}
}

func BenchmarkSet(b *testing.B) {
	v := make(map[string][]byte, b.N)
	for i := 0; i < b.N; i++ {
		var n uint64
		binary.Read(crand.Reader, binary.LittleEndian, &n)
		s := strconv.FormatUint(n, 36)
		v[s] = []byte(s)
	}

	buckets, _ = NewBuckets(10)
	b.ResetTimer()
	for k, v := range v {
		Set(k, v)
	}
}

func TestExec(t *testing.T) {
	func() {
		var b []byte
		enc := codec.NewEncoderBytes(&b, &mh)
		if err := enc.Encode(map[string]interface{}{}); err != nil {
			t.Fatal("command encode error")
		}

		r := Exec(b)

		res := make(map[string]interface{})
		dec := codec.NewDecoderBytes(r, &mh)
		if err := dec.Decode(res); err != nil {
			t.Errorf("response decode error: %v", err)
		}

		s, ok := res["status"]
		if !ok {
			t.Error("response key 'status' not found")
		}
		sb, ok := s.(bool)
		if !ok {
			t.Error("response key 'status' not bool type")
		}
		if sb != false {
			t.Errorf("got: %v, want: %v", sb, false)
		}

		c, ok := res["code"]
		if !ok {
			t.Error("response key 'code' not found")
		}
		ci, ok := c.(uint64)
		if !ok {
			t.Error("response key 'code' not uint64 type")
		}
		if ci != ErrorCodeCommandFormatError {
			t.Errorf("got: %v, want: %v", ci, ErrorCodeCommandFormatError)
		}
	}()

	func() {
		var b []byte
		enc := codec.NewEncoderBytes(&b, &mh)
		if err := enc.Encode(map[string]interface{}{"cmd": "get"}); err != nil {
			t.Fatal("command encode error")
		}

		r := Exec(b)

		res := make(map[string]interface{})
		dec := codec.NewDecoderBytes(r, &mh)
		if err := dec.Decode(res); err != nil {
			t.Errorf("response decode error: %v", err)
		}

		s, ok := res["status"]
		if !ok {
			t.Error("response key 'status' not found")
		}
		sb, ok := s.(bool)
		if !ok {
			t.Error("response key 'status' not bool type")
		}
		if sb != false {
			t.Errorf("got: %v, want: %v", sb, false)
		}

		c, ok := res["code"]
		if !ok {
			t.Error("response key 'code' not found")
		}
		ci, ok := c.(uint64)
		if !ok {
			t.Error("response key 'code' not uint64 type")
		}
		if ci != ErrorCodeCommandFormatError {
			t.Errorf("got: %v, want: %v", ci, ErrorCodeCommandFormatError)
		}
	}()
}

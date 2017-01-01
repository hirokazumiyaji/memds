package memds

import (
	"reflect"
	"testing"
	"time"

	"github.com/ugorji/go/codec"
)

func TestNewValue(t *testing.T) {
	testCase := []struct {
		Value  interface{}
		Expire int64
		Out    Value
	}{
		{
			Value:  int64(10),
			Expire: int64(10),
			Out: Value{
				"value":  int64(10),
				"expire": time.Now().UTC().Unix() + int64(10),
			},
		},
		{
			Value:  int64(10),
			Expire: int64(0),
			Out: Value{
				"value":  int64(10),
				"expire": int64(0),
			},
		},
	}

	for _, tc := range testCase {
		v := NewValue(tc.Value, tc.Expire)
		if !reflect.DeepEqual(v["value"], tc.Out["value"]) {
			t.Errorf("got: %v, want: %v", v["value"], tc.Out["value"])
		}
		e, ok := v["expire"]
		if !ok {
			t.Error("expire field not found")
		}
		actual, _ := tc.Out["expire"].(int64)
		switch e := e.(type) {
		case int64:
			if e < actual {
				t.Error("expire less than")
			}
		default:
			t.Error("expire not int64 type")
		}
	}
}

func TestValueEncode(t *testing.T) {
	var expected []byte
	codec.NewEncoderBytes(&expected, &mh).Encode(
		map[string]interface{}{
			"value":  "value",
			"expire": int64(1483268026),
		},
	)
	v := Value{
		"value":  "value",
		"expire": int64(1483268026),
	}
	b, err := v.Encode()
	if err != nil {
		t.Errorf("encode should not be error: %v", err)
	}
	if !reflect.DeepEqual(expected, b) {
		t.Errorf("got: %v, want: %v", b, expected)
	}
}

func TestValueDecode(t *testing.T) {
	expected := Value{
		"value":  int64(10),
		"expire": int64(10),
	}
	var b []byte
	codec.NewEncoderBytes(&b, &mh).Encode(expected)

	v := Value{}
	err := v.Decode(b)
	if err != nil {
		t.Errorf("decode should not be error: %v", err)
	}

	if !reflect.DeepEqual(v, expected) {
		t.Errorf("got: %v, want: %v", v, expected)
	}
}

func TestValueBool(t *testing.T) {
	testCase := []struct {
		In  Value
		Out NullBool
	}{
		{
			In: Value{},
			Out: NullBool{
				Bool:  false,
				Valid: false,
			},
		},
		{
			In: Value{
				"value": nil,
			},
			Out: NullBool{
				Bool:  false,
				Valid: false,
			},
		},
		{
			In: Value{
				"value": "value",
			},
			Out: NullBool{
				Bool:  false,
				Valid: false,
			},
		},
		{
			In: Value{
				"value": true,
			},
			Out: NullBool{
				Bool:  true,
				Valid: true,
			},
		},
		{
			In: Value{
				"value": "true",
			},
			Out: NullBool{
				Bool:  true,
				Valid: true,
			},
		},
	}

	for _, tc := range testCase {
		got := tc.In.Bool()
		if !reflect.DeepEqual(got, tc.Out) {
			t.Errorf("got: %v, want: %v", got, tc.Out)
		}
	}
}

func TestValueString(t *testing.T) {
	testCase := []struct {
		In  Value
		Out NullString
	}{
		{
			In: Value{},
			Out: NullString{
				String: "",
				Valid:  false,
			},
		},
		{
			In: Value{
				"value": nil,
			},
			Out: NullString{
				String: "",
				Valid:  false,
			},
		},
		{
			In: Value{
				"value": struct{}{},
			},
			Out: NullString{
				String: "",
				Valid:  false,
			},
		},
		{
			In: Value{
				"value": "value",
			},
			Out: NullString{
				String: "value",
				Valid:  true,
			},
		},
		{
			In: Value{
				"value": int64(10),
			},
			Out: NullString{
				String: "10",
				Valid:  true,
			},
		},
		{
			In: Value{
				"value": 10.0,
			},
			Out: NullString{
				String: "1E+01",
				Valid:  true,
			},
		},
		{
			In: Value{
				"value": true,
			},
			Out: NullString{
				String: "true",
				Valid:  true,
			},
		},
	}

	for _, tc := range testCase {
		got := tc.In.String()
		if !reflect.DeepEqual(got, tc.Out) {
			t.Errorf("got: %v, want: %v", got, tc.Out)
		}
	}
}

func TestValueInt(t *testing.T) {
	testCase := []struct {
		In  Value
		Out NullInt
	}{
		{
			In: Value{},
			Out: NullInt{
				Int:   int64(0),
				Valid: false,
			},
		},
		{
			In: Value{
				"value": nil,
			},
			Out: NullInt{
				Int:   int64(0),
				Valid: false,
			},
		},
		{
			In: Value{
				"value": "value",
			},
			Out: NullInt{
				Int:   int64(0),
				Valid: false,
			},
		},
		{
			In: Value{
				"value": int64(10),
			},
			Out: NullInt{
				Int:   int64(10),
				Valid: true,
			},
		},
	}

	for _, tc := range testCase {
		got := tc.In.Int()
		if !reflect.DeepEqual(got, tc.Out) {
			t.Errorf("got: %v, want: %v", got, tc.Out)
		}
	}
}

func TestValueFloat(t *testing.T) {
	testCase := []struct {
		In  Value
		Out NullFloat
	}{
		{
			In: Value{},
			Out: NullFloat{
				Float: float64(0),
				Valid: false,
			},
		},
		{
			In: Value{
				"value": nil,
			},
			Out: NullFloat{
				Float: float64(0),
				Valid: false,
			},
		},
		{
			In: Value{
				"value": "value",
			},
			Out: NullFloat{
				Float: float64(0),
				Valid: false,
			},
		},
		{
			In: Value{
				"value": float64(10),
			},
			Out: NullFloat{
				Float: float64(10),
				Valid: true,
			},
		},
		{
			In: Value{
				"value": "10.0",
			},
			Out: NullFloat{
				Float: float64(10.0),
				Valid: true,
			},
		},
		{
			In: Value{
				"value": int64(10),
			},
			Out: NullFloat{
				Float: float64(10),
				Valid: true,
			},
		},
	}

	for _, tc := range testCase {
		got := tc.In.Float()
		if !reflect.DeepEqual(got, tc.Out) {
			t.Errorf("got: %v, want: %v", got, tc.Out)
		}
	}
}

func TestValueExpireAt(t *testing.T) {
	testCase := []struct {
		In  Value
		Out NullInt
	}{
		{
			In: Value{
				"value": "value",
			},
			Out: NullInt{
				Int:   int64(0),
				Valid: false,
			},
		},
		{
			In: Value{
				"value":  "value",
				"expire": int64(1000),
			},
			Out: NullInt{
				Int:   int64(1000),
				Valid: true,
			},
		},
	}

	for _, tc := range testCase {
		got := tc.In.ExpireAt()
		if !reflect.DeepEqual(got, tc.Out) {
			t.Errorf("got: %v, want: %v", got, tc.Out)
		}
	}
}

func TestIsExpire(t *testing.T) {
	testCase := []struct {
		In  Value
		Out bool
	}{
		{
			In: Value{
				"value":  "value",
				"expire": int64(0),
			},
			Out: false,
		},
		{
			In: Value{
				"value": "value",
			},
			Out: false,
		},
		{
			In: Value{
				"value":  "value",
				"expire": time.Now().UTC().Unix() + 1000,
			},
			Out: false,
		},
		{
			In: Value{
				"value":  "value",
				"expire": time.Now().UTC().Unix() - 1000,
			},
			Out: true,
		},
	}

	for _, tc := range testCase {
		is := tc.In.IsExpire()
		if is != tc.Out {
			t.Errorf("got: %v, want: %v", is, tc.Out)
		}
	}
}

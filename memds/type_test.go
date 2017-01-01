package memds

import (
	"reflect"
	"testing"
)

func TestNullBoolScan(t *testing.T) {
	testCase := []struct {
		In  interface{}
		Out NullBool
	}{
		{
			In: nil,
			Out: NullBool{
				Bool:  false,
				Valid: false,
			},
		},
		{
			In: true,
			Out: NullBool{
				Bool:  true,
				Valid: true,
			},
		},
		{
			In: "true",
			Out: NullBool{
				Bool:  true,
				Valid: true,
			},
		},
		{
			In: "nullbool",
			Out: NullBool{
				Bool:  false,
				Valid: false,
			},
		},
	}

	for _, tc := range testCase {
		n := NullBool{}
		if _ = n.Scan(tc.In); !reflect.DeepEqual(n, tc.Out) {
			t.Errorf("got: %v, want: %v", n, tc.Out)
		}
	}
}

func TestNullBoolValue(t *testing.T) {
	testCase := []struct {
		In  NullBool
		Out interface{}
	}{
		{
			In: NullBool{
				Bool:  false,
				Valid: false,
			},
			Out: nil,
		},
		{
			In: NullBool{
				Bool:  true,
				Valid: true,
			},
			Out: true,
		},
	}

	for _, tc := range testCase {
		if v := tc.In.Value(); v != tc.Out {
			t.Errorf("got: %v, want: %v", v, tc.Out)
		}
	}
}

func TestNullStringScan(t *testing.T) {
	testCase := []struct {
		In  interface{}
		Out NullString
	}{
		{
			In: nil,
			Out: NullString{
				String: "",
				Valid:  false,
			},
		},
		{
			In: "nullstring",
			Out: NullString{
				String: "nullstring",
				Valid:  true,
			},
		},
		{
			In: []uint8{110, 117, 108, 108, 115, 116, 114, 105, 110, 103},
			Out: NullString{
				String: "nullstring",
				Valid:  true,
			},
		},
		{
			In: int(10),
			Out: NullString{
				String: "10",
				Valid:  true,
			},
		},
		{
			In: int64(100),
			Out: NullString{
				String: "100",
				Valid:  true,
			},
		},
		{
			In: float64(3.1415926535),
			Out: NullString{
				String: "3.1415926535E+00",
				Valid:  true,
			},
		},
		{
			In: true,
			Out: NullString{
				String: "true",
				Valid:  true,
			},
		},
		{
			In: struct {
				A int
			}{},
			Out: NullString{
				String: "",
				Valid:  false,
			},
		},
	}

	for _, tc := range testCase {
		n := NullString{}
		if _ = n.Scan(tc.In); !reflect.DeepEqual(n, tc.Out) {
			t.Errorf("got: %v, want: %v", n, tc.Out)
		}
	}
}

func TestNullStringValue(t *testing.T) {
	testCase := []struct {
		In  NullString
		Out interface{}
	}{
		{
			In: NullString{
				String: "",
				Valid:  false,
			},
			Out: nil,
		},
		{
			In: NullString{
				String: "nullstring",
				Valid:  true,
			},
			Out: "nullstring",
		},
	}

	for _, tc := range testCase {
		if v := tc.In.Value(); v != tc.Out {
			t.Errorf("got: %v, want: %v", v, tc.Out)
		}
	}
}

func TestNullIntScan(t *testing.T) {
	testCase := []struct {
		In  interface{}
		Out NullInt
	}{
		{
			In: nil,
			Out: NullInt{
				Int:   int64(0),
				Valid: false,
			},
		},
		{
			In: int64(10),
			Out: NullInt{
				Int:   int64(10),
				Valid: true,
			},
		},
		{
			In: uint64(10),
			Out: NullInt{
				Int:   int64(10),
				Valid: true,
			},
		},
		{
			In: int(10),
			Out: NullInt{
				Int:   int64(10),
				Valid: true,
			},
		},
		{
			In: float64(10.0),
			Out: NullInt{
				Int:   int64(10),
				Valid: true,
			},
		},
		{
			In: "10",
			Out: NullInt{
				Int:   int64(10),
				Valid: true,
			},
		},
		{
			In: []uint8{49, 48},
			Out: NullInt{
				Int:   int64(10),
				Valid: true,
			},
		},
		{
			In: true,
			Out: NullInt{
				Int:   int64(0),
				Valid: false,
			},
		},
	}

	for _, tc := range testCase {
		n := NullInt{}
		if _ = n.Scan(tc.In); !reflect.DeepEqual(n, tc.Out) {
			t.Errorf("got: %v, want: %v", n, tc.Out)
		}
	}
}

func TestNullIntValue(t *testing.T) {
	testCase := []struct {
		In  NullInt
		Out interface{}
	}{
		{
			In: NullInt{
				Int:   int64(0),
				Valid: false,
			},
			Out: nil,
		},
		{
			In: NullInt{
				Int:   int64(10),
				Valid: true,
			},
			Out: int64(10),
		},
	}

	for _, tc := range testCase {
		if v := tc.In.Value(); v != tc.Out {
			t.Errorf("got: %v, want: %v", v, tc.Out)
		}
	}
}

func TestNullFloatScan(t *testing.T) {
	testCase := []struct {
		In  interface{}
		Out NullFloat
	}{
		{
			In: nil,
			Out: NullFloat{
				Float: float64(0),
				Valid: false,
			},
		},
		{
			In: float64(10),
			Out: NullFloat{
				Float: float64(10),
				Valid: true,
			},
		},
		{
			In: int(10),
			Out: NullFloat{
				Float: float64(10),
				Valid: true,
			},
		},
		{
			In: int64(10),
			Out: NullFloat{
				Float: float64(10),
				Valid: true,
			},
		},
		{
			In: uint64(10),
			Out: NullFloat{
				Float: float64(10),
				Valid: true,
			},
		},
		{
			In: "10.0",
			Out: NullFloat{
				Float: float64(10),
				Valid: true,
			},
		},
		{
			In: []uint8{49, 48, 46, 48},
			Out: NullFloat{
				Float: float64(10),
				Valid: true,
			},
		},
		{
			In: true,
			Out: NullFloat{
				Float: float64(0),
				Valid: false,
			},
		},
	}

	for _, tc := range testCase {
		n := NullFloat{}
		if _ = n.Scan(tc.In); !reflect.DeepEqual(n, tc.Out) {
			t.Errorf("got: %v, want: %v", n, tc.Out)
		}
	}
}

func TestNullFloatValue(t *testing.T) {
	testCase := []struct {
		In  NullFloat
		Out interface{}
	}{
		{
			In: NullFloat{
				Float: float64(0),
				Valid: false,
			},
			Out: nil,
		},
		{
			In: NullFloat{
				Float: float64(10),
				Valid: true,
			},
			Out: float64(10),
		},
	}

	for _, tc := range testCase {
		if v := tc.In.Value(); v != tc.Out {
			t.Errorf("got: %v, want: %v", v, tc.Out)
		}
	}
}

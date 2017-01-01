package memds

import (
	"time"

	"github.com/ugorji/go/codec"
)

type Value map[string]interface{}

func NewValue(v interface{}, es int64) Value {
	expire := int64(0)
	if es > 0 {
		expire = time.Now().UTC().Unix() + es
	}
	return Value{
		"value":  v,
		"expire": expire,
	}
}

func (v Value) Encode() ([]byte, error) {
	var b []byte
	enc := codec.NewEncoderBytes(&b, &mh)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return b, nil
}

func (v Value) Decode(b []byte) error {
	dec := codec.NewDecoderBytes(b, &mh)
	return dec.Decode(&v)
}

func (v Value) Value() interface{} {
	value, ok := v["value"]
	if !ok {
		return nil
	}
	return value
}

func (v Value) Bool() NullBool {
	n := NullBool{
		Bool:  false,
		Valid: false,
	}
	value := v.Value()
	n.Scan(value)
	return n
}

func (v Value) String() NullString {
	n := NullString{
		String: "",
		Valid:  false,
	}
	value := v.Value()
	n.Scan(value)
	return n
}

func (v Value) Int() NullInt {
	n := NullInt{
		Int:   int64(0),
		Valid: false,
	}
	value := v.Value()
	n.Scan(value)
	return n
}

func (v Value) Float() NullFloat {
	n := NullFloat{
		Float: float64(0),
		Valid: false,
	}
	value := v.Value()
	n.Scan(value)
	return n
}

func (v Value) ExpireAt() NullInt {
	n := NullInt{
		Int:   int64(0),
		Valid: false,
	}
	e, ok := v["expire"]
	if !ok {
		return n
	}
	n.Scan(e)
	return n
}

func (v Value) IsExpire() bool {
	e := v.ExpireAt()
	if e.Valid == false {
		return false
	}
	if e.Int == 0 {
		return false
	}
	now := time.Now().UTC().Unix()
	return e.Int <= now
}

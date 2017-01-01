package memds

import (
	"errors"
	"strconv"
)

var (
	NotConvertTypeError = errors.New("convert type error")
)

type Nullable interface {
	Scan(interface{}) error
	Value() interface{}
}

type NullBool struct {
	Bool  bool
	Valid bool
}

func (n *NullBool) Scan(v interface{}) error {
	if v == nil {
		n.Bool, n.Valid = false, false
		return nil
	}

	switch v := v.(type) {
	case bool:
		n.Bool, n.Valid = v, true
		return nil
	case string:
		b, err := strconv.ParseBool(v)
		if err != nil {
			n.Bool, n.Valid = false, false
			return err
		}
		n.Bool, n.Valid = b, true
		return err
	default:
		n.Bool, n.Valid = false, false
		return NotConvertTypeError
	}
}

func (n *NullBool) Value() interface{} {
	if n.Valid == false {
		return nil
	}
	return n.Bool
}

type NullString struct {
	String string
	Valid  bool
}

func (n *NullString) Scan(v interface{}) error {
	if v == nil {
		n.String, n.Valid = "", false
		return nil
	}
	switch v := v.(type) {
	case string:
		n.String, n.Valid = v, true
		return nil
	case []uint8:
		s := Uint8ArrayToString(v)
		n.String, n.Valid = s, true
		return nil
	case int:
		s := strconv.Itoa(v)
		n.String, n.Valid = s, true
		return nil
	case int64:
		s := strconv.FormatInt(v, 10)
		n.String, n.Valid = s, true
		return nil
	case float64:
		s := strconv.FormatFloat(v, 'E', -1, 64)
		n.String, n.Valid = s, true
		return nil
	case bool:
		s := strconv.FormatBool(v)
		n.String, n.Valid = s, true
		return nil
	default:
		n.String, n.Valid = "", false
		return NotConvertTypeError
	}
}

func (n *NullString) Value() interface{} {
	if n.Valid == false {
		return nil
	}
	return n.String
}

type NullInt struct {
	Int   int64
	Valid bool
}

func (n *NullInt) Scan(v interface{}) error {
	if v == nil {
		n.Int, n.Valid = int64(0), false
		return nil
	}
	switch v := v.(type) {
	case int:
		n.Int, n.Valid = int64(v), true
		return nil
	case int64:
		n.Int, n.Valid = v, true
		return nil
	case uint64:
		n.Int, n.Valid = int64(v), true
		return nil
	case float64:
		n.Int, n.Valid = int64(v), true
		return nil
	case string:
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			n.Int, n.Valid = int64(0), false
			return err
		}
		n.Int, n.Valid = int64(i), true
		return nil
	case []uint8:
		s := Uint8ArrayToString(v)
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			n.Int, n.Valid = int64(0), false
			return err
		}
		n.Int, n.Valid = i, true
		return nil
	default:
		n.Int, n.Valid = int64(0), false
		return NotConvertTypeError
	}
}

func (n *NullInt) Value() interface{} {
	if n.Valid == false {
		return nil
	}
	return n.Int
}

type NullFloat struct {
	Float float64
	Valid bool
}

func (n *NullFloat) Scan(v interface{}) error {
	if v == nil {
		n.Float, n.Valid = float64(0), false
		return nil
	}
	switch v := v.(type) {
	case int:
		n.Float, n.Valid = float64(v), true
		return nil
	case int64:
		n.Float, n.Valid = float64(v), true
		return nil
	case uint64:
		n.Float, n.Valid = float64(v), true
		return nil
	case float64:
		n.Float, n.Valid = v, true
		return nil
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			n.Float, n.Valid = float64(0), false
			return err
		}
		n.Float, n.Valid = f, true
		return nil
	case []uint8:
		s := Uint8ArrayToString(v)
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			n.Float, n.Valid = float64(0), false
			return err
		}
		n.Float, n.Valid = f, true
		return nil
	default:
		n.Float, n.Valid = float64(0), false
		return nil
	}
}

func (n *NullFloat) Value() interface{} {
	if n.Valid == false {
		return nil
	}
	return n.Float
}

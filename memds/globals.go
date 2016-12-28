package memds

import (
	"reflect"

	"github.com/ugorji/go/codec"
)

var (
	mh      codec.MsgpackHandle
	buckets Buckets
)

func init() {
	mh.MapType = reflect.TypeOf(map[string]interface{}(nil))
}

package memds

import (
	"fmt"

	"github.com/ugorji/go/codec"
)

func Exec(b []byte) []byte {
	cmd := make(map[string]interface{})
	dec := codec.NewDecoderBytes(b, &mh)
	if err := dec.Decode(&cmd); err != nil {
		return responseCmdDecodeError(err.Error())
	}

	c, ok := cmd["cmd"]
	if !ok {
		return responseCmdFormatError(fmt.Sprintf("key 'cmd' not found"))
	}
	var cs string
	switch v := c.(type) {
	case string:
		cs = v
	case []uint8:
		cs = Uint8ArrayToString(v)
	default:
		return responseCmdFormatError(fmt.Sprintf("key 'cmd' not type string"))
	}

	switch cs {
	case "get":
		k, ok := cmd["key"]
		if !ok {
			return responseCmdFormatError(fmt.Sprintf("key 'key' not found"))
		}
		var ks string
		switch v := k.(type) {
		case string:
			ks = v
		case []uint8:
			ks = Uint8ArrayToString(v)
		default:
			return responseCmdFormatError(fmt.Sprintf("key 'key' not type string"))
		}

		v, err := Get(ks)
		if err != nil && err != ValueNotFoundError {
			return responseCmdExecuteError(err.Error())
		}

		return response(map[string]interface{}{"value": v})
	case "set":
		k, ok := cmd["key"]
		if !ok {
			return responseCmdFormatError(fmt.Sprintf("key 'key' not found"))
		}
		var ks string
		switch v := k.(type) {
		case string:
			ks = v
		case []uint8:
			ks = Uint8ArrayToString(v)
		default:
			return responseCmdFormatError(fmt.Sprintf("key 'key' not type string"))
		}

		v, ok := cmd["value"]
		if !ok {
			return responseCmdFormatError(fmt.Sprintf("key 'value' not found"))
		}

		err := Set(ks, v)
		if err != nil {
			return responseCmdExecuteError(err.Error())
		}

		return responseOK()
	case "del":
		k, ok := cmd["key"]
		if !ok {
			return responseCmdFormatError(fmt.Sprintf("key 'key' not found"))
		}
		var ks string
		switch v := k.(type) {
		case string:
			ks = v
		case []uint8:
			ks = Uint8ArrayToString(v)
		default:
			return responseCmdFormatError(fmt.Sprintf("key 'key' not type string"))
		}

		err := Del(ks)
		if err != nil {
			return responseCmdExecuteError(err.Error())
		}
		return responseOK()
	default:
		return responseCmdNotFoundError()
	}
}

func Get(k string) (interface{}, error) {
	b := buckets.Get(k)
	if b == nil {
		return nil, BucketNotFoundError
	}
	return b.Get(k)
}

func Set(k string, v interface{}) error {
	b := buckets.Get(k)
	if b == nil {
		return BucketNotFoundError
	}
	return b.Set(k, v)
}

func Del(k string) error {
	b := buckets.Get(k)
	if b == nil {
		return BucketNotFoundError
	}
	b.Del(k)
	return nil
}

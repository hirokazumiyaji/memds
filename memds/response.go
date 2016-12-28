package memds

import (
	"fmt"

	"github.com/ugorji/go/codec"
)

func response(m map[string]interface{}) []byte {
	if _, ok := m["status"]; !ok {
		m["status"] = true
	}
	var rb []byte
	enc := codec.NewEncoderBytes(&rb, &mh)
	if err := enc.Encode(m); err != nil {
		Error(fmt.Sprintf("response encode error: %v", err))
	}
	rb = append(rb, '\n')
	return rb
}

func errorResponse(m map[string]interface{}) []byte {
	m["status"] = false
	return response(m)
}

func responseOK() []byte {
	return response(
		map[string]interface{}{
			"msg": "OK",
		},
	)
}

func responseCmdDecodeError(e string) []byte {
	return errorResponse(
		map[string]interface{}{
			"code": ErrorCodeCommandDecodeError,
			"msg":  e,
		},
	)
}

func responseCmdFormatError(e string) []byte {
	return errorResponse(
		map[string]interface{}{
			"code": ErrorCodeCommandFormatError,
			"msg":  e,
		},
	)
}

func responseCmdNotFoundError() []byte {
	return errorResponse(
		map[string]interface{}{
			"code": ErrorCodeCommandNotFoundError,
			"msg":  fmt.Sprintf("cmd format error: not found cmd"),
		},
	)
}

func responseCmdExecuteError(e string) []byte {
	return errorResponse(
		map[string]interface{}{
			"code": ErrorCodeCommandExecuteError,
			"msg":  e,
		},
	)
}

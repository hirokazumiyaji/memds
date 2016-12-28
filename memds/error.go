package memds

import "errors"

const (
	ErrorCodeCommandDecodeError   = 100
	ErrorCodeCommandFormatError   = 200
	ErrorCodeCommandNotFoundError = 300
	ErrorCodeCommandExecuteError  = 400
)

var (
	BucketsLEZeroError   = errors.New("bucket num can't le 0")
	BucketNotFoundError  = errors.New("bucket not found")
	ValueNotFoundError   = errors.New("value not found")
	CommandNotFoundError = errors.New("command not found")
)

package protocol

import "errors"

var (
	PayloadTooLargeError = errors.New("payload size exceeds maximum allowed size")
)

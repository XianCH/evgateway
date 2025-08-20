package protocol

import "errors"

var (
	PayloadTooLargeError = errors.New("payload size exceeds maximum allowed size")

	ErrNeedMoreData = errors.New("need more data to parse frame")

	ErrCRCMismatch = errors.New("crc mismatch")
)

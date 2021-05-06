package cos

import "errors"

var (
	// ErrServerSideEncryptionCustomerKey will be returned while server-side encryption customer key is invalid.
	ErrServerSideEncryptionCustomerKey = errors.New("invalid server-side encryption customer key")
)

// ErrCode
//
// ref: https://cloud.tencent.com/document/product/436/7730
const (
	// NoSuchKey the specified key does not exist.
	NoSuchKey = "NoSuchKey"
)

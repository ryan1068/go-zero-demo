package model

import "errors"

// ErrAbortTransaction signals a business rejection that must roll back the
// current transaction. Logic maps it to a business response.
var ErrAbortTransaction = errors.New("abort transaction")

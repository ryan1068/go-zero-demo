package model

import (
	"errors"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// ErrAbortTransaction signals a business rejection that must roll back the
// current transaction. Logic maps it to a business response.
var ErrAbortTransaction = errors.New("abort transaction")

// sessionOrConn returns the transaction session when one is provided,
// otherwise the model's raw connection. Both satisfy sqlx.Session.
func sessionOrConn(conn sqlx.SqlConn, session sqlx.Session) sqlx.Session {
	if session != nil {
		return session
	}
	return conn
}

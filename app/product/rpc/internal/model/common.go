package model

import "database/sql"

// NullString builds a valid nullable string for generated model fields.
func NullString(value string) sql.NullString {
	return sql.NullString{String: value, Valid: true}
}

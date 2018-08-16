package ftnlib

import (
	"database/sql"
)

// ContainsString - method to define is string value in array
func ContainsString(needle string, haystack []string) bool {
	for _, elem := range haystack {
		if elem == needle {
			return true
		}
	}
	return false
}

// CountSQLResultRows - count the sql rows
func CountSQLResultRows(rows *sql.Rows) (cnt int) {
	for rows.Next() {
		cnt++
	}
	return
}

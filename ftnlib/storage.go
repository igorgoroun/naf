package ftnlib

import (
	"database/sql"
	"fmt"
	"strings"
	// SQLite3 driver
	_ "github.com/mattn/go-sqlite3"
)

var (
	storeTypes = []string{
		"config",
		"nodelist",
		"echomail",
		"netmail",
	}
)

// GetStorage - returns db for querying
func GetStorage(path string, store string) (db *sql.DB, err error) {
	// Check requested storetype
	if !ContainsString(store, storeTypes) {
		return nil, fmt.Errorf("Invalid store type requested")
	}
	// Init DB File (separate for some data)
	filepath := strings.Join([]string{path, "naf.db"}, "")
	switch store {
	case "nodelist":
		filepath = strings.Join([]string{path, "nodelist.db"}, "")
	}

	db, err = sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}
	//defer db.Close()

	// Check storetype tables exists
	switch store {
	case "nodelist":
		rows, _ := db.Query(sqNodelistCheckTables)
		rowsCnt := CountSQLResultRows(rows)
		if rowsCnt != 2 {
			// create tables
			for _, q := range sqNodelistCreateTables {
				_, err := db.Exec(q)
				if err != nil {
					return nil, fmt.Errorf("Cannot create nodelist tables")
				}
			}
		}
	}

	return
}

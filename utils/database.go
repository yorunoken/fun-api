package utils

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type DatabaseData struct {
	Key   string
	Value string
}

func EntryExists(db *sql.DB, table string, id string) bool {
	var count int
	row := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE id = ?", table), id)
	err := row.Scan(&count)
	if err != nil {
		fmt.Println("Error scanning row:", err)
		return false
	}
	return count > 0
}

func GetEntry(db *sql.DB, table string, id string) *sql.Row {
	return db.QueryRow(fmt.Sprintf("SELECT * FROM %s WHERE id = ?", table), id)
}

func AddEntry(db *sql.DB, table string, id string, dataArr []DatabaseData) (sql.Result, error) {
	setClause := ""
	values := []string{}

	for index, data := range dataArr {
		values = append(values, data.Value)
		if index == len(dataArr)-1 {
			setClause += data.Key + " = ?"
		} else {
			setClause += data.Key + " = ?, "
		}
	}

	valueArgs := make([]interface{}, len(values))
	for i, v := range values {
		valueArgs[i] = v
	}

	rowExists := EntryExists(db, table, id)
	if !rowExists {
		fields := []string{}
		placeholders := ""

		for index, data := range dataArr {
			fields = append(fields, data.Key)

			if index == len(dataArr)-1 {
				placeholders += "?"
			} else {
				placeholders += "?, "
			}

		}

		args := append([]interface{}{id}, valueArgs...)
		query := fmt.Sprintf("INSERT INTO %s (id, %s) values (?, %s)", table, strings.Join(fields, ", "), placeholders)
		return db.Exec(query, args...)
	}

	args := append(valueArgs, id)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE ID = ?", table, setClause)
	return db.Exec(query, args)
}

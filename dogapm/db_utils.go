package dogapm

import (
	"database/sql"
)
type dbUtil struct{}
var DBUtil = &dbUtil{}

func (d *dbUtil) Query(rows *sql.Rows, err error) []map[string]any {
	if err != nil {
		panic(err)
	}
	if rows == nil {
		return []map[string]any{}
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}
	result := make([]map[string]any, 0)
	for rows.Next() {
		rowMap := make(map[string]any)
		columnPointers := make([]any, len(columns))
		for i := range columnPointers {
			columnPointers[i] = new(any)
		}
		if err := rows.Scan(columnPointers...); err != nil {
			panic(err)
		}
		for i, colName := range columns {
			val := *(columnPointers[i].(*any))
			rowMap[colName] = val
		}
		result = append(result, rowMap)
	}
	return result
}

func (d *dbUtil) QueryFirst(rows *sql.Rows, err error) map[string]any {
	if err != nil {
		panic(err)
	}
	results := d.Query(rows, nil)
	if len(results) == 0 {
		return map[string]any{}
	}
	return results[0]
}
package pageable

import (
	"database/sql"
	"fmt"
)

type Driver struct {
	Type string
}

func NewDriver(dbType string) (*Driver, error) {
	if dbType == "SQL" || dbType == "Mongo" {
		return &Driver{Type: dbType}, nil
	}
	return nil, fmt.Errorf("DB not supported, Please read Docs for more insights.")
}

func (d *Driver) Paginate(db *sql.DB, tableName string, page, size int, sort string) (*Response, error) {
	if d.Type == "SQL" {
		response, err := PaginateSQL(db, tableName, page, size, sort)
		if err != nil {
			return nil, err
		}
		return &response, nil
	}
	return nil, nil
}

func PaginateSQL(db *sql.DB, tableName string, page, size int, sort string) (Response, error) {
	response := Response{}
	var count int
	countStatement := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	errCount := db.QueryRow(countStatement).Scan(&count)
	if errCount != nil {
		return response, errCount
	}
	offset := (page - 1) * size
	sqlStatement := fmt.Sprintf("SELECT * FROM %s LIMIT %d OFFSET %d", tableName, size, offset)
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return response, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return response, err
	}
	var result []map[string]interface{}
	for rows.Next() {
		values := make(map[string]interface{})
		scanArgs := make([]interface{}, len(columns))
		for i := range columns {
			scanArgs[i] = new(interface{})
		}
		if err := rows.Scan(scanArgs...); err != nil {
			return response, err
		}
		for i, col := range columns {
			value := *(scanArgs[i].(*interface{}))
			values[col] = value
		}
		result = append(result, values)
	}
	if err := rows.Err(); err != nil {
		return response, err
	}
	pr := NewPaginatedResponse(result, len(result), int64(count), page, size, "")
	return *pr, nil
}

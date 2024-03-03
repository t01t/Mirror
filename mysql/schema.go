package mysql

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/t01t/mirror/helpers"
)

func (ms *MysqlServer) GetDatabases() error {
	ms.Databases = make(map[string]*DB)
	rows, err := ms.db.Query("SHOW DATABASES")
	if err != nil {
		return err
	}
	ms.IsConnected = true
	for rows.Next() {
		var dbname string
		err := rows.Scan(&dbname)
		if err != nil {
			log.Println("getting database name err:", err)
			continue
		}

		if !helpers.IsInArray(dbname, ms.Dbs) {
			continue
		}

		tables, err := ms.GetTables(dbname)
		if err != nil {
			log.Println("getting", dbname, "database tables err:", err)
			continue
		}
		ms.Databases[dbname] = &DB{
			Name:   dbname,
			Tables: tables,
		}
	}
	return nil
}

func (ms *MysqlServer) GetTables(db string) (map[string]Table, error) {
	_, err := ms.db.Exec("USE `" + db + "`")
	if err != nil {
		return nil, err
	}

	rows, err := ms.db.Query("SHOW FULL TABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tables := make(map[string]Table)
	for rows.Next() {
		var tableName, tableType string
		err := rows.Scan(&tableName, &tableType)
		if err != nil {
			return nil, err
		}
		columnNames, err := ms.GetColumns(db, tableName)
		if err != nil {
			return nil, err
		}
		var primarys []int
		for i, c := range columnNames {
			if c.Key == "PRI" {
				primarys = append(primarys, i)
			}
		}
		tables[tableName] = Table{
			Name:     tableName,
			Columns:  columnNames,
			Primarys: primarys,
			Type:     tableType,
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

func (ms *MysqlServer) GetColumns(db, table string) ([]Column, error) {
	_, err := ms.db.Exec("USE `" + db + "`")
	if err != nil {
		return nil, err
	}

	rows, err := ms.db.Query("SHOW COLUMNS FROM `" + table + "`")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []Column
	var primarysCount = 0
	for rows.Next() {
		var name, datatype, null, key, datadefault, extra sql.NullString
		err := rows.Scan(&name, &datatype, &null, &key, &datadefault, &extra)
		if err != nil {
			return nil, err
		}

		var isnull bool
		if null.String != "NO" {
			isnull = true
		}
		columns = append(columns, Column{
			Name:    name.String,
			Type:    datatype.String,
			Null:    isnull,
			Key:     key.String,
			Default: datadefault.String,
			Extra:   extra.String,
		})

		if key.String == "PRI" {
			primarysCount++
		}
	}
	// when the table does not have primary key
	if primarysCount == 0 {
		for i := range columns {
			columns[i].Key = "PRI"
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return columns, nil
}

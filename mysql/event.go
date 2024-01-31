package mysql

import (
	"fmt"

	"github.com/t01t/mirror/helpers"
)

const (
	INSERT = iota
	UPDATE
	DELETE
	QUERY
)

func operationType(e string) uint8 {
	var operation uint8
	switch e {
	case "WriteRowsEventV2", "WriteRowsEventV1":
		operation = INSERT
	case "UpdateRowsEventV2", "UpdateRowsEventV1":
		operation = UPDATE
	case "DeleteRowsEventV2", "DeleteRowsEventV1":
		operation = DELETE
	}

	return operation
}

func getPrimaryValues(rows [][]interface{}, primaryIndexs []int) [][]interface{} {
	var primarys [][]interface{}
	for _, r := range rows {
		rowPrimarys := []interface{}{}
		for j, c := range r {
			if helpers.IsInArray(j, primaryIndexs) {
				rowPrimarys = append(rowPrimarys, interfaceToString(c, true))
			}
		}
		primarys = append(primarys, rowPrimarys)
	}
	return primarys
}

func getUpdateChanges(rows [][]interface{}, primary []int) ([]map[int]interface{}, [][]interface{}) {

	rowsCount := len(rows)
	updates := []map[int]interface{}{}
	primaryValues := [][]interface{}{}
	for i := 0; i < rowsCount; i += 2 {
		rowPrimaryValues := []interface{}{}
		oldRow := rows[i]
		newRow := rows[i+1]
		changes := make(map[int]interface{})
		for j, v := range newRow {
			switch v.(type) {
			case []uint8:
				v = fmt.Sprintf("%s", v)
			}
			if helpers.IsInArray(j, primary) {
				rowPrimaryValues = append(rowPrimaryValues, interfaceToString(oldRow[j], true))
			}
			old := interfaceToString(oldRow[j])
			new := interfaceToString(newRow[j])

			if old == new {
				continue
			}
			changes[j] = v
		}
		updates = append(updates, changes)
		primaryValues = append(primaryValues, rowPrimaryValues)
	}
	return updates, primaryValues
}

package mysql

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/t01t/mirror/helpers"
)

func LogToSql(event []interface{}, database *DB) string {

	switch event[0].(float64) {
	case INSERT:
		tableType := database.Tables[event[1].(string)].Type
		return buildInsertStatment(event[1].(string), event[2].([]interface{}), tableType)
	case UPDATE:
		columns := database.Tables[event[1].(string)].Columns
		primarys := database.Tables[event[1].(string)].Primarys
		return buildUpdateStatment(event[1].(string), event[2].([]interface{}), event[3].([]interface{}), columns, primarys)
	case DELETE:
		columns := database.Tables[event[1].(string)].Columns
		primarys := database.Tables[event[1].(string)].Primarys
		return buildDeleteStatment(event[1].(string), event[2].([]interface{}), columns, primarys)
	case QUERY:
		return event[1].(string)
	default:
		return ""
	}
}

func interfaceToString(v interface{}, noQuotes ...bool) string {
	var str string
	if v != nil {
		switch val := v.(type) {
		case []byte, string:
			if len(noQuotes) > 0 {
				str = replacer.Replace(fmt.Sprintf("%s", val))
			} else {
				str = "'" + replacer.Replace(fmt.Sprintf("%s", val)) + "'"
			}
		case float64:
			str = strconv.FormatFloat(val, 'f', -1, 64)
		case float32:
			str = strconv.FormatFloat(float64(val), 'f', -1, 32)
		default:
			str = fmt.Sprint(val)
		}
	} else {
		str = "NULL"
	}

	return str
}

func buildDeleteStatment(table string, primarys []interface{}, columns []Column, primaryIndexList []int) string {
	var str strings.Builder
	str.WriteString("DELETE FROM `")
	str.WriteString(table)
	str.WriteString("` WHERE ")
	for i, p := range primaryIndexList {
		primaryKeyName := columns[p].Name
		if i > 0 {
			str.WriteString(" AND ")
		}

		tmpValues := []interface{}{}
		for _, p := range primarys {
			val := p.([]interface{})[i]
			if helpers.IsInArray(val, tmpValues) {
				continue
			}
			tmpValues = append(tmpValues, val)
		}

		if len(tmpValues) == 1 {
			if tmpValues[0] == "NULL" {
				str.WriteString("`")
				str.WriteString(primaryKeyName)
				str.WriteString("` IS NULL")
			} else {
				str.WriteString("`")
				str.WriteString(primaryKeyName)
				str.WriteString("`=")
				str.WriteString(interfaceToString(tmpValues[0]))
			}
		} else {
			hasNull := false
			j := 0
			var tmpStr strings.Builder
			tmpStr.WriteString("`")
			tmpStr.WriteString(primaryKeyName)
			tmpStr.WriteString("` IN (")
			for _, val := range tmpValues {
				if val == "NULL" {
					hasNull = true
				} else {
					if j > 0 {
						tmpStr.WriteString(",")
					}
					tmpStr.WriteString(interfaceToString(val))
					j++
				}
			}
			tmpStr.WriteString(")")
			if hasNull {
				tmpStr.WriteString(" OR `")
				tmpStr.WriteString(primaryKeyName)
				tmpStr.WriteString("` IS NULL")
				str.WriteString("(" + tmpStr.String() + ")")
			} else {
				str.WriteString(tmpStr.String())
			}
		}
	}
	str.WriteString(";")
	return str.String()
}

func buildInsertStatment(table string, rows []interface{}, tableType string) string {
	var str strings.Builder
	switch tableType {
	case "BASE TABLE":
		str.WriteString("INSERT INTO `")
		str.WriteString(table)
		str.WriteString("` VALUES ")
		for i, r := range rows {
			if i > 0 {
				str.WriteString(",")
			}
			str.WriteString("(")
			for j, c := range r.([]interface{}) {
				if j > 0 {
					str.WriteString(",")
				}
				str.WriteString(interfaceToString(c))
			}
			str.WriteString(")")
		}
	case "SEQUENCE":
		str.WriteString("ALTER SEQUENCE `")
		str.WriteString(table)
		str.WriteString("` RESTART WITH ")
		str.WriteString(interfaceToString(rows[0].([]interface{})[0]))
	}
	str.WriteString(";")
	return str.String()
}

func buildUpdateStatment(table string, rows []interface{}, primaryList []interface{}, columns []Column, primaryIndexList []int) string {
	var str strings.Builder
	for i, r := range rows {
		str.WriteString("UPDATE `")
		str.WriteString(table)
		str.WriteString("` SET ")
		row := r.(map[string]interface{})

		counter := 0
		for j, c := range row {
			if counter > 0 {
				str.WriteString(",")
			}
			j, _ := strconv.Atoi(j)
			new := interfaceToString(c)
			str.WriteString("`")
			str.WriteString(columns[j].Name)
			str.WriteString("`=")
			str.WriteString(new)

			counter++
		}

		str.WriteString(" WHERE ")
		for j, p := range primaryIndexList {
			if j > 0 {
				str.WriteString(" AND ")
			}
			if primaryList[i].([]interface{})[j] == "NULL" {
				str.WriteString("`")
				str.WriteString(columns[p].Name)
				str.WriteString("` IS NULL")
				continue
			}

			primaryValue := interfaceToString(primaryList[i].([]interface{})[j])
			str.WriteString("`")
			str.WriteString(columns[p].Name)
			str.WriteString("`=")
			str.WriteString(primaryValue)
		}
		str.WriteString(";")
	}
	return str.String()
}

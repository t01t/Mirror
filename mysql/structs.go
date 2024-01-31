package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlServer struct {
	Name        string
	Host        string
	Port        string
	User        string
	Pass        string `json:"-"`
	Databases   map[string]*DB
	db          *sql.DB
	Dbs         []string
	BinLogFile  string
	BinLogPos   uint32
	IsConnected bool
	State       string
}
type DB struct {
	Name   string
	Tables map[string]Table `json:"tables"`
}

type Table struct {
	Name     string   `json:"name"`
	Columns  []Column `json:"columns"`
	Primarys []int
}

type Column struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Null    bool   `json:"null"`
	Key     string `json:"key"`
	Default string `json:"default"`
	Extra   string `json:"extra"`
}

type Event struct {
	Server, Database, Table, Query string
	Type                           uint8
	Rows                           [][]interface{}
}

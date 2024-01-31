package mysql

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

func (ms *MysqlServer) Connect() error {
	connectionAttempts := 0
	for {
		connectionAttempts++
		if connectionAttempts > 60 {
			return errors.New("faild to conenct to server " + ms.Host)
		}

		conn, _ := sql.Open("mysql", ms.User+":"+ms.Pass+"@tcp("+ms.Host+":"+ms.Port+")/")
		err := conn.Ping()
		if err != nil {
			log.Println(err, "retrying after 10s")
			time.Sleep(10 * time.Second)
			continue
		}

		ms.db = conn
		log.Println("connected")
		return nil
	}
}

func (ms *MysqlServer) Close() {
	ms.db.Close()
}

func (ms *MysqlServer) Status() error {
	if ms.db == nil {
		return errors.New("no db is initilized")
	}
	path := "servers/" + ms.Name + "/status.yml"
	if _, err := os.Stat(path); err == nil {
		return ms.getMasterStatusFromCache()
	}
	return ms.getMasterStatusFromServer()
}

func (ms *MysqlServer) getMasterStatusFromServer() error {
	path := "servers/" + ms.Name + "/status.yml"
	rows, err := ms.db.Query("SHOW MASTER STATUS")

	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		cols, err := rows.Columns()
		if err != nil {
			return err
		}
		var file, position, binlogdodb, binlogignordb, executedgtid sql.NullString
		if len(cols) == 5 {
			err = rows.Scan(&file, &position, &binlogdodb, &binlogignordb, &executedgtid)
			if err != nil {
				return err
			}
		} else if len(cols) == 4 {
			err = rows.Scan(&file, &position, &binlogdodb, &binlogignordb)
			if err != nil {
				return err
			}
		} else {
			return errors.New("sql accespt params")
		}

		pos, err := strconv.Atoi(position.String)
		if err != nil {
			return err
		}
		ms.BinLogFile = file.String
		ms.BinLogPos = uint32(pos)
	}

	if ms.Name != "" {
		status := make(map[string]interface{})
		status["binlogname"] = ms.BinLogFile
		status["binlogpos"] = ms.BinLogPos
		status["time"] = time.Now().Format("2006-01-02 15:04:05")

		yml, err := yaml.Marshal(&status)
		if err != nil {
			return err
		}

		if _, err := os.Stat("servers/" + ms.Name); os.IsNotExist(err) {
			err := os.MkdirAll("servers/"+ms.Name, 0744)
			if err != nil {
				return err
			}
		}

		err2 := os.WriteFile(path, yml, 0744)
		if err2 != nil {
			return err2
		}
	}
	return nil
}

func (ms *MysqlServer) getMasterStatusFromCache() error {
	path := "servers/" + ms.Name + "/status.yml"
	state := make(map[string]string)
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, state)
	if err != nil {
		return err
	}
	pos, err := strconv.Atoi(state["binlogpos"])
	if err != nil {
		return err
	}
	ms.BinLogFile = state["binlogname"]
	ms.BinLogPos = uint32(pos)

	return nil
}

func (ms *MysqlServer) updatePosition(name string, pos uint32) error {
	path := "servers/" + ms.Name + "/status.yml"

	status := make(map[string]interface{})
	status["binlogname"] = name
	status["binlogpos"] = pos
	status["time"] = time.Now().Format("2006-01-02 15:04:05")

	yml, err := yaml.Marshal(&status)
	if err != nil {
		return err
	}

	if _, err := os.Stat("servers/" + ms.Name); os.IsNotExist(err) {
		err := os.MkdirAll("servers/"+ms.Name, 0744)
		if err != nil {
			return err
		}
	}

	err2 := os.WriteFile(path, yml, 0744)
	if err2 != nil {
		return err2
	}
	return nil
}

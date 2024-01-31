package mysql

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (ms *MysqlServer) EventToLog(e Event) []interface{} {
	// type | tableName | rows | primarys | time
	log := make([]interface{}, 0, 5)

	switch e.Type {
	case INSERT:
		// type int | tableName string | rows array<map<column index, value>>
		for i, r := range e.Rows {
			for j, c := range r {
				switch c.(type) {
				case []uint8:
					e.Rows[i][j] = fmt.Sprintf("%s", c)
				}
			}
		}
		log = append(log, e.Type, e.Table, e.Rows)

	case UPDATE:
		// type int | tableName string | rows array<map<column index, value>> | primarys array<array<values>>
		primarys := ms.Databases[e.Database].Tables[e.Table].Primarys
		u, p := getUpdateChanges(e.Rows, primarys)
		log = append(log, e.Type, e.Table, u, p)

	case DELETE:
		primarys := ms.Databases[e.Database].Tables[e.Table].Primarys
		// type int | tableName string | primarys array<array<values>>
		log = append(log, e.Type, e.Table, getPrimaryValues(e.Rows, primarys))

	case QUERY:
		if e.Query == "BEGIN" || e.Query == "COMMIT" {
			return log
		}
		if !strings.Contains(e.Query, ";") {
			e.Query += ";"
		}
		log = append(log, e.Type, e.Query)
	}
	return log
}

func addToLogFile(server, database string, event *[]interface{}) error {
	t := time.Now()
	*event = append(*event, t.Format("150405"))
	res, err := json.Marshal(event)
	if err != nil {
		return err
	}
	path := filepath.Join(AppPath, "servers", server, database, t.Format("2006-01-02")+".mrr")
	var file *os.File
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err = os.Create(path)
		if err != nil {
			return err
		}
		if err != nil {
			log.Println("Failed to write to file:", err)
			return err
		}
	} else {
		file, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("Failed to open file:", err)
			return err
		}
	}
	defer file.Close()

	_, err = file.Write(append(res, '\n'))
	if err != nil {
		log.Println("Failed to write to file:", err)
		return err
	}
	return nil
}

func LogFiles(server, database string) ([]map[string]interface{}, error) {
	var list []map[string]interface{}
	logsPath := filepath.Join(AppPath, "servers", server, database)
	files, err := os.ReadDir(logsPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if strings.Contains(file.Name(), ".mrr") {
			info, err := file.Info()
			if err != nil {
				return nil, err
			}
			details := make(map[string]interface{})
			details["name"] = strings.Replace(file.Name(), ".mrr", "", 1)
			details["size"] = info.Size()
			details["modification"] = info.ModTime().Format("2006-01-02 15:04:05")

			list = append(list, details)
		}
	}
	return list, nil
}

func ReadLogFile(server, database, path string) ([][]interface{}, error) {
	if path == "" {
		path = time.Now().Format("2006-01-02")
	}
	path = filepath.Join(AppPath, "servers", server, database, path+".mrr")
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var logs [][]interface{}
	l := 1
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		var logLine []interface{}
		err = json.Unmarshal(line, &logLine)
		if err != nil {
			log.Println("failed to parse log", path, "line", l)
			continue
		}
		var log = []interface{}{
			logLine[0],
			logLine[1],
			logLine[len(logLine)-1],
		}
		logs = append(logs, log)
		l++
	}

	return logs, nil
}

func ReadTableLogFileAsStream(server string, database string, table string, filename string, lineCh chan string) {
	if filename == "" {
		filename = time.Now().Format("2006-01-02")
	}
	path := filepath.Join(AppPath, "servers", server, database, filename+".mrr")
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	l := 1
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		if bytes.Contains(line, []byte("\""+table+"\"")) {
			lineCh <- string(line)
		}
		l++
	}
	lineCh <- ""
}

func LogFileToSqlAsStream(server string, database *DB, date string, lineCh chan string) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	path := filepath.Join(AppPath, "servers", server, database.Name, date+".mrr")
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	l := 1
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		var logLine []interface{}
		err = json.Unmarshal(line, &logLine)
		if err != nil {
			log.Println("failed to parse log", path, "line", l)
			continue
		}
		lineCh <- LogToSql(logLine, database)
		l++
	}
	lineCh <- ""
}

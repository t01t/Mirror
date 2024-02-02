package mysql

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/t01t/mirror/helpers"
	"github.com/t01t/mirror/sse"

	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
)

var (
	replacer *strings.Replacer
	AppPath  string
)

func (ms *MysqlServer) CDCinit() error {
	replacer = strings.NewReplacer("\\", "\\\\", "'", "\\'")
	err := ms.Connect()
	if err != nil {
		return err
	}
	err = ms.GetDatabases()
	if err != nil {
		log.Println("getting '", ms.Name, "' databases err: ", err)
		return err
	}

	err = ms.Status()
	if err != nil {
		log.Println("server status failed '", ms.Name, "' err:", err)
		return err
	}

	for _, db := range ms.Dbs {
		path := AppPath + "servers/" + ms.Name + "/" + db + "/" + db + ".sql"
		if _, err := os.Stat(path); err != nil {
			ms.Backup(db)
		}
	}

	err = ms.listen()
	if err != nil {
		return err
	}

	return nil
}

func (ms *MysqlServer) listen() error {
	port, err := strconv.Atoi(ms.Port)
	if err != nil {
		return err
	}
	cfg := replication.BinlogSyncerConfig{
		ServerID: 1001,
		Host:     ms.Host,
		Port:     uint16(port),
		User:     ms.User,
		Password: ms.Pass,
	}

RESTART:
	syncer := replication.NewBinlogSyncer(cfg)
	streamer, err := syncer.StartSync(mysql.Position{
		Name: ms.BinLogFile,
		Pos:  ms.BinLogPos,
	})
	if err != nil {
		fmt.Println("Failed to start binlog sync:", err)
		return err
	}
	today := time.Now().Format("2006-01-02")

	for {
		ev, err := streamer.GetEvent(context.Background())
		if err != nil {
			fmt.Println("Failed to get event from binlog:", err)
			if strings.Contains(err.Error(), "Could not find first log file name in binary log index file") {
				ms.getMasterStatusFromServer()
				syncer.Close()
				goto RESTART
			}
			return err
		}
		currentDate := time.Now().Format("2006-01-02")
		if currentDate != today {
			ms.DailyServerDBsBackup(today)
			today = currentDate
		}

		event := Event{
			Server: ms.Name,
		}

		switch ev.Event.(type) {
		case *replication.RowsEvent:
			e := ev.Event.(*replication.RowsEvent)
			event.Database = string(e.Table.Schema)
			event.Table = string(e.Table.Table)
			if !helpers.IsInArray(event.Database, ms.Dbs) {
				continue
			}
			event.Type = operationType(ev.Header.EventType.String())
			event.Rows = e.Rows
			fmt.Println("event:", ev.Header.EventType.String(), "database:", event.Database, "table:", event.Table)
		case *replication.QueryEvent:
			e := ev.Event.(*replication.QueryEvent)
			event.Database = string(e.Schema)
			event.Type = QUERY
			if !helpers.IsInArray(event.Database, ms.Dbs) {
				continue
			}
			event.Query = string(e.Query)
		default:
			continue
		}

		logs := ms.EventToLog(event)

		if len(logs) != 0 {
			err := addToLogFile(event.Server, event.Database, &logs)
			if err != nil {
				log.Println(err)
				continue
			}
			message := []interface{}{
				[]string{event.Server, event.Database, event.Table},
				logs,
			}
			sse.SendMessage(message)
		}

		// update status.yml with this position
		err = ms.updatePosition(syncer.GetNextPosition().Name, syncer.GetNextPosition().Pos)
		if err != nil {
			log.Println("updating log position faild", err)
		}
	}
}

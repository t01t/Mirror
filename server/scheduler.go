package server

import (
	"log"
	"os"
	"time"
)

func StartDailyBackupScheduler() {

	for {
		t := time.Now()
		if t.Format("15") == "05" {
			for _, server := range Servers {
				for _, db := range server.Dbs {
					path := "servers/" + server.Name + "/" + db + "/daily_fullbackup"
					if _, err := os.Stat(path); os.IsNotExist(err) {
						err := os.MkdirAll(path, 0744)
						if err != nil {
							log.Println(server.Name, db, "FULLBACKUP ERROR making dir", path, ":", err)
							continue
						}
					}

					path += "/" + t.Format("2006-01-02") + ".log"
					if _, err := os.Stat(path); err != nil {
						server.Backup(db, path)
					}
				}
			}
		}
		time.Sleep(time.Minute * 10)
	}
}

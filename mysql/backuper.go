package mysql

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func (ms *MysqlServer) Backup(dbName string, dir ...string) error {
	log.Println(ms.Name, dbName, "FULLBACKUP started ...")

	path := AppPath + "/servers/" + ms.Name + "/" + dbName //filepath.Join(currentDir, "/servers/"+ms.Name+"/"+dbName)
	var dumpFile string
	if len(dir) == 0 {
		dumpFile = filepath.Join(AppPath, path, dbName+".sql")
	} else {
		dumpFile = dir[0]
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0744)
		if err != nil {
			log.Println(ms.Name, dbName, "FULLBACKUP ERROR making dir", path, ":", err)
			return err
		}
	}

	cmd := exec.Command("mysqldump",
		"-h", ms.Host,
		"-P", ms.Port,
		"-u"+ms.User,
		"--databases", dbName,
		"--result-file="+dumpFile,
	)
	cmd.Env = append(os.Environ(), "MYSQL_PWD="+ms.Pass)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(ms.Name, dbName, "FULLBACKUP ERROR execute mysqldump command:", err)
		log.Println(string(output))
		return err
	}
	log.Println(ms.Name, dbName, "FULLBACKUP completed")

	return nil
}
func (ms *MysqlServer) DailyServerDBsBackup(date string) {
	for _, db := range ms.Dbs {
		path := "servers/" + ms.Name + "/" + db + "/daily_fullbackup"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			err := os.MkdirAll(path, 0744)
			if err != nil {
				log.Println(ms.Name, db, "FULLBACKUP ERROR making dir", path, ":", err)
				continue
			}
		}

		path += "/" + date + ".log"
		if _, err := os.Stat(path); err != nil {
			ms.Backup(db, path)
		}
	}
}

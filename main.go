package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/skratchdot/open-golang/open"
	"github.com/t01t/mirror/mysql"
	"github.com/t01t/mirror/server"
	"github.com/t01t/mirror/sse"
	"gopkg.in/yaml.v2"
)

func main() {
	loadServersFromYamlFile()

	var serversNames []string
	for _, server := range server.Servers {
		go server.CDCinit()
		serversNames = append(serversNames, server.Name)
	}

	go server.StartDailyBackupScheduler()
	sse.Init(serversNames)
	open.Run("http://127.0.0.1:2345")
	err := server.HttpServerStart()
	if err != nil {
		log.Fatal(err)
	}
}

func loadServersFromYamlFile() {
	server.Servers = make(map[string]*mysql.MysqlServer)
	exePath, err := os.Getwd()
	if err != nil {
		log.Println("Error:", err)
		return
	}
	exePath += "/Mirror"
	exeDir := filepath.Dir(exePath)

	var serversPath string

	if runtime.GOOS == "darwin" {
		//macOSDir := filepath.Join(exeDir, "..")
		//contentsDir := filepath.Join(macOSDir, "..")
		//appDir := filepath.Join(contentsDir, "..")
		mysql.AppPath = exeDir
		serversPath = "servers.yml" //filepath.Join(exeDir, "servers.yml")
	} else {
		mysql.AppPath = exeDir
		serversPath = filepath.Join(exeDir, "servers.yml")
	}

	yamlFile, err := os.ReadFile(serversPath)
	if err != nil {
		log.Fatal("can't find servers file")
	}
	err = yaml.Unmarshal(yamlFile, server.Servers)
	if err != nil {
		log.Fatal("settings file currepted")
	}

	for name, server := range server.Servers {
		server.Name = name
		if server.Host == "" {
			log.Fatal("server", name, "host is empty")
		}
		if server.Port == "" {
			log.Fatal("server", name, "port is empty")
		}
		if server.User == "" {
			log.Fatal("server", name, "user is empty")
		}
		if server.Pass == "" {
			log.Fatal("server", name, "password is empty")
		}
	}
}

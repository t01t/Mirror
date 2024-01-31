package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/t01t/mirror/frontend"
	"github.com/t01t/mirror/mysql"
	"github.com/t01t/mirror/sse"
)

var Servers map[string]*mysql.MysqlServer

func shutdown(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("done"))
	os.Exit(0)
}
func appPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("embeddedHTML").Parse(frontend.HomeHTML))
	data := map[string]interface{}{
		"logo": base64.StdEncoding.EncodeToString(frontend.LogoPNG),
	}
	// Execute the template and write the output to the response writer
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func cssHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	w.Write(frontend.StyleCSS)
}
func appJsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/js")
	w.Write(frontend.AppJS)
}
func bgJsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/js")
	w.Write(frontend.GradientJS)
}
func routeJsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/js")
	w.Write(frontend.RoutesJS)
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	sse.SSE.ServeHTTP(w, r)
}
func serversHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	response := make(map[string]*mysql.MysqlServer)
	for name, s := range Servers {
		res := *s
		res.Databases = nil
		res.Pass = "****"
		response[name] = &res
	}
	json, err := json.Marshal(response)
	if err != nil {
		return
	}
	fmt.Fprintln(w, string(json))
}
func serverHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	s := params["server"]
	server := Servers[s]

	json, err := json.Marshal(server)
	if err != nil {
		return
	}
	fmt.Fprintln(w, string(json))
}
func databaseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	server := params["server"]
	db := params["database"]
	tmpserver := *Servers[server]
	tmpserver.Pass = "****"

	res := make(map[string]interface{})
	res["server"] = tmpserver
	res["database"] = Servers[server].Databases[db]
	json, err := json.Marshal(res)
	if err != nil {
		return
	}
	fmt.Fprintln(w, string(json))
}
func tableHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	server := params["server"]
	db := params["database"]
	table := params["table"]

	tmpserver := *Servers[server]
	response := make(map[string]interface{})
	tmpserver.Pass = "****"
	tmpserver.Databases = nil
	response["server"] = tmpserver
	response["database"] = Servers[server].Databases[db]
	response["table"] = Servers[server].Databases[db].Tables[table]

	json, err := json.Marshal(response)
	if err != nil {
		return
	}
	fmt.Fprintln(w, string(json))
}
func logHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	server := params["server"]
	db := params["database"]
	date := params["date"]

	logs, err := mysql.ReadLogFile(server, db, date)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintln(w, err)
		return
	}
	res, err := json.Marshal(logs)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, err)
		return
	}
	w.Write(res)
}

func logFilesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	server := params["server"]
	db := params["database"]

	files, err := mysql.LogFiles(server, db)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, err)
		return
	}
	res, err := json.Marshal(files)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, err)
		return
	}
	w.Write(res)
}

func sqlStreamHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	server := params["server"]
	db := params["database"]
	date := params["date"]
	database := Servers[server].Databases[db]
	logCh := make(chan string)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/x-sql; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename="+server+"."+db+"_"+date+".sql")
	w.Header().Set("Content-Type", "application/octet-stream")
	defer close(logCh)
	fmt.Fprintln(w, "USE "+database.Name+";")
	go mysql.LogFileToSqlAsStream(server, database, date, logCh)
	var querys strings.Builder
	var i = 0
	for query := range logCh {
		if query == "" {
			break
		}
		querys.WriteString(query + "\n")
		if i > 100 {
			fmt.Fprintln(w, querys.String())
			i = 0
			querys.Reset()
		}
		i++
	}
	fmt.Fprintln(w, querys.String())
	w.(http.Flusher).Flush()
}

func tableLogshandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	server := params["server"]
	db := params["database"]
	table := params["table"]
	date := params["date"]
	logCh := make(chan string)
	defer close(logCh)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var logs strings.Builder
	logs.WriteString("[")
	go mysql.ReadTableLogFileAsStream(server, db, table, date, logCh)
	var i = 0
	for line := range logCh {
		if line == "" {
			break
		}
		if i > 0 {
			logs.WriteString(",")
		}
		logs.WriteString(line)
		i++
	}
	logs.WriteString("]")
	fmt.Fprint(w, logs.String())
	w.(http.Flusher).Flush()
}

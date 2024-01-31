package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func HttpServerStart() error {
	mux := mux.NewRouter()
	mux.HandleFunc("/", appPageHandler)
	// mux.Handle("/", http.FileServer(http.Dir("./frontend")))

	mux.HandleFunc("/assets/style.css", cssHandler)
	mux.HandleFunc("/assets/app.js", appJsHandler)
	mux.HandleFunc("/assets/route.js", routeJsHandler)
	mux.HandleFunc("/assets/bg.js", bgJsHandler)

	mux.HandleFunc("/api/servers", serversHandler)
	mux.HandleFunc("/api/servers/{server}", serverHandler)
	mux.HandleFunc("/api/servers/{server}/databases/{database}", databaseHandler)
	mux.HandleFunc("/api/servers/{server}/databases/{database}/files", logFilesHandler)
	mux.HandleFunc("/api/servers/{server}/databases/{database}/logs", logHandler)
	mux.HandleFunc("/api/servers/{server}/databases/{database}/logs/{date}", logHandler)
	mux.HandleFunc("/api/servers/{server}/databases/{database}/sql", sqlStreamHandler)
	mux.HandleFunc("/api/servers/{server}/databases/{database}/sql/{date}", sqlStreamHandler)
	mux.HandleFunc("/api/servers/{server}/databases/{database}/tables/{table}", tableHandler)
	mux.HandleFunc("/api/servers/{server}/databases/{database}/tables/{table}/logs", tableLogshandler)
	mux.HandleFunc("/api/servers/{server}/databases/{database}/tables/{table}/logs/{date}", tableLogshandler)

	mux.HandleFunc("/api/app/shutdown", shutdown)

	mux.HandleFunc("/events", sseHandler)
	return http.ListenAndServe(":2345", mux)
}

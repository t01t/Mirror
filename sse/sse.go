package sse

import (
	"encoding/json"

	"github.com/r3labs/sse/v2"
)

var SSE *sse.Server

func SendMessage(payload interface{}) {
	message, err := json.Marshal(payload)
	if err != nil {
		return
	}
	SSE.Publish("log", &sse.Event{
		Data: message,
	})
}

func Init(servers []string) {
	SSE = sse.New()
	SSE.BufferSize = 1
	SSE.CreateStream("log")
}

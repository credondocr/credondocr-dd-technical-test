package server

import (
	"credondocr-dd-technical-test/config"
	"fmt"
)

func Init() {
	config := config.GetConfig()
	r := NewRouter()
	r.Run(fmt.Sprintf(":%s", config.GetString("server.port")))
}

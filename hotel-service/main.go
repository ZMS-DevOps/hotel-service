package main

import (
	"github.com/ZMS-DevOps/hotel-service/startup"
	cfg "github.com/ZMS-DevOps/hotel-service/startup/config"
)

func main() {
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}

package main

import (
	"github.com/mmmajder/zms-devops-hotel-service/startup"
	cfg "github.com/mmmajder/zms-devops-hotel-service/startup/config"
)

func main() {
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}

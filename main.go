package main

import (
	"log"
	"strings"

	"github.com/charmbracelet/soft/config"
	"github.com/charmbracelet/soft/server"
)

func main() {
	cfg := config.DefaultConfig()
	s := server.NewServer(cfg)
	log.Printf("Starting SSH server on %s:%d\n", cfg.Host, cfg.SSHPort)
	log.Printf("Starting %s server on %s:%d\n", strings.ToUpper(cfg.HTTPScheme), cfg.Host, cfg.HTTPPort)
	err := s.Start()
	if err != nil {
		log.Fatalln(err)
	}
}

package server

import (
	"log"

	"github.com/charmbracelet/soft/config"
	appCfg "github.com/charmbracelet/soft/internal/config"
)

type Server struct {
	HTTPServer *HTTPServer
	SSHServer  *SSHServer
	Config     *config.Config
	ac         *appCfg.Config
}

func NewServer(cfg *config.Config) *Server {
	ac, err := appCfg.NewConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}
	return &Server{
		HTTPServer: NewHTTPServer(cfg, ac),
		SSHServer:  NewSSHServer(cfg, ac),
		Config:     cfg,
		ac:         ac,
	}
}

func (srv *Server) Reload() error {
	return srv.ac.Reload()
}

func (srv *Server) Start() error {
	go func() {
		if err := srv.HTTPServer.Start(); err != nil {
			log.Fatal(err)
		}
	}()
	return srv.SSHServer.Start()
}

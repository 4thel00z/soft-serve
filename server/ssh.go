package server

import (
	"fmt"
	"log"

	"github.com/charmbracelet/soft/config"
	appCfg "github.com/charmbracelet/soft/internal/config"
	"github.com/charmbracelet/soft/internal/tui"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	gm "github.com/charmbracelet/wish/git"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/gliderlabs/ssh"
)

type SSHServer struct {
	s   *ssh.Server
	cfg *config.Config
	ac  *appCfg.Config
}

// NewSSHServer returns a new *ssh.Server configured to serve Soft Serve. The SSH
// server key-pair will be created if none exists. An initial admin SSH public
// key can be provided with authKey. If authKey is provided, access will be
// restricted to that key. If authKey is not provided, the server will be
// publicly writable until configured otherwise by cloning the `config` repo.
func NewSSHServer(cfg *config.Config, ac *appCfg.Config) *SSHServer {
	mw := []wish.Middleware{
		bm.Middleware(tui.SessionHandler(ac)),
		gm.Middleware(cfg.RepoPath, ac),
		lm.Middleware(),
	}
	s, err := wish.NewServer(
		ssh.PublicKeyAuth(ac.PublicKeyHandler),
		ssh.PasswordAuth(ac.PasswordHandler),
		wish.WithAddress(fmt.Sprintf("%s:%d", cfg.Host, cfg.SSHPort)),
		wish.WithHostKeyPath(cfg.KeyPath),
		wish.WithMiddleware(mw...),
	)
	if err != nil {
		log.Fatalln(err)
	}
	return &SSHServer{
		s:   s,
		cfg: cfg,
		ac:  ac,
	}
}

func (s *SSHServer) Start() error {
	return s.s.ListenAndServe()
}

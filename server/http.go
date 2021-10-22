package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/charmbracelet/soft/config"
	appCfg "github.com/charmbracelet/soft/internal/config"
	"github.com/charmbracelet/wish/git"
	"goji.io"
	"goji.io/pat"
	"goji.io/pattern"
)

type HTTPServer struct {
	server     *http.Server
	gitHandler http.Handler
	cfg        *config.Config
	ac         *appCfg.Config
}

func NewHTTPServer(cfg *config.Config, ac *appCfg.Config) *HTTPServer {
	h := goji.NewMux()
	s := &HTTPServer{
		cfg:        cfg,
		ac:         ac,
		gitHandler: http.FileServer(http.Dir(cfg.RepoPath)),
		server: &http.Server{
			Addr:      fmt.Sprintf(":%d", cfg.HTTPPort),
			Handler:   h,
			TLSConfig: cfg.TLSConfig,
		},
	}
	h.HandleFunc(pat.Get("/:repo"), s.handleGit)
	h.HandleFunc(pat.Get("/:repo/*"), s.handleGit)
	return s
}

func (s *HTTPServer) Start() error {
	if s.cfg.HTTPScheme == "https" {
		return s.server.ListenAndServeTLS("", "")
	} else {
		return s.server.ListenAndServe()
	}
}

func (s *HTTPServer) handleGit(w http.ResponseWriter, r *http.Request) {
	repo := pat.Param(r, "repo")
	access := s.ac.AuthRepo(repo, nil)
	if access < git.ReadOnlyAccess || !s.ac.AllowKeyless {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	path := pattern.Path(r.Context())
	stat, err := os.Stat(filepath.Join(s.cfg.RepoPath, repo, path))
	// Restrict access to files
	if err != nil || stat.IsDir() {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	r.URL.Path = fmt.Sprintf("/%s/%s", repo, path)
	s.gitHandler.ServeHTTP(w, r)
}

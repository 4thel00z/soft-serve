package config

import (
	"log"
	"strings"

	gm "github.com/charmbracelet/wish/git"
	"github.com/gliderlabs/ssh"
)

// Push registers Git push functionality for the given repo and key.
func (cfg *Config) Push(repo string, pk ssh.PublicKey) {
	err := cfg.Reload()
	if err != nil {
		log.Printf("error reloading after push: %s", err)
	}
	if cfg.Cfg.Callbacks != nil {
		cfg.Cfg.Callbacks.Push(repo)
	}
}

// Fetch registers Git fetch functionality for the given repo and key.
func (cfg *Config) Fetch(repo string, pk ssh.PublicKey) {
	if cfg.Cfg.Callbacks != nil {
		cfg.Cfg.Callbacks.Fetch(repo)
	}
}

// AuthRepo grants repo authorization to the given key.
func (cfg *Config) AuthRepo(repo string, pk ssh.PublicKey) gm.AccessLevel {
	return cfg.accessForKey(repo, pk)
}

// PasswordHandler returns whether or not password access is allowed.
func (cfg *Config) PasswordHandler(ctx ssh.Context, password string) bool {
	return (cfg.AnonAccess != "no-access") && cfg.AllowKeyless
}

// PublicKeyHandler returns whether or not the given public key may access the
// repo.
func (cfg *Config) PublicKeyHandler(ctx ssh.Context, pk ssh.PublicKey) bool {
	return cfg.accessForKey("", pk) != gm.NoAccess
}

func (cfg *Config) accessForKey(repo string, pk ssh.PublicKey) gm.AccessLevel {
	private := cfg.isPrivate(repo)
	if repo == "config" {
		private = true
	}
	for _, u := range cfg.Users {
		for _, k := range u.PublicKeys {
			apk, _, _, _, err := ssh.ParseAuthorizedKey([]byte(strings.TrimSpace(k)))
			if err != nil {
				log.Printf("error: malformed authorized key: '%s'", k)
				return gm.NoAccess
			}
			if ssh.KeysEqual(pk, apk) {
				if u.Admin {
					return gm.AdminAccess
				}
				for _, r := range u.CollabRepos {
					if repo == r {
						return gm.ReadWriteAccess
					}
				}
				if !private {
					return gm.ReadOnlyAccess
				}
			}
		}
	}
	if private && (cfg.AnonAccess != "read-write") {
		return gm.NoAccess
	}
	switch cfg.AnonAccess {
	case "no-access":
		return gm.NoAccess
	case "read-only":
		return gm.ReadOnlyAccess
	case "read-write":
		return gm.ReadWriteAccess
	default:
		return gm.NoAccess
	}
}

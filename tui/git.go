package tui

import (
	"sort"

	gitypes "github.com/charmbracelet/soft-serve/tui/bubbles/git/types"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Repo struct {
	name   string
	repo   *git.Repository
	readme string
	ref    *plumbing.Reference
}

func (r *Repo) Name() string {
	return r.name
}

func (r *Repo) GetRef() *plumbing.Reference {
	return r.ref
}

func (r *Repo) SetRef(ref *plumbing.Reference) {
	r.ref = ref
}

func (r *Repo) Repository() *git.Repository {
	return r.repo
}

func (r *Repo) GetCommits(limit int) gitypes.Commits {
	commits := gitypes.Commits{}
	ci, err := r.repo.CommitObjects()
	if err != nil {
		return nil
	}
	err = ci.ForEach(func(c *object.Commit) error {
		commits = append(commits, &gitypes.Commit{c})
		return nil
	})
	if err != nil {
		return nil
	}
	sort.Sort(commits)
	if limit <= 0 || limit > len(commits) {
		limit = len(commits)
	}
	return commits[:limit]
}

func (r *Repo) GetReadme() string {
	if r.readme != "" {
		return r.readme
	}
	md, err := r.readFile("README.md")
	if err != nil {
		return ""
	}
	return md
}

func (r *Repo) readFile(path string) (string, error) {
	lg, err := r.repo.Log(&git.LogOptions{
		From: r.ref.Hash(),
	})
	if err != nil {
		return "", err
	}
	c, err := lg.Next()
	if err != nil {
		return "", err
	}
	f, err := c.File(path)
	if err != nil {
		return "", err
	}
	content, err := f.Contents()
	if err != nil {
		return "", err
	}
	return content, nil
}

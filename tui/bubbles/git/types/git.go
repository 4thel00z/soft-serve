package types

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type ReferenceName plumbing.ReferenceName

func (n ReferenceName) String() string {
	return plumbing.ReferenceName(n).String()
}

func (n ReferenceName) Short() string {
	return plumbing.ReferenceName(n).Short()
}

type Repo interface {
	Name() string
	GetReference() *plumbing.Reference
	SetReference(*plumbing.Reference)
	GetReadme() string
	GetCommits(limit int) Commits
	Repository() *git.Repository
}

type Commit struct {
	*object.Commit
}

type Commits []*Commit

func (cl Commits) Len() int      { return len(cl) }
func (cl Commits) Swap(i, j int) { cl[i], cl[j] = cl[j], cl[i] }
func (cl Commits) Less(i, j int) bool {
	return cl[i].Author.When.After(cl[j].Author.When)
}

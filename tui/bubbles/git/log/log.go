package log

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	gansi "github.com/charmbracelet/glamour/ansi"
	"github.com/charmbracelet/soft-serve/internal/git"
	"github.com/charmbracelet/soft-serve/tui/bubbles/git/style"
	"github.com/charmbracelet/soft-serve/tui/bubbles/git/types"
	vp "github.com/charmbracelet/soft-serve/tui/bubbles/git/viewport"
	"github.com/muesli/termenv"
)

const glamourMaxWidth = 120

var (
	diffChroma = &gansi.CodeBlockElement{
		Code:     "",
		Language: "diff",
	}
)

type pageView int

const (
	logView pageView = iota
	commitView
)

type item git.RepoCommit

func (i item) Title() string {
	lines := strings.Split(i.Commit.Message, "\n")
	if len(lines) > 0 {
		return lines[0]
	}
	return ""
}

func (i item) FilterValue() string { return i.Title() }

type itemDelegate struct {
	style *style.Styles
}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	if index == m.Index() {
		fmt.Fprint(w, d.style.LogItemSelector.MarginLeft(1).Render(">")+
			d.style.LogItemHash.MarginLeft(1).Bold(true).Render(i.Commit.Hash.String()[:7])+
			d.style.LogItemActive.MarginLeft(1).Render(i.Title()))
	} else {
		fmt.Fprint(w, d.style.LogItemSelector.MarginLeft(1).Render(" ")+
			d.style.LogItemHash.MarginLeft(1).Render(i.Commit.Hash.String()[:7])+
			d.style.LogItemInactive.MarginLeft(1).Render(i.Title()))
	}
}

type Bubble struct {
	repo           *git.Repo
	list           list.Model
	pageView       pageView
	commitViewport *vp.ViewportBubble
	style          *style.Styles
	width          int
	widthMargin    int
	height         int
	heightMargin   int
	rctx           gansi.RenderContext
}

// TODO enable filter
func NewBubble(repo *git.Repo, style *style.Styles, width, widthMargin, height, heightMargin int) *Bubble {
	items := make([]list.Item, 0)
	for _, c := range repo.GetCommits(0) {
		items = append(items, item(c))
	}
	l := list.NewModel(items, itemDelegate{style}, width-widthMargin, height-heightMargin)
	l.SetShowFilter(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	styles := "light"
	if termenv.HasDarkBackground() {
		styles = "dark"
	}
	b := &Bubble{
		commitViewport: &vp.ViewportBubble{
			Viewport: &viewport.Model{},
		},
		repo:         repo,
		list:         l,
		style:        style,
		pageView:     logView,
		width:        width,
		widthMargin:  widthMargin,
		height:       height,
		heightMargin: heightMargin,
		rctx: gansi.NewRenderContext(gansi.Options{
			ColorProfile: termenv.TrueColor,
			Styles:       *glamour.DefaultStyles[styles],
		}),
	}
	b.SetSize(width, height)
	return b
}

func (b *Bubble) Help() []types.HelpEntry {
	switch b.pageView {
	case logView:
		return []types.HelpEntry{
			{"enter", "select"},
		}
	case commitView:
		return []types.HelpEntry{
			{"esc", "back"},
		}
	default:
		return []types.HelpEntry{}
	}
}

func (b *Bubble) GotoTop() {
	b.commitViewport.Viewport.GotoTop()
}

func (b *Bubble) Init() tea.Cmd {
	return nil
}

func (b *Bubble) SetSize(width, height int) {
	b.width = width
	b.height = height
	b.commitViewport.Viewport.Width = width - b.widthMargin
	b.commitViewport.Viewport.Height = height - b.heightMargin
	b.list.SetSize(width-b.widthMargin, height-b.heightMargin)
}

func (b *Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch msg.String() {
		case "down", "j":
			if b.pageView == logView {
				b.list.CursorDown()
			}
		case "up", "k":
			if b.pageView == logView {
				b.list.CursorUp()
			}
		case "enter":
			if b.pageView == logView {
				b.pageView = commitView
				b.commitViewport.Viewport.SetContent(b.commitView())
				b.GotoTop()
			}
		case "esc":
			if b.pageView == commitView {
				b.pageView = logView
			}
		}
	}
	rv, cmd := b.commitViewport.Update(msg)
	b.commitViewport = rv.(*vp.ViewportBubble)
	cmds = append(cmds, cmd)
	return b, tea.Batch(cmds...)
}

func (b *Bubble) commitView() string {
	s := strings.Builder{}
	commits := b.repo.GetCommits(0)
	commit := commits[b.list.Index()]
	s.WriteString(fmt.Sprintf("%s\n%s\n%s\n%s\n",
		b.style.LogCommitHash.Render("commit "+commit.Commit.Hash.String()),
		b.style.LogCommitAuthor.Render("Author: "+commit.Commit.Author.String()),
		b.style.LogCommitDate.Render("Date:   "+commit.Commit.Committer.When.Format(time.UnixDate)),
		b.style.LogCommitBody.Render(strings.TrimSpace(commit.Commit.Message)),
	))
	stats, err := commit.Commit.Stats()
	if err == nil {
		s.WriteString(fmt.Sprintf("\n%s", b.renderStats(stats.String())))
	}
	if commit.Commit.NumParents() > 0 {
		parent, err := commit.Commit.Parent(0)
		if err == nil {
			ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
			defer cancel()
			patch, err := commit.Commit.PatchContext(ctx, parent)
			if err == nil {
				diffChroma.Code = patch.String()
				p := strings.Builder{}
				err := diffChroma.Render(&p, b.rctx)
				if err == nil {
					s.WriteString(fmt.Sprintf("\n%s", p.String()))
				}
			}
		}
	}
	w := b.width - b.widthMargin
	if w > glamourMaxWidth {
		w = glamourMaxWidth
	}
	return b.style.LogCommit.MaxWidth(w).Render(s.String())
}

func (b *Bubble) renderStats(stats string) string {
	rv := strings.Builder{}
	s := bufio.NewScanner(strings.NewReader(stats))
	for s.Scan() {
		line := s.Text()
		for _, c := range line {
			if c == '+' {
				rv.WriteString(b.style.LogCommitStatsAdd.Render(string(c)))
			} else if c == '-' {
				rv.WriteString(b.style.LogCommitStatsDel.Render(string(c)))
			} else {
				rv.WriteString(string(c))
			}
		}
		rv.WriteString("\n")
	}
	return rv.String()
}

func (b *Bubble) View() string {
	switch b.pageView {
	case logView:
		return b.list.View()
	case commitView:
		return b.commitViewport.View()
	default:
		return ""
	}
}

package git

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/soft-serve/internal/git"
	"github.com/charmbracelet/soft-serve/tui/bubbles/git/about"
	"github.com/charmbracelet/soft-serve/tui/bubbles/git/log"
	"github.com/charmbracelet/soft-serve/tui/bubbles/git/style"
	"github.com/charmbracelet/soft-serve/tui/bubbles/git/types"
	"github.com/muesli/reflow/truncate"
	"github.com/muesli/reflow/wrap"
)

const (
	repoNameMaxWidth = 32
)

type pageState int

const (
	aboutPage pageState = iota
	refsPage
	logPage
	treePage
)

func pageString(p pageState) string {
	switch p {
	case aboutPage:
		return "About"
	case refsPage:
		return "Refs"
	case logPage:
		return "Log"
	case treePage:
		return "Tree"
	default:
		return "Unknown"
	}
}

type Bubble struct {
	name         string
	host         string
	port         int
	page         pageState
	repoSource   *git.RepoSource
	height       int
	heightMargin int
	width        int
	widthMargin  int
	style        *style.Styles
	boxes        []tea.Model
	Active       bool
}

func NewBubble(host string, port int, name string, rs *git.RepoSource, styles *style.Styles, width, wm, height, hm int) *Bubble {
	b := &Bubble{
		repoSource:   rs,
		page:         aboutPage,
		width:        width,
		widthMargin:  wm,
		height:       height,
		heightMargin: hm,
		style:        styles,
		host:         host,
		port:         port,
		name:         name,
		boxes:        make([]tea.Model, 4),
	}
	repo, err := rs.GetRepo(name)
	if err != nil {
		return nil
	}
	heightMargin := hm + lipgloss.Height(b.headerView())
	b.boxes[aboutPage] = about.NewBubble(repo, b.style, b.width, wm, b.height, heightMargin)
	b.boxes[logPage] = log.NewBubble(repo, b.style, width, wm, height, heightMargin)
	return b
}

func (b *Bubble) Init() tea.Cmd {
	return b.setupCmd
}

func (b *Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "A":
			b.page = aboutPage
		case "L":
			b.page = logPage
		}
	case tea.WindowSizeMsg:
		b.width = msg.Width
		b.height = msg.Height
	}
	m, cmd := b.boxes[b.page].Update(msg)
	b.boxes[b.page] = m
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	b.boxes[b.page] = m
	return b, tea.Batch(cmds...)
}

func (b *Bubble) Help() []types.HelpEntry {
	h := []types.HelpEntry{}
	if b.page != aboutPage {
		h = append(h, types.HelpEntry{"A", "about"})
	}
	if b.page != logPage {
		h = append(h, types.HelpEntry{"L", "log"})
	}
	switch b.page {
	case aboutPage:
		h = append(h, types.HelpEntry{"f/b", "pgup/pgdown"})
	case logPage:
		h = append(h, b.boxes[logPage].(*log.Bubble).Help()...)
	}
	return h
}

func (b *Bubble) Styles() *style.Styles {
	return b.style
}

func (b *Bubble) headerView() string {
	// TODO better header, tabs?
	// Render repo title
	title := b.name
	if title == "config" {
		title = "Home"
	}
	title = truncate.StringWithTail(title, repoNameMaxWidth, "…")
	title = b.style.RepoTitle.Render(title)

	// Render clone command
	var note string
	if b.name == "config" {
		note = ""
	} else {
		note = fmt.Sprintf("git clone %s", b.sshAddress())
	}
	noteWidth := b.width -
		b.widthMargin -
		lipgloss.Width(title) -
		b.style.RepoTitleBox.GetHorizontalFrameSize()
	// Hard-wrap the clone command only, without the usual word-wrapping. since
	// a long repo name isn't going to be a series of space-separated "words",
	// we'll always want it to be perfectly hard-wrapped.
	note = wrap.String(note, noteWidth-b.style.RepoNote.GetHorizontalFrameSize())
	note = b.style.RepoNote.Copy().Width(noteWidth).Render(note)

	// Render borders on name and command
	height := max(lipgloss.Height(title), lipgloss.Height(note))
	titleBoxStyle := b.style.RepoTitleBox.Copy().Height(height)
	noteBoxStyle := b.style.RepoNoteBox.Copy().Height(height)
	if b.Active {
		titleBoxStyle = titleBoxStyle.BorderForeground(b.style.ActiveBorderColor)
		noteBoxStyle = noteBoxStyle.BorderForeground(b.style.ActiveBorderColor)
	}
	title = titleBoxStyle.Render(title)
	note = noteBoxStyle.Render(note)

	// Render
	return lipgloss.JoinHorizontal(lipgloss.Top, title, note)
}

func (b *Bubble) View() string {
	header := b.headerView()
	bs := b.style.RepoBody.Copy()
	if b.Active {
		bs = bs.BorderForeground(b.style.ActiveBorderColor)
	}
	body := bs.Width(b.width - b.widthMargin - b.style.RepoBody.GetVerticalFrameSize()).
		Height(b.height - b.heightMargin - lipgloss.Height(header)).
		Render(b.boxes[b.page].View())
	return header + body
}

func (b *Bubble) setupCmd() tea.Msg {
	cmds := make([]tea.Cmd, 0)
	for _, bx := range b.boxes {
		if bx != nil {
			initCmd := bx.Init()
			if initCmd != nil {
				msg := initCmd()
				switch msg := msg.(type) {
				case types.ErrMsg:
					return msg
				}
			}
			cmds = append(cmds, initCmd)
		}
	}
	return tea.Batch(cmds...)
}

func (b Bubble) sshAddress() string {
	p := ":" + strconv.Itoa(int(b.port))
	if p == ":22" {
		p = ""
	}
	return fmt.Sprintf("ssh://%s%s/%s", b.host, p, b.name)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

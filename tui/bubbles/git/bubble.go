package git

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/soft-serve/tui/bubbles/git/about"
	"github.com/charmbracelet/soft-serve/tui/bubbles/git/log"
	"github.com/charmbracelet/soft-serve/tui/bubbles/git/style"
	"github.com/charmbracelet/soft-serve/tui/bubbles/git/types"
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
	page         pageState
	repo         types.Repo
	height       int
	heightMargin int
	width        int
	widthMargin  int
	style        *style.Styles
	boxes        []tea.Model
}

func NewBubble(r types.Repo, styles *style.Styles, width, wm, height, hm int) *Bubble {
	b := &Bubble{
		repo:         r,
		page:         aboutPage,
		width:        width,
		widthMargin:  wm,
		height:       height,
		heightMargin: hm,
		style:        styles,
		boxes:        make([]tea.Model, 4),
	}
	heightMargin := hm + lipgloss.Height(b.headerView())
	b.boxes[aboutPage] = about.NewBubble(r, b.style, b.width, wm, b.height, heightMargin)
	b.boxes[logPage] = log.NewBubble(r, b.style, width, wm, height, heightMargin)
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
	h = append(h, b.boxes[logPage].(types.HelpableBubble).Help()...)
	return h
}

func (b *Bubble) Styles() *style.Styles {
	return b.style
}

func (b *Bubble) headerView() string {
	// TODO better header, tabs?
	return ""
}

func (b *Bubble) View() string {
	header := b.headerView()
	return header + b.boxes[b.page].View()
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

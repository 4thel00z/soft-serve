package about

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/soft-serve/internal/git"
	"github.com/charmbracelet/soft-serve/tui/bubbles/git/style"
	"github.com/charmbracelet/soft-serve/tui/bubbles/git/types"
	vp "github.com/charmbracelet/soft-serve/tui/bubbles/git/viewport"
)

const (
	glamourMaxWidth  = 120
	repoNameMaxWidth = 32
)

type Bubble struct {
	readmeViewport *vp.ViewportBubble
	repo           *git.Repo
	styles         *style.Styles
	height         int
	heightMargin   int
	width          int
	widthMargin    int
}

func NewBubble(repo *git.Repo, styles *style.Styles, width, wm, height, hm int) *Bubble {
	b := &Bubble{
		readmeViewport: &vp.ViewportBubble{
			Viewport: &viewport.Model{},
		},
		repo:         repo,
		styles:       styles,
		widthMargin:  wm,
		heightMargin: hm,
	}
	b.SetSize(width, height)
	return b
}
func (b *Bubble) Init() tea.Cmd {
	return b.setupCmd
}

func (b *Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.SetSize(msg.Width, msg.Height)
		// XXX: if we find that longer readmes take more than a few
		// milliseconds to render we may need to move Glamour rendering into a
		// command.
		md, err := b.glamourize(b.repo.Readme)
		if err != nil {
			return b, nil
		}
		b.readmeViewport.Viewport.SetContent(md)
	}
	rv, cmd := b.readmeViewport.Update(msg)
	b.readmeViewport = rv.(*vp.ViewportBubble)
	cmds = append(cmds, cmd)
	return b, tea.Batch(cmds...)
}

func (b *Bubble) SetSize(w, h int) {
	b.width = w
	b.height = h
	b.readmeViewport.Viewport.Width = w - b.widthMargin
	b.readmeViewport.Viewport.Height = h - b.heightMargin
}

func (b *Bubble) GotoTop() {
	b.readmeViewport.Viewport.GotoTop()
}

func (b *Bubble) View() string {
	return b.readmeViewport.View()
}

func (b *Bubble) Help() []types.HelpEntry {
	return []types.HelpEntry{
		{"f/b", "pgup/pgdown"},
	}
}

func (b *Bubble) setupCmd() tea.Msg {
	md, err := b.glamourize(b.repo.Readme)
	if err != nil {
		return types.ErrMsg{err}
	}
	b.readmeViewport.Viewport.SetContent(md)
	b.GotoTop()
	return nil
}

func (b *Bubble) glamourize(md string) (string, error) {
	w := b.width - b.widthMargin // - repoBody.GetHorizontalFrameSize()
	if w > glamourMaxWidth {
		w = glamourMaxWidth
	}
	tr, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(w),
	)

	if err != nil {
		return "", err
	}
	mdt, err := tr.Render(md)
	if err != nil {
		return "", err
	}
	// For now, truncate long lines in Glamour that would otherwise break the
	// layout when wrapping. This is very likely due to #43 in Reflow, which
	// has to do with a bug in the way lines longer than the given width are
	// wrapped.
	//
	//     https://github.com/muesli/reflow/issues/43
	//
	// TODO: solve this upstream in Glamour/Reflow.
	mdt = lipgloss.NewStyle().MaxWidth(w).Render(mdt)
	return mdt, nil
}

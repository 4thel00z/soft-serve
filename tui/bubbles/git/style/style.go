package style

import (
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	ActiveBorderColor   lipgloss.Color
	InactiveBorderColor lipgloss.Color

	RepoTitle       lipgloss.Style
	RepoTitleBox    lipgloss.Style
	RepoNote        lipgloss.Style
	RepoNoteBox     lipgloss.Style
	RepoBody        lipgloss.Style
	RepoBodyBorder  lipgloss.Border
	RepoTitleBorder lipgloss.Border
	RepoNoteBorder  lipgloss.Border

	LogItemSelector   lipgloss.Style
	LogItemActive     lipgloss.Style
	LogItemInactive   lipgloss.Style
	LogItemHash       lipgloss.Style
	LogCommit         lipgloss.Style
	LogCommitHash     lipgloss.Style
	LogCommitAuthor   lipgloss.Style
	LogCommitDate     lipgloss.Style
	LogCommitBody     lipgloss.Style
	LogCommitStatsAdd lipgloss.Style
	LogCommitStatsDel lipgloss.Style

	TreeFileDir  lipgloss.Style
	TreeFileMode lipgloss.Style
	TreeFileSize lipgloss.Style
}

func DefaultStyles() *Styles {
	s := &Styles{}

	s.ActiveBorderColor = lipgloss.Color("62")
	s.InactiveBorderColor = lipgloss.Color("236")

	s.RepoTitle = lipgloss.NewStyle().
		Padding(0, 2).
		Foreground(lipgloss.Color("252"))

	s.RepoTitleBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "┬",
		BottomLeft:  "├",
		BottomRight: "┴",
	}

	s.RepoNoteBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "┬",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┤",
	}

	s.RepoBodyBorder = lipgloss.Border{
		Top:         "",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "",
		TopRight:    "",
		BottomLeft:  "╰",
		BottomRight: "╯",
	}

	s.RepoTitle = lipgloss.NewStyle().
		Padding(0, 2).
		Foreground(lipgloss.Color("252"))

	s.RepoTitleBox = lipgloss.NewStyle().
		BorderStyle(s.RepoTitleBorder).
		BorderForeground(s.InactiveBorderColor)

	s.RepoNote = lipgloss.NewStyle().
		Padding(0, 2).
		Foreground(lipgloss.Color("168"))

	s.RepoNoteBox = lipgloss.NewStyle().
		BorderStyle(s.RepoNoteBorder).
		BorderForeground(s.InactiveBorderColor).
		BorderTop(true).
		BorderRight(true).
		BorderBottom(true).
		BorderLeft(false)

	s.RepoBody = lipgloss.NewStyle().
		BorderStyle(s.RepoBodyBorder).
		BorderForeground(s.InactiveBorderColor).
		PaddingRight(1)

	s.LogItemInactive = lipgloss.NewStyle()

	s.LogItemSelector = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#B083EA"))

	s.LogItemActive = s.LogItemInactive.Copy().
		Bold(true)

	s.LogItemHash = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E8E8A8"))

	s.LogCommit = lipgloss.NewStyle().
		Margin(0, 2)

	s.LogCommitHash = s.LogItemHash.Copy().
		Bold(true)

	s.LogCommitBody = lipgloss.NewStyle().
		MarginTop(1).
		MarginLeft(2).
		Width(80)

	s.LogCommitStatsAdd = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00D787")).
		Bold(true)

	s.LogCommitStatsDel = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FD5B5B")).
		Bold(true)

	s.TreeFileDir = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00AAFF"))

	s.TreeFileMode = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#777777"))

	s.TreeFileSize = lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	return s
}

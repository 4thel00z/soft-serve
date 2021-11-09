package types

import (
	"github.com/charmbracelet/glamour"
	gansi "github.com/charmbracelet/glamour/ansi"
	"github.com/muesli/termenv"
)

var (
	RenderCtx = DefaultRenderCtx()
)

func DefaultRenderCtx() gansi.RenderContext {
	var styles = "light"
	if termenv.HasDarkBackground() {
		styles = "dark"
	}
	return gansi.NewRenderContext(gansi.Options{
		ColorProfile: termenv.TrueColor,
		Styles:       *glamour.DefaultStyles[styles],
	})
}

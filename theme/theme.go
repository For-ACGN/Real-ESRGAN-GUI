package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type Chinese struct{}

var _ fyne.Theme = (*Chinese)(nil)

func (Chinese) Font(fyne.TextStyle) fyne.Resource {
	return resourceMsyhTtc
}

func (Chinese) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
}

func (Chinese) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (Chinese) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}

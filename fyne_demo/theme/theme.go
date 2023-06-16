package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	"fyne_demo/resource/font"
	"fyne_demo/resource/image"
)

var _ fyne.Theme = (*ChineseTheme)(nil)

type ChineseTheme struct {
}

func (c ChineseTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (c ChineseTheme) Font(style fyne.TextStyle) fyne.Resource {
	return &fyne.StaticResource{
		StaticName:    "default.ttf",
		StaticContent: font.Font,
	}
}

func (c ChineseTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(fyne.ThemeIconName(image.Icon))
}

func (c ChineseTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

package widgets

import (
	"retroart-sdl2/internal/theme"

	"github.com/TotallyGamerJet/clay"
)

// Cached design system instance to avoid repeated calls
var cachedDesignSystem *theme.DesignSystem

// getDesignSystem returns a cached instance of the design system
func getDesignSystem() *theme.DesignSystem {
	if cachedDesignSystem == nil {
		ds := theme.DefaultDesignSystem()
		cachedDesignSystem = &ds
	}
	return cachedDesignSystem
}

// Text creates a text element using the theme's typography system
func Text(content string, fontSize uint16, color clay.Color) {
	fontId := theme.GetFontIdForSize(fontSize)

	textConfig := &clay.TextElementConfig{
		FontId:    fontId,
		FontSize:  fontSize,
		TextColor: color,
	}

	clay.Text(content, textConfig)
}

// Convenience functions using the design system typography sizes

// TextXSmall creates text with XSmall typography size
func TextXSmall(content string, color clay.Color) {
	ds := getDesignSystem()
	Text(content, ds.Typography.XSmall, color)
}

// TextSmall creates text with Small typography size
func TextSmall(content string, color clay.Color) {
	ds := getDesignSystem()
	Text(content, ds.Typography.Small, color)
}

// TextBase creates text with Base typography size
func TextBase(content string, color clay.Color) {
	ds := getDesignSystem()
	Text(content, ds.Typography.Base, color)
}

// TextLarge creates text with Large typography size
func TextLarge(content string, color clay.Color) {
	ds := getDesignSystem()
	Text(content, ds.Typography.Large, color)
}

// TextXLarge creates text with XLarge typography size
func TextXLarge(content string, color clay.Color) {
	ds := getDesignSystem()
	Text(content, ds.Typography.XLarge, color)
}

// Convenience functions with common colors

// TextPrimary creates text with primary color
func TextPrimary(content string, fontSize uint16) {
	ds := getDesignSystem()
	Text(content, fontSize, ds.Colors.TextPrimary)
}

// TextSecondary creates text with secondary color
func TextSecondary(content string, fontSize uint16) {
	ds := getDesignSystem()
	Text(content, fontSize, ds.Colors.TextSecondary)
}

// TextMuted creates text with muted color
func TextMuted(content string, fontSize uint16) {
	ds := getDesignSystem()
	Text(content, fontSize, ds.Colors.TextMuted)
}

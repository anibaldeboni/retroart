package theme

import "github.com/TotallyGamerJet/clay"

// VirtualKeyboardStyle contém configurações para o teclado virtual
type VirtualKeyboardStyle struct {
	BackgroundColor clay.Color
	Padding         clay.Padding
	CornerRadius    float32
	KeySpacing      uint16
	MaxWidth        float32
	MaxHeight       float32
	KeyButtonStyle  KeyButtonStyle
	FontSize        uint16
}

// KeyButtonStyle contém configurações para botões individuais do teclado
type KeyButtonStyle struct {
	Width        float32
	Height       float32
	Padding      clay.Padding
	FontSize     uint16
	CornerRadius float32
	// Estados normais
	BackgroundColor clay.Color
	BorderColor     clay.Color
	TextColor       clay.Color
	// Estados focados
	FocusedBackgroundColor clay.Color
	FocusedBorderColor     clay.Color
	FocusedTextColor       clay.Color
	// Estados pressionados
	PressedBackgroundColor clay.Color
	PressedBorderColor     clay.Color
	PressedTextColor       clay.Color
}

// GetVirtualKeyboardStyle retorna a configuração de estilo para o teclado virtual
func (ds DesignSystem) GetVirtualKeyboardStyle() VirtualKeyboardStyle {
	return VirtualKeyboardStyle{
		BackgroundColor: ds.Colors.SurfaceTertiary,
		Padding:         clay.Padding{Left: ds.Spacing.SM, Right: ds.Spacing.SM, Top: ds.Spacing.SM, Bottom: ds.Spacing.SM},
		CornerRadius:    ds.Border.Radius.Large,
		KeySpacing:      ds.Spacing.XS,
		MaxWidth:        480,
		MaxHeight:       200,
		KeyButtonStyle: KeyButtonStyle{
			Width:        33,
			Height:       36,
			Padding:      clay.Padding{Left: 4, Right: 4, Top: 4, Bottom: 4},
			FontSize:     ds.Typography.Base,
			CornerRadius: ds.Border.Radius.Large,

			// Estado normal
			BackgroundColor: ds.Colors.SurfaceSecondary,
			BorderColor:     ds.Colors.Border,
			TextColor:       ds.Colors.TextPrimary,

			// Estado focado
			FocusedBackgroundColor: ds.Colors.Primary,
			FocusedBorderColor:     ds.Colors.PrimaryHover,
			FocusedTextColor:       ds.Colors.TextOnPrimary,

			// Estado pressionado
			PressedBackgroundColor: ds.Colors.PrimaryActive,
			PressedBorderColor:     ds.Colors.Primary,
			PressedTextColor:       ds.Colors.TextOnPrimary,
		},
	}
}

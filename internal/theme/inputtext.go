package theme

import "github.com/TotallyGamerJet/clay"

// ComponentStyleType define os tipos de estilo disponíveis
type ComponentStyleType string

const (
	StylePrimary   ComponentStyleType = "primary"
	StyleSecondary ComponentStyleType = "secondary"
	StyleDanger    ComponentStyleType = "danger"
	StyleSuccess   ComponentStyleType = "success"
	StyleWarning   ComponentStyleType = "warning"
	StyleInfo      ComponentStyleType = "info"
)

// InputTextStyle contém configurações para campos de texto de entrada
type InputTextStyle struct {
	Sizing           clay.Sizing
	Padding          clay.Padding
	FontSize         uint16
	CornerRadius     float32
	BorderWidth      uint16
	BackgroundColor  clay.Color
	BorderColor      clay.Color
	TextColor        clay.Color
	CursorColor      clay.Color
	PlaceholderColor clay.Color
	// Estados
	FocusedBackgroundColor clay.Color
	FocusedBorderColor     clay.Color
	FocusedTextColor       clay.Color
}

// GetInputTextStyle retorna a configuração de estilo para campos de texto
func (ds DesignSystem) GetInputTextStyle() InputTextStyle {
	return InputTextStyle{
		Padding:          clay.Padding{Left: ds.Spacing.MD, Right: ds.Spacing.MD, Top: ds.Spacing.SM, Bottom: ds.Spacing.SM},
		FontSize:         ds.Typography.Base,
		CornerRadius:     ds.Border.Radius.Large,
		BorderWidth:      ds.Border.Width.XSmall,
		BackgroundColor:  ds.Colors.InputBackground,
		BorderColor:      ds.Colors.InputBorder,
		TextColor:        ds.Colors.TextPrimary,
		CursorColor:      ds.Colors.Primary,
		PlaceholderColor: ds.Colors.TextPlaceholder,

		// Estados focados
		FocusedBackgroundColor: ds.Colors.InputBackgroundFocused,
		FocusedBorderColor:     ds.Colors.InputBorderFocused,
		FocusedTextColor:       ds.Colors.TextPrimary,
	}
}

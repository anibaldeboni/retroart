package theme

import "github.com/TotallyGamerJet/clay"

// ButtonColor define a aparência de um estado específico do botão
type ButtonColor struct {
	BackgroundColor clay.Color
	TextColor       clay.Color
}

// ButtonStyle contém configurações completas para botões
type ButtonStyle struct {
	Padding      clay.Padding
	FontSize     uint16
	CornerRadius float32
	Normal       ButtonColor
	Focused      ButtonColor
}

// GetButtonStyle retorna a configuração de estilo para um botão baseado no tipo
func (ds DesignSystem) GetButtonStyle(styleType ComponentStyleType) ButtonStyle {
	baseStyle := ButtonStyle{
		Padding:      clay.Padding{Left: ds.Spacing.LG, Right: ds.Spacing.LG, Top: ds.Spacing.MD, Bottom: ds.Spacing.MD},
		FontSize:     ds.Typography.Large,
		CornerRadius: ds.Border.Radius.Large,
	}

	switch styleType {
	case StylePrimary:
		baseStyle.Normal = ButtonColor{
			BackgroundColor: ds.Colors.Primary,
			TextColor:       ds.Colors.TextOnPrimary,
		}
		baseStyle.Focused = ButtonColor{
			BackgroundColor: ds.Colors.PrimaryHover,
			TextColor:       ds.Colors.TextPrimary,
		}

	case StyleSecondary:
		baseStyle.Normal = ButtonColor{
			BackgroundColor: ds.Colors.Secondary,
			TextColor:       ds.Colors.TextOnSecondary,
		}
		baseStyle.Focused = ButtonColor{
			BackgroundColor: ds.Colors.SecondaryHover,
			TextColor:       ds.Colors.TextPrimary,
		}

	case StyleDanger:
		baseStyle.Normal = ButtonColor{
			BackgroundColor: ds.Colors.Danger,
			TextColor:       ds.Colors.TextOnDanger,
		}
		baseStyle.Focused = ButtonColor{
			BackgroundColor: ds.Colors.DangerHover,
			TextColor:       ds.Colors.TextPrimary,
		}

	case StyleSuccess:
		baseStyle.Normal = ButtonColor{
			BackgroundColor: ds.Colors.Success,
			TextColor:       ds.Colors.TextPrimary,
		}
		baseStyle.Focused = ButtonColor{
			BackgroundColor: ds.Colors.SuccessHover,
			TextColor:       ds.Colors.TextPrimary,
		}

	default:
		// Fallback para primary
		return ds.GetButtonStyle(StylePrimary)
	}

	return baseStyle
}

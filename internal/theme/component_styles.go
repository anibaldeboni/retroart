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

// ButtonState define a aparência de um estado específico do botão
type ButtonState struct {
	BackgroundColor clay.Color
	TextColor       clay.Color
}

// ButtonStyle contém configurações completas para botões
type ButtonStyle struct {
	Sizing       clay.Sizing
	Padding      clay.Padding
	TextSize     uint16
	CornerRadius float32
	Normal       ButtonState
	Focused      ButtonState
}

// CheckboxListStyle contém configurações para checkbox lists
type CheckboxListStyle struct {
	Sizing          clay.Sizing
	Padding         clay.Padding
	ChildGap        uint16
	BackgroundColor clay.Color
	MaxHeight       float32
	ScrollOffset    int
	CheckboxSize    float32
	ItemHeight      float32

	// Cores para diferentes estados dos itens
	ItemNormalBg     clay.Color
	ItemSelectedBg   clay.Color
	ItemFocusedBg    clay.Color
	ItemNormalText   clay.Color
	ItemSelectedText clay.Color
	ItemFocusedText  clay.Color

	// Cores do checkbox
	CheckboxNormal   clay.Color
	CheckboxSelected clay.Color
	CheckboxMark     clay.Color
}

// ContainerStyle define estilos para containers
type ContainerStyle struct {
	BackgroundColor clay.Color
	Padding         clay.Padding
	CornerRadius    float32
}

// GetButtonStyle retorna a configuração de estilo para um botão baseado no tipo
func (ds DesignSystem) GetButtonStyle(styleType ComponentStyleType) ButtonStyle {
	baseStyle := ButtonStyle{
		Sizing: clay.Sizing{
			Width:  clay.SizingFixed(220),
			Height: clay.SizingFixed(45),
		},
		Padding:      clay.Padding{Left: ds.Spacing.LG, Right: ds.Spacing.LG, Top: ds.Spacing.MD, Bottom: ds.Spacing.MD},
		TextSize:     ds.Typography.Base,
		CornerRadius: ds.BorderRadius.Large,
	}

	switch styleType {
	case StylePrimary:
		baseStyle.Normal = ButtonState{
			BackgroundColor: ds.Colors.Primary,
			TextColor:       ds.Colors.TextOnPrimary,
		}
		baseStyle.Focused = ButtonState{
			BackgroundColor: ds.Colors.PrimaryHover,
			TextColor:       ds.Colors.TextPrimary,
		}

	case StyleSecondary:
		baseStyle.Normal = ButtonState{
			BackgroundColor: ds.Colors.Secondary,
			TextColor:       clay.Color{R: 220, G: 200, B: 255, A: 255}, // Texto roxo claro específico
		}
		baseStyle.Focused = ButtonState{
			BackgroundColor: ds.Colors.SecondaryHover,
			TextColor:       ds.Colors.TextPrimary,
		}

	case StyleDanger:
		baseStyle.Normal = ButtonState{
			BackgroundColor: ds.Colors.Danger,
			TextColor:       ds.Colors.TextOnDanger,
		}
		baseStyle.Focused = ButtonState{
			BackgroundColor: ds.Colors.DangerHover,
			TextColor:       ds.Colors.TextPrimary,
		}

	case StyleSuccess:
		baseStyle.Normal = ButtonState{
			BackgroundColor: ds.Colors.Success,
			TextColor:       ds.Colors.TextPrimary,
		}
		baseStyle.Focused = ButtonState{
			BackgroundColor: ds.Colors.SuccessHover,
			TextColor:       ds.Colors.TextPrimary,
		}

	default:
		// Fallback para primary
		return ds.GetButtonStyle(StylePrimary)
	}

	return baseStyle
}

// GetCheckboxListStyle retorna a configuração de estilo para checkbox lists
func (ds DesignSystem) GetCheckboxListStyle() CheckboxListStyle {
	return CheckboxListStyle{
		Sizing: clay.Sizing{
			Width:  clay.SizingGrow(1),
			Height: clay.SizingFixed(300),
		},
		Padding:         clay.Padding{Left: ds.Spacing.MD, Right: ds.Spacing.MD, Top: ds.Spacing.SM, Bottom: ds.Spacing.SM},
		ChildGap:        ds.Spacing.XS,
		BackgroundColor: ds.Colors.SurfaceSecondary,
		MaxHeight:       300,
		ScrollOffset:    0,
		CheckboxSize:    18,
		ItemHeight:      35,

		// Estados dos itens
		ItemNormalBg:     clay.Color{R: 0, G: 0, B: 0, A: 0}, // Transparente
		ItemSelectedBg:   ds.Colors.Success,
		ItemFocusedBg:    ds.Colors.Info,
		ItemNormalText:   clay.Color{R: 220, G: 230, B: 245, A: 255}, // Cinza claro específico
		ItemSelectedText: clay.Color{R: 180, G: 255, B: 200, A: 255}, // Verde claro específico
		ItemFocusedText:  ds.Colors.TextPrimary,

		// Checkbox
		CheckboxNormal:   ds.Colors.Border,
		CheckboxSelected: ds.Colors.SuccessHover,
		CheckboxMark:     ds.Colors.TextPrimary,
	}
}

// GetMainContainerStyle retorna o estilo para o container principal
func (ds DesignSystem) GetMainContainerStyle() ContainerStyle {
	return ContainerStyle{
		BackgroundColor: ds.Colors.Background,
		Padding:         clay.Padding{Left: 0, Right: 0, Top: 0, Bottom: 0},
		CornerRadius:    0,
	}
}

// InputTextStyle contém configurações para campos de texto de entrada
type InputTextStyle struct {
	Sizing           clay.Sizing
	Padding          clay.Padding
	TextSize         uint16
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

// VirtualKeyboardStyle contém configurações para o teclado virtual
type VirtualKeyboardStyle struct {
	BackgroundColor clay.Color
	Padding         clay.Padding
	CornerRadius    float32
	KeySpacing      uint16
	MaxWidth        float32
	MaxHeight       float32
	KeyButtonStyle  KeyButtonStyle
}

// KeyButtonStyle contém configurações para botões individuais do teclado
type KeyButtonStyle struct {
	Width        float32
	Height       float32
	Padding      clay.Padding
	TextSize     uint16
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

// GetInputTextStyle retorna a configuração de estilo para campos de texto
func (ds DesignSystem) GetInputTextStyle() InputTextStyle {
	return InputTextStyle{
		Sizing: clay.Sizing{
			Width:  clay.SizingFixed(300),
			Height: clay.SizingFixed(40),
		},
		Padding:          clay.Padding{Left: ds.Spacing.MD, Right: ds.Spacing.MD, Top: ds.Spacing.SM, Bottom: ds.Spacing.SM},
		TextSize:         ds.Typography.Base,
		CornerRadius:     ds.BorderRadius.Medium,
		BorderWidth:      1,
		BackgroundColor:  clay.Color{R: 80, G: 85, B: 95, A: 120}, // Lighter background
		BorderColor:      ds.Colors.Border,
		TextColor:        ds.Colors.TextPrimary,
		CursorColor:      ds.Colors.Primary,
		PlaceholderColor: ds.Colors.TextMuted,

		// Estados focados
		FocusedBackgroundColor: clay.Color{R: 90, G: 95, B: 105, A: 140}, // Even lighter when focused
		FocusedBorderColor:     ds.Colors.Primary,
		FocusedTextColor:       ds.Colors.TextPrimary,
	}
}

// GetVirtualKeyboardStyle retorna a configuração de estilo para o teclado virtual
func (ds DesignSystem) GetVirtualKeyboardStyle() VirtualKeyboardStyle {
	return VirtualKeyboardStyle{
		BackgroundColor: clay.Color{R: 45, G: 47, B: 59, A: 180}, // Semi-transparent dark
		Padding:         clay.Padding{Left: ds.Spacing.LG, Right: ds.Spacing.LG, Top: ds.Spacing.LG, Bottom: ds.Spacing.LG},
		CornerRadius:    ds.BorderRadius.Large,
		KeySpacing:      ds.Spacing.XS,
		MaxWidth:        480,
		MaxHeight:       200,
		KeyButtonStyle: KeyButtonStyle{
			Width:        33,
			Height:       36,
			Padding:      clay.Padding{Left: 4, Right: 4, Top: 4, Bottom: 4},
			TextSize:     ds.Typography.Small,
			CornerRadius: ds.BorderRadius.Small,

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

// GetContentContainerStyle retorna o estilo para containers de conteúdo
func (ds DesignSystem) GetContentContainerStyle() ContainerStyle {
	return ContainerStyle{
		BackgroundColor: ds.Colors.Surface,
		Padding:         clay.Padding{Left: ds.Spacing.LG, Right: ds.Spacing.LG, Top: ds.Spacing.LG, Bottom: ds.Spacing.LG},
		CornerRadius:    ds.BorderRadius.Medium,
	}
}

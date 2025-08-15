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

// ButtonColor define a aparência de um estado específico do botão
type ButtonColor struct {
	BackgroundColor clay.Color
	TextColor       clay.Color
}

// ButtonStyle contém configurações completas para botões
type ButtonStyle struct {
	Padding      clay.Padding
	TextSize     uint16
	CornerRadius float32
	Normal       ButtonColor
	Focused      ButtonColor
}

type Checkbox struct {
	Size            float32
	Border          clay.BorderElementConfig
	CornerRadius    float32
	Background      clay.Color
	Mark            CheckboxMark
	ScrollIndicator CheckboxScrollIndicator
	Color           CheckboxColor
}

type CheckboxMark struct {
	Symbol string
	Size   uint16
}

type CheckboxScrollIndicator struct {
	Size       uint16
	UpSymbol   string
	DownSymbol string
}

type CheckboxColor struct {
	Normal   clay.Color
	Selected clay.Color
	Mark     clay.Color
}

// CheckboxListStyle contém configurações para checkbox lists
type CheckboxListStyle struct {
	FontSize        uint16
	Padding         clay.Padding
	ChildGap        uint16
	BackgroundColor clay.Color
	ScrollOffset    int
	Checkbox        Checkbox
	CornerRadius    float32

	// Cores para diferentes estados dos itens
	ItemNormalBg     clay.Color
	ItemSelectedBg   clay.Color
	ItemFocusedBg    clay.Color
	ItemNormalText   clay.Color
	ItemSelectedText clay.Color
	ItemFocusedText  clay.Color
}

// ContainerStyle define estilos para containers
type ContainerStyle struct {
	Image           clay.ImageElementConfig
	BackgroundColor clay.Color
	Padding         clay.Padding
	CornerRadius    float32
	Border          clay.BorderElementConfig
}

// GetButtonStyle retorna a configuração de estilo para um botão baseado no tipo
func (ds DesignSystem) GetButtonStyle(styleType ComponentStyleType) ButtonStyle {
	baseStyle := ButtonStyle{
		Padding:      clay.Padding{Left: ds.Spacing.LG, Right: ds.Spacing.LG, Top: ds.Spacing.MD, Bottom: ds.Spacing.MD},
		TextSize:     ds.Typography.Base,
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
			TextColor:       clay.Color{R: 220, G: 200, B: 255, A: 255}, // Texto roxo claro específico
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

// GetCheckboxListStyle retorna a configuração de estilo para checkbox lists
func (ds DesignSystem) GetCheckboxListStyle() CheckboxListStyle {
	return CheckboxListStyle{
		Padding:         clay.Padding{Left: ds.Spacing.SM, Right: ds.Spacing.SM, Top: ds.Spacing.MD, Bottom: ds.Spacing.MD},
		ChildGap:        ds.Spacing.SM,
		BackgroundColor: ds.Colors.SurfaceSecondary,
		ScrollOffset:    0,
		FontSize:        ds.Typography.Small,
		Checkbox: Checkbox{
			Size:         22,
			CornerRadius: ds.Border.Radius.Small,
			Color: CheckboxColor{
				Normal:   ds.Colors.Border,
				Selected: ds.Colors.SuccessHover,
				Mark:     ds.Colors.TextPrimary,
			},
			Mark: CheckboxMark{
				Symbol: "◣",
				Size:   ds.Typography.Small,
			},
			ScrollIndicator: CheckboxScrollIndicator{
				Size:       ds.Typography.Small,
				UpSymbol:   "▲",
				DownSymbol: "▼",
			},
		},
		CornerRadius: ds.Border.Radius.Large,

		// Estados dos itens
		ItemNormalBg:     clay.Color{R: 0, G: 0, B: 0, A: 0}, // Transparente
		ItemSelectedBg:   ds.Colors.Success,
		ItemFocusedBg:    ds.Colors.Info,
		ItemNormalText:   clay.Color{R: 220, G: 230, B: 245, A: 255}, // Cinza claro específico
		ItemSelectedText: clay.Color{R: 180, G: 255, B: 200, A: 255}, // Verde claro específico
		ItemFocusedText:  ds.Colors.TextPrimary,
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
		Padding:          clay.Padding{Left: ds.Spacing.MD, Right: ds.Spacing.MD, Top: ds.Spacing.SM, Bottom: ds.Spacing.SM},
		TextSize:         ds.Typography.Base,
		CornerRadius:     ds.Border.Radius.Large,
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
		CornerRadius:    ds.Border.Radius.Large,
		KeySpacing:      ds.Spacing.XS,
		MaxWidth:        480,
		MaxHeight:       200,
		KeyButtonStyle: KeyButtonStyle{
			Width:        33,
			Height:       36,
			Padding:      clay.Padding{Left: 4, Right: 4, Top: 4, Bottom: 4},
			TextSize:     ds.Typography.Small,
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

// GetContentContainerStyle retorna o estilo para containers de conteúdo
func (ds DesignSystem) GetContentContainerStyle() ContainerStyle {
	return ContainerStyle{
		BackgroundColor: ds.Colors.Surface,
		Padding:         clay.Padding{Left: ds.Spacing.LG, Right: ds.Spacing.LG, Top: ds.Spacing.LG, Bottom: ds.Spacing.LG},
		CornerRadius:    ds.Border.Radius.Large,
		Border: clay.BorderElementConfig{
			Width: clay.BorderWidth{Left: ds.Border.Width.Small, Right: ds.Border.Width.Small, Top: ds.Border.Width.Small, Bottom: ds.Border.Width.Small},
			Color: ds.Colors.Border,
		},
	}
}

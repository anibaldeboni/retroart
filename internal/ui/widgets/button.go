package widgets

import (
	"log"
	"retroart-sdl2/internal/input"

	"github.com/TotallyGamerJet/clay"
)

type Button struct {
	ID      string
	Label   string
	Config  ButtonConfig
	OnClick func()
	focused bool
	enabled bool
}

// ButtonConfig contém configurações para diferentes estados do botão
type ButtonConfig struct {
	Sizing       clay.Sizing
	Padding      clay.Padding
	TextSize     uint16
	CornerRadius float32

	Normal  ButtonState
	Focused ButtonState
}

// ButtonState define a aparência de um estado específico do botão
type ButtonState struct {
	BackgroundColor clay.Color
	TextColor       clay.Color
}

// DefaultButtonConfig retorna uma configuração padrão para botões com estado
func DefaultButtonConfig() ButtonConfig {
	return ButtonConfig{
		Sizing: clay.Sizing{
			Width:  clay.SizingFixed(220),
			Height: clay.SizingFixed(45),
		},
		Padding:      clay.Padding{Left: 20, Right: 20, Top: 12, Bottom: 12},
		TextSize:     16,
		CornerRadius: 12,
		Normal: ButtonState{
			BackgroundColor: clay.Color{R: 70, G: 130, B: 220, A: 200},
			TextColor:       clay.Color{R: 220, G: 230, B: 255, A: 255},
		},
		Focused: ButtonState{
			BackgroundColor: clay.Color{R: 30, G: 150, B: 255, A: 255},
			TextColor:       clay.Color{R: 255, G: 255, B: 255, A: 255},
		},
	}
}

// CreateButtonConfig cria uma configuração personalizada para botão com estado
func CreateButtonConfig(normalBg, focusedBg, normalText, focusedText clay.Color) ButtonConfig {
	config := DefaultButtonConfig()
	config.Normal.BackgroundColor = normalBg
	config.Normal.TextColor = normalText
	config.Focused.BackgroundColor = focusedBg
	config.Focused.TextColor = focusedText
	return config
}

// Configurações predefinidas para diferentes tipos de botões
func PrimaryButtonConfig() ButtonConfig {
	return CreateButtonConfig(
		clay.Color{R: 50, G: 120, B: 200, A: 200},  // Azul normal
		clay.Color{R: 30, G: 150, B: 255, A: 255},  // Azul vibrante focado
		clay.Color{R: 220, G: 230, B: 255, A: 255}, // Texto azul claro
		clay.Color{R: 255, G: 255, B: 255, A: 255}, // Texto branco focado
	)
}

func DangerButtonConfig() ButtonConfig {
	return CreateButtonConfig(
		clay.Color{R: 200, G: 60, B: 60, A: 180},   // Vermelho normal
		clay.Color{R: 255, G: 80, B: 80, A: 255},   // Vermelho vibrante focado
		clay.Color{R: 255, G: 200, B: 200, A: 255}, // Texto vermelho claro
		clay.Color{R: 255, G: 255, B: 255, A: 255}, // Texto branco focado
	)
}

func SecondaryButtonConfig() ButtonConfig {
	return CreateButtonConfig(
		clay.Color{R: 120, G: 60, B: 200, A: 180},  // Roxo normal
		clay.Color{R: 150, G: 80, B: 255, A: 255},  // Roxo vibrante focado
		clay.Color{R: 220, G: 200, B: 255, A: 255}, // Texto roxo claro
		clay.Color{R: 255, G: 255, B: 255, A: 255}, // Texto branco focado
	)
}

// NewButton cria um novo botão focável
func NewButton(id, label string, config ButtonConfig, onClick func()) *Button {
	return &Button{
		ID:      id,
		Label:   label,
		Config:  config,
		OnClick: onClick,
		enabled: true,
	}
}

// Interface Focusable implementation
func (b *Button) GetID() string {
	return b.ID
}

func (b *Button) IsFocused() bool {
	return b.focused
}

func (b *Button) OnFocusChanged(focused bool) {
	b.focused = focused
}

func (b *Button) CanFocus() bool {
	return b.enabled
}

func (b *Button) HandleInput(inputType input.InputType) bool {
	if inputType == input.InputConfirm && b.OnClick != nil {
		b.OnClick()
		return true
	}
	return false
}

// Render renderiza o botão usando o sistema stateful
func (b *Button) Render() {
	// Determinar estado atual baseado no foco
	var currentState ButtonState
	if b.focused {
		currentState = b.Config.Focused
	} else {
		currentState = b.Config.Normal
	}

	log.Printf("Creating new button: %s (focused: %t)", b.ID, b.focused)
	clay.UI()(clay.ElementDeclaration{
		Id: clay.ID(b.ID),
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  b.Config.Sizing.Width,
				Height: b.Config.Sizing.Height,
			},
			Padding: b.Config.Padding,
			ChildAlignment: clay.ChildAlignment{
				X: clay.ALIGN_X_CENTER,
				Y: clay.ALIGN_Y_CENTER,
			},
		},
		CornerRadius:    clay.CornerRadiusAll(b.Config.CornerRadius),
		BackgroundColor: currentState.BackgroundColor,
	}, func() {
		// Texto do botão centralizado com cor do estado atual
		clay.Text(b.Label, &clay.TextElementConfig{
			FontSize:  b.Config.TextSize,
			TextColor: currentState.TextColor,
		})
	})
	log.Printf("New button created: %s", b.ID)
}

// SetEnabled habilita/desabilita o botão
func (b *Button) SetEnabled(enabled bool) {
	b.enabled = enabled
}

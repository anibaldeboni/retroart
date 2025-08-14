package widgets

import (
	"log"
	"retroart-sdl2/internal/input"
	"retroart-sdl2/internal/theme"

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
type ButtonConfig = theme.ButtonStyle

// ButtonState define a aparência de um estado específico do botão
type ButtonState = theme.ButtonState

// NewButton cria um novo botão usando o design system
func NewButton(id, label string, styleType theme.ComponentStyleType, onClick func()) *Button {
	return &Button{
		ID:      id,
		Label:   label,
		Config:  theme.GetButtonStyle(styleType),
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
		Text(b.Label, b.Config.TextSize, currentState.TextColor)
	})
	log.Printf("New button created: %s", b.ID)
}

// SetEnabled habilita/desabilita o botão
func (b *Button) SetEnabled(enabled bool) {
	b.enabled = enabled
}

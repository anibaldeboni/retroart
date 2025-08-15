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
	Width   clay.SizingAxis
	Height  clay.SizingAxis
	Config  theme.ButtonStyle
	OnClick func()
	focused bool
	enabled bool
}

// NewButton cria um novo botão usando o design system
func NewButton(id, label string, width, height clay.SizingAxis, styleType theme.ComponentStyleType, onClick func()) *Button {
	return &Button{
		ID:      id,
		Label:   label,
		Width:   width,
		Height:  height,
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

func (n *Button) buttonColor() theme.ButtonColor {
	if n.focused {
		return n.Config.Focused
	}
	return n.Config.Normal
}

// Render draws the button widget on the UI using the current configuration and state.
// It sets up the layout, background color, corner radius, and centers the label text
// with the appropriate color based on the button's state. Debug logs are emitted
// before and after rendering to track button creation and focus state.
func (b *Button) Render() {
	color := b.buttonColor()

	clay.UI()(clay.ElementDeclaration{
		Id: clay.ID(b.ID),
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  b.Width,
				Height: b.Height,
			},
			Padding: b.Config.Padding,
			ChildAlignment: clay.ChildAlignment{
				X: clay.ALIGN_X_CENTER,
				Y: clay.ALIGN_Y_CENTER,
			},
		},
		CornerRadius:    clay.CornerRadiusAll(b.Config.CornerRadius),
		BackgroundColor: color.BackgroundColor,
	}, func() {
		// Texto do botão centralizado com cor do estado atual
		Text(b.Label, b.Config.TextSize, color.TextColor)
	})

	log.Printf("Button: render id=%s focused=%t", b.ID, b.focused)
}

// SetEnabled habilita/desabilita o botão
func (b *Button) SetEnabled(enabled bool) {
	b.enabled = enabled
}

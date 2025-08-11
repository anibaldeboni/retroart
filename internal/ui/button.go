package ui

import (
	"log"

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
func (fb *Button) GetID() string {
	return fb.ID
}

func (fb *Button) IsFocused() bool {
	return fb.focused
}

func (fb *Button) OnFocusChanged(focused bool) {
	fb.focused = focused
}

func (fb *Button) CanFocus() bool {
	return fb.enabled
}

func (fb *Button) HandleInput(direction InputDirection) bool {
	if direction == DirectionConfirm && fb.OnClick != nil {
		fb.OnClick()
		return true
	}
	return false
}

// Render renderiza o botão usando o sistema stateful
func (fb *Button) Render() {
	// Determinar estado atual baseado no foco
	var currentState ButtonState
	if fb.focused {
		currentState = fb.Config.Focused
	} else {
		currentState = fb.Config.Normal
	}

	log.Printf("Creating new button: %s (focused: %t)", fb.ID, fb.focused)
	clay.UI()(clay.ElementDeclaration{
		Id: clay.ID(fb.ID),
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  fb.Config.Sizing.Width,
				Height: fb.Config.Sizing.Height,
			},
			Padding: fb.Config.Padding,
			ChildAlignment: clay.ChildAlignment{
				X: clay.ALIGN_X_CENTER,
				Y: clay.ALIGN_Y_CENTER,
			},
		},
		CornerRadius:    clay.CornerRadiusAll(fb.Config.CornerRadius),
		BackgroundColor: currentState.BackgroundColor,
	}, func() {
		// Texto do botão centralizado com cor do estado atual
		clay.Text(fb.Label, &clay.TextElementConfig{
			FontSize:  fb.Config.TextSize,
			TextColor: currentState.TextColor,
		})
	})
	log.Printf("New button created: %s", fb.ID)
}

// SetEnabled habilita/desabilita o botão
func (fb *Button) SetEnabled(enabled bool) {
	fb.enabled = enabled
}

// FocusableCheckboxList é um wrapper que torna CheckboxList compatível com o sistema de foco
type FocusableCheckboxList[T any] struct {
	*CheckboxList[T]
	focused bool
}

// NewFocusableCheckboxList cria uma nova checkbox list focável
func NewFocusableCheckboxList[T any](id string, items []CheckboxListItem[T], config CheckboxListConfig) *FocusableCheckboxList[T] {
	return &FocusableCheckboxList[T]{
		CheckboxList: NewCheckboxList(id, items, config),
		focused:      false,
	}
}

// Interface Focusable implementation
func (fcl *FocusableCheckboxList[T]) GetID() string {
	return fcl.ID
}

func (fcl *FocusableCheckboxList[T]) IsFocused() bool {
	return fcl.focused
}

func (fcl *FocusableCheckboxList[T]) OnFocusChanged(focused bool) {
	fcl.focused = focused
	fcl.SetFocused(focused) // Delegar para implementação existente
}

func (fcl *FocusableCheckboxList[T]) CanFocus() bool {
	return len(fcl.Items) > 0
}

func (fcl *FocusableCheckboxList[T]) HandleInput(direction InputDirection) bool {
	if !fcl.focused {
		return false
	}

	switch direction {
	case DirectionUp:
		return fcl.MoveFocusUp()
	case DirectionDown:
		return fcl.MoveFocusDown()
	case DirectionConfirm:
		fcl.ToggleFocusedItem()
		return true
	default:
		return false
	}
}

// Render renderiza a checkbox list
func (fcl *FocusableCheckboxList[T]) Render(parentHeight float32) {
	fcl.RenderCheckboxList(parentHeight)
}

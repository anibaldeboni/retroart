package ui

// FocusableButton é um wrapper que torna Button compatível com o sistema de foco
type FocusableButton struct {
	ID      string
	Label   string
	Config  ButtonConfig
	OnClick func()
	focused bool
	enabled bool
}

// NewFocusableButton cria um novo botão focável
func NewFocusableButton(id, label string, config ButtonConfig, onClick func()) *FocusableButton {
	return &FocusableButton{
		ID:      id,
		Label:   label,
		Config:  config,
		OnClick: onClick,
		enabled: true,
	}
}

// Interface Focusable implementation
func (fb *FocusableButton) GetID() string {
	return fb.ID
}

func (fb *FocusableButton) IsFocused() bool {
	return fb.focused
}

func (fb *FocusableButton) OnFocusChanged(focused bool) {
	fb.focused = focused
}

func (fb *FocusableButton) CanFocus() bool {
	return fb.enabled
}

func (fb *FocusableButton) HandleInput(direction InputDirection) bool {
	if direction == DirectionConfirm && fb.OnClick != nil {
		fb.OnClick()
		return true
	}
	return false
}

// Render renderiza o botão usando o sistema stateful
func (fb *FocusableButton) Render(claySystem *ClayLayoutSystem) {
	claySystem.CreateButton(fb.ID, fb.Label, fb.Config, fb.focused, fb.OnClick)
}

// SetEnabled habilita/desabilita o botão
func (fb *FocusableButton) SetEnabled(enabled bool) {
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
func (fcl *FocusableCheckboxList[T]) Render(claySystem *ClayLayoutSystem, parentHeight float32) {
	fcl.RenderCheckboxList(claySystem, parentHeight)
}

package ui

import "retroart-sdl2/internal/input"

// Focusable interface que todos os widgets focáveis devem implementar
type Focusable interface {
	// GetID retorna um identificador único para este widget
	GetID() string

	// IsFocused retorna se este widget está atualmente focado
	IsFocused() bool

	// OnFocusChanged é chamado quando o estado de foco muda
	OnFocusChanged(focused bool)

	// CanFocus retorna se este widget pode receber foco
	CanFocus() bool

	// HandleInput processa input quando focado, retorna true se consumiu o input
	HandleInput(inputType input.InputType) bool
}

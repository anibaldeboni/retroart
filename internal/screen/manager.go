package screen

import (
	"retroart-sdl2/internal/ui"

	"github.com/veandco/go-sdl2/sdl"
)

// Tipos de input suportados
type InputType int

const (
	InputUp InputType = iota
	InputDown
	InputLeft
	InputRight
	InputConfirm
	InputBack
)

// Interface para telas
type Screen interface {
	Update()
	Render()
	HandleInput(input InputType)
	OnEnter()
	OnExit()
}

// Gerenciador de telas
type Manager struct {
	screens       map[string]Screen
	currentScreen Screen
	currentName   string
	layout        *ui.Layout
}

func NewManager(layout *ui.Layout) *Manager {
	return &Manager{
		screens: make(map[string]Screen),
		layout:  layout,
	}
}

func (sm *Manager) AddScreen(name string, screen Screen) {
	sm.screens[name] = screen
}

func (sm *Manager) SetCurrentScreen(name string) {
	if screen, exists := sm.screens[name]; exists {
		if sm.currentScreen != nil {
			sm.currentScreen.OnExit()
		}
		sm.currentScreen = screen
		sm.currentName = name
		sm.currentScreen.OnEnter()
	}
}

func (sm *Manager) GetCurrentScreenName() string {
	return sm.currentName
}

func (sm *Manager) Update() {
	if sm.currentScreen != nil {
		sm.currentScreen.Update()
	}
}

func (sm *Manager) Render() {
	if sm.currentScreen != nil {
		sm.layout.BeginLayout()
		sm.currentScreen.Render()
		sm.layout.Render() // Renderizar o sistema ClayLayout após a tela
	}
}

func (sm *Manager) HandleInput(keycode sdl.Keycode) {
	if sm.currentScreen == nil {
		return
	}

	var input InputType
	switch keycode {
	case sdl.K_UP:
		input = InputUp
	case sdl.K_DOWN:
		input = InputDown
	case sdl.K_LEFT:
		input = InputLeft
	case sdl.K_RIGHT:
		input = InputRight
	case sdl.K_RETURN, sdl.K_SPACE:
		input = InputConfirm
	case sdl.K_ESCAPE:
		input = InputBack
	default:
		return // Tecla não mapeada
	}

	sm.currentScreen.HandleInput(input)
}

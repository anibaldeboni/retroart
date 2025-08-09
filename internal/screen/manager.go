package screen

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
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
	Render(renderer *sdl.Renderer)
	HandleInput(input InputType)
	OnEnter()
	OnExit()
}

// Gerenciador de telas
type Manager struct {
	screens       map[string]Screen
	currentScreen Screen
	currentName   string
	renderer      *sdl.Renderer
	font          *ttf.Font
}

func NewManager(renderer *sdl.Renderer, font *ttf.Font) *Manager {
	return &Manager{
		screens:  make(map[string]Screen),
		renderer: renderer,
		font:     font,
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

func (sm *Manager) Render(renderer *sdl.Renderer) {
	if sm.currentScreen != nil {
		sm.currentScreen.Render(renderer)
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

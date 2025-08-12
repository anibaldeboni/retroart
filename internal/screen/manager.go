package screen

import (
	"retroart-sdl2/internal/input"
	"retroart-sdl2/internal/ui"

	"github.com/veandco/go-sdl2/sdl"
)

// Navigator interface para navegação entre telas
// Permite que telas naveguem sem conhecer o ScreenManager diretamente
type Navigator interface {
	NavigateTo(screenName string)
	GoBack()
	GetCurrentScreenName() string
}

// Interface para telas
type Screen interface {
	Update()
	Render()
	HandleInput(inputType input.InputType)
	OnEnter(navigator Navigator) // Recebe navigator para navegação
	OnExit()
}

// Gerenciador de telas
type Manager struct {
	screens       map[string]Screen
	currentScreen Screen
	currentName   string
	layout        *ui.Layout
	history       []string // Histórico de navegação para GoBack()
}

func NewManager(layout *ui.Layout) *Manager {
	return &Manager{
		screens: make(map[string]Screen),
		layout:  layout,
		history: make([]string, 0),
	}
}

// Navigator interface implementation
func (sm *Manager) NavigateTo(screenName string) {
	sm.SetCurrentScreen(screenName)
}

func (sm *Manager) GoBack() {
	if len(sm.history) > 1 {
		// Remove current screen from history
		sm.history = sm.history[:len(sm.history)-1]
		// Go to previous screen
		previousScreen := sm.history[len(sm.history)-1]
		sm.setCurrentScreenInternal(previousScreen, false) // Don't add to history
	}
}

func (sm *Manager) GetCurrentScreenName() string {
	return sm.currentName
}

func (sm *Manager) AddScreen(name string, screen Screen) {
	sm.screens[name] = screen
}

func (sm *Manager) SetCurrentScreen(name string) {
	sm.setCurrentScreenInternal(name, true)
}

// setCurrentScreenInternal handles screen changes with optional history tracking
func (sm *Manager) setCurrentScreenInternal(name string, addToHistory bool) {
	if screen, exists := sm.screens[name]; exists {
		if sm.currentScreen != nil {
			sm.currentScreen.OnExit()
		}

		// Add to history if requested and not already current
		if addToHistory && sm.currentName != name {
			sm.history = append(sm.history, name)
		}

		sm.currentScreen = screen
		sm.currentName = name
		sm.currentScreen.OnEnter(sm) // Pass self as Navigator
	}
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

	var inputType input.InputType
	switch keycode {
	case sdl.K_UP:
		inputType = input.InputUp
	case sdl.K_DOWN:
		inputType = input.InputDown
	case sdl.K_LEFT:
		inputType = input.InputLeft
	case sdl.K_RIGHT:
		inputType = input.InputRight
	case sdl.K_RETURN, sdl.K_SPACE:
		inputType = input.InputConfirm
	case sdl.K_ESCAPE:
		inputType = input.InputBack
	default:
		return // Tecla não mapeada
	}

	sm.currentScreen.HandleInput(inputType)
}

package screen

import (
	"retroart-sdl2/internal/ui"
)

// BaseScreen fornece funcionalidade básica de foco para todas as telas
type BaseScreen struct {
	focusManager *ui.FocusManager
	screenID     string
}

// NewBaseScreen cria uma nova tela base
func NewBaseScreen(screenID string) *BaseScreen {
	return &BaseScreen{
		focusManager: ui.NewFocusManager(),
		screenID:     screenID,
	}
}

// GetFocusManager retorna o gerenciador de foco da tela
func (bs *BaseScreen) GetFocusManager() *ui.FocusManager {
	return bs.focusManager
}

// AddFocusGroup adiciona um grupo de foco à tela
func (bs *BaseScreen) AddFocusGroup(group *ui.FocusGroup) {
	bs.focusManager.AddGroup(group)
}

// HandleInput processa entrada direcionais, delegando para o focus manager
func (bs *BaseScreen) HandleInput(direction ui.InputDirection) bool {
	return bs.focusManager.HandleInput(direction)
}

// GetCurrentGroup retorna o grupo atualmente ativo
func (bs *BaseScreen) GetCurrentGroup() *ui.FocusGroup {
	return bs.focusManager.GetCurrentGroup()
}

// GetCurrentFocusable retorna o widget atualmente focado
func (bs *BaseScreen) GetCurrentFocusable() ui.Focusable {
	return bs.focusManager.GetCurrentFocusable()
}

// FocusableScreen interface que define o contrato para telas que usam sistema de foco
type FocusableScreen interface {
	// InitializeFocus configura os grupos de foco iniciais da tela
	InitializeFocus()

	// HandleInput processa entrada direcional
	HandleInput(direction ui.InputDirection) bool

	// Render renderiza a tela
	Render()

	// GetFocusManager retorna o gerenciador de foco
	GetFocusManager() *ui.FocusManager
}

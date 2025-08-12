package screen

import (
	"retroart-sdl2/internal/input"
	"retroart-sdl2/internal/ui"
)

// BaseScreen fornece funcionalidade básica de navegação espacial para todas as telas
type BaseScreen struct {
	spatialManager *ui.SpatialNavigationManager
	screenID       string
}

// NewBaseScreen cria uma nova tela base
func NewBaseScreen(screenID string) *BaseScreen {
	return &BaseScreen{
		spatialManager: ui.NewSpatialNavigationManager(),
		screenID:       screenID,
	}
}

// GetSpatialManager retorna o gerenciador de navegação espacial da tela
func (bs *BaseScreen) GetSpatialManager() *ui.SpatialNavigationManager {
	return bs.spatialManager
}

// RegisterWidget registra um widget focável na tela
func (bs *BaseScreen) RegisterWidget(widget ui.Focusable) {
	bs.spatialManager.RegisterWidget(widget)
}

// HandleInput processa entrada direcional, delegando para o spatial manager
func (bs *BaseScreen) HandleInput(inputType input.InputType) bool {
	return bs.spatialManager.HandleInput(inputType)
}

// GetCurrentWidget retorna o widget atualmente focado
func (bs *BaseScreen) GetCurrentWidget() ui.Focusable {
	return bs.spatialManager.GetCurrentWidget()
}

// GetCurrentFocus retorna o ID do widget atualmente focado
func (bs *BaseScreen) GetCurrentFocus() string {
	return bs.spatialManager.GetCurrentFocus()
}

// SpatialScreen interface que define o contrato para telas que usam navegação espacial
type SpatialScreen interface {
	// InitializeWidgets configura os widgets focáveis iniciais da tela
	InitializeWidgets()

	// HandleInput processa entrada direcional
	HandleInput(inputType input.InputType) bool

	// Render renderiza a tela
	Render()

	// GetSpatialManager retorna o gerenciador de navegação espacial
	GetSpatialManager() *ui.SpatialNavigationManager
}

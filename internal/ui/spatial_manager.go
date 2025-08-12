package ui

import (
	"log"
)

// SpatialNavigationManager gerencia a navegação espacial integrada com o Layout
type SpatialNavigationManager struct {
	layout  *Layout
	widgets map[string]Focusable
}

// NewSpatialNavigationManager cria um novo gerenciador de navegação espacial
func NewSpatialNavigationManager() *SpatialNavigationManager {
	return &SpatialNavigationManager{
		layout:  GetLayout(),
		widgets: make(map[string]Focusable),
	}
}

// RegisterWidget registra um widget focável
func (snm *SpatialNavigationManager) RegisterWidget(widget Focusable) {
	id := widget.GetID()
	snm.widgets[id] = widget

	// Registrar no sistema de navegação espacial do layout
	if snm.layout != nil {
		snm.layout.RegisterFocusable(widget)
	}

	log.Printf("SpatialNavigationManager: Registered widget '%s'", id)
}

// UnregisterWidget remove um widget focável
func (snm *SpatialNavigationManager) UnregisterWidget(id string) {
	delete(snm.widgets, id)

	// Remover do sistema de navegação espacial do layout
	if snm.layout != nil {
		snm.layout.UnregisterFocusable(id)
	}

	log.Printf("SpatialNavigationManager: Unregistered widget '%s'", id)
}

// HandleInput processa input de navegação
func (snm *SpatialNavigationManager) HandleInput(direction InputDirection) bool {
	if snm.layout == nil {
		return false
	}

	// Delegar para o sistema de navegação espacial do layout
	return snm.layout.HandleSpatialInput(direction)
}

// GetCurrentWidget retorna o widget atualmente focado
func (snm *SpatialNavigationManager) GetCurrentWidget() Focusable {
	if snm.layout == nil || snm.layout.spatialNav == nil {
		return nil
	}

	return snm.layout.spatialNav.GetCurrentWidget()
}

// GetCurrentFocus retorna o ID do widget atualmente focado
func (snm *SpatialNavigationManager) GetCurrentFocus() string {
	if snm.layout == nil || snm.layout.spatialNav == nil {
		return ""
	}

	return snm.layout.spatialNav.GetCurrentFocus()
}

// SetFocus define o foco em um widget específico
func (snm *SpatialNavigationManager) SetFocus(elementID string) bool {
	if snm.layout == nil || snm.layout.spatialNav == nil {
		return false
	}

	// Encontrar o elemento e definir foco
	for i := range snm.layout.spatialNav.elements {
		if snm.layout.spatialNav.elements[i].ID == elementID {
			return snm.layout.spatialNav.setFocus(elementID)
		}
	}

	return false
}

// DebugPrint imprime informações de debug
func (snm *SpatialNavigationManager) DebugPrint() {
	log.Printf("SpatialNavigationManager: %d registered widgets", len(snm.widgets))
	for id := range snm.widgets {
		log.Printf("  - %s", id)
	}

	if snm.layout != nil && snm.layout.spatialNav != nil {
		snm.layout.spatialNav.DebugPrintElements()
	}
}

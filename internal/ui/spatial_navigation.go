package ui

import (
	"log"
	"math"

	"github.com/TotallyGamerJet/clay"
)

// ElementPosition armazena a posição de um elemento focável
type ElementPosition struct {
	ID          string
	BoundingBox clay.BoundingBox
	Widget      Focusable
}

// SpatialNavigation gerencia navegação baseada em posições espaciais
type SpatialNavigation struct {
	elements     []ElementPosition
	currentFocus string
	focusables   map[string]Focusable
	enabled      bool
}

// NewSpatialNavigation cria um novo sistema de navegação espacial
func NewSpatialNavigation() *SpatialNavigation {
	return &SpatialNavigation{
		elements:   make([]ElementPosition, 0),
		focusables: make(map[string]Focusable),
		enabled:    true,
	}
}

// RegisterFocusable registra um widget focável
func (sn *SpatialNavigation) RegisterFocusable(widget Focusable) {
	id := widget.GetID()
	sn.focusables[id] = widget
	log.Printf("SpatialNavigation: Registered focusable '%s'", id)
}

// UnregisterFocusable remove um widget focável
func (sn *SpatialNavigation) UnregisterFocusable(id string) {
	delete(sn.focusables, id)
	log.Printf("SpatialNavigation: Unregistered focusable '%s'", id)
}

// UpdateLayout atualiza as posições dos elementos baseado nos comandos de renderização do Clay
func (sn *SpatialNavigation) UpdateLayout(commands clay.RenderCommandArray) {
	log.Printf("SpatialNavigation: Processing %d commands, have %d registered focusables", int(commands.Length), len(sn.focusables))

	// List registered focusables for debugging
	log.Println("SpatialNavigation: Registered focusables:")
	for id := range sn.focusables {
		log.Printf("  - '%s'", id)
	}

	// Clear current elements
	sn.elements = sn.elements[:0]

	// Process each command to find focusable elements
	for i := int32(0); i < commands.Length; i++ {
		command := clay.RenderCommandArray_Get(&commands, i)

		log.Printf("SpatialNavigation: Command %d - ID=%d, Type=%d", i, command.Id, command.CommandType)

		// Try to extract element ID
		elementID := sn.extractElementID(command)
		if elementID == "" {
			continue
		}

		// Check if this element is registered as focusable
		widget, exists := sn.focusables[elementID]
		if !exists {
			log.Printf("SpatialNavigation: Element '%s' found but not registered as focusable", elementID)
			continue
		}

		// Add to elements list
		sn.elements = append(sn.elements, ElementPosition{
			ID:          elementID,
			BoundingBox: command.BoundingBox,
			Widget:      widget,
		})

		log.Printf("SpatialNavigation: Added focusable element '%s' at (%f,%f %fx%f)",
			elementID, command.BoundingBox.X, command.BoundingBox.Y,
			command.BoundingBox.Width, command.BoundingBox.Height)
	}

	log.Printf("SpatialNavigation: Updated layout with %d focusable elements", len(sn.elements))
} // clayHashString implements the same hashing algorithm as Clay's Clay__HashString function
func clayHashString(key string, offset uint32, seed uint32) uint32 {
	var hash uint32
	base := seed

	for i := 0; i < len(key); i++ {
		base += uint32(key[i])
		base += (base << 10)
		base ^= (base >> 6)
	}
	hash = base
	hash += offset
	hash += (hash << 10)
	hash ^= (hash >> 6)

	hash += (hash << 3)
	base += (base << 3)
	hash ^= (hash >> 11)
	base ^= (base >> 11)
	hash += (hash << 15)
	base += (base << 15)

	return hash + 1 // Reserve the hash result of zero as "null id"
}

// extractElementID tries to extract a widget ID from a Clay command
func (sn *SpatialNavigation) extractElementID(cmd *clay.RenderCommand) string {
	// Convert the numeric ID back to the original string ID
	cmdID := cmd.Id

	// Try to match against all registered focusable widgets
	for widgetID := range sn.focusables {
		expectedID := clayHashString(widgetID, 0, 0)
		if expectedID == cmdID {
			log.Printf("SpatialNavigation: Found match for ID '%s' (hash=%d)", widgetID, cmdID)
			return widgetID
		}
	}

	log.Printf("SpatialNavigation: No focusable widget found for ID '%d'", cmdID)
	return ""
}

// HandleInput processa input de navegação espacial usando Chain of Responsibility pattern
func (sn *SpatialNavigation) HandleInput(direction InputDirection) bool {
	if !sn.enabled || len(sn.elements) == 0 {
		return false
	}

	currentElement := sn.getCurrentElement()
	if currentElement == nil {
		// Se não há foco atual, focar no primeiro elemento
		return sn.focusFirst()
	}

	// Chain of Responsibility: dar ao widget focado a primeira chance de processar o input
	if currentElement.Widget != nil {
		// Se o widget consome o input, processo termina aqui
		if consumed := currentElement.Widget.HandleInput(direction); consumed {
			log.Printf("SpatialNavigation: Input %d consumed by widget '%s'", direction, currentElement.ID)
			return true
		}
		log.Printf("SpatialNavigation: Input %d not consumed by widget '%s', falling back to spatial navigation", direction, currentElement.ID)
	}

	// Fallback: se widget não consumiu input, usar navegação espacial
	switch direction {
	case DirectionUp:
		return sn.navigateInDirection(currentElement, 0, -1)
	case DirectionDown:
		return sn.navigateInDirection(currentElement, 0, 1)
	case DirectionLeft:
		return sn.navigateInDirection(currentElement, -1, 0)
	case DirectionRight:
		return sn.navigateInDirection(currentElement, 1, 0)
	case DirectionConfirm, DirectionBack:
		// Para Confirm e Back, se chegou até aqui é porque widget não processou
		// Não há fallback espacial para estes, então retorna false
		return false
	}

	return false
}

// navigateInDirection encontra o próximo elemento na direção especificada
func (sn *SpatialNavigation) navigateInDirection(current *ElementPosition, dirX, dirY float32) bool {
	bestElement := sn.findBestElementInDirection(current, dirX, dirY)
	if bestElement != nil {
		return sn.setFocus(bestElement.ID)
	}
	return false
}

// findBestElementInDirection encontra o melhor elemento na direção especificada
func (sn *SpatialNavigation) findBestElementInDirection(current *ElementPosition, dirX, dirY float32) *ElementPosition {
	var bestElement *ElementPosition
	bestDistance := float32(math.Inf(1))

	currentCenter := sn.getCenter(current.BoundingBox)

	for i := range sn.elements {
		element := &sn.elements[i]

		// Pular o elemento atual
		if element.ID == current.ID {
			continue
		}

		// Verificar se o elemento está na direção correta
		if !sn.isInDirection(currentCenter, element.BoundingBox, dirX, dirY) {
			continue
		}

		// Calcular distância
		distance := sn.calculateDistance(currentCenter, element.BoundingBox, dirX, dirY)

		if distance < bestDistance {
			bestDistance = distance
			bestElement = element
		}
	}

	return bestElement
}

// isInDirection verifica se um elemento está na direção especificada
func (sn *SpatialNavigation) isInDirection(fromCenter clay.Vector2, toBoundingBox clay.BoundingBox, dirX, dirY float32) bool {
	toCenter := sn.getCenter(toBoundingBox)

	// Para movimentos horizontais
	if dirX != 0 {
		// Verificar se está na direção horizontal correta
		if dirX > 0 && toCenter.X <= fromCenter.X {
			return false // Queremos ir para direita, mas elemento está à esquerda
		}
		if dirX < 0 && toCenter.X >= fromCenter.X {
			return false // Queremos ir para esquerda, mas elemento está à direita
		}
	}

	// Para movimentos verticais
	if dirY != 0 {
		// Verificar se está na direção vertical correta
		if dirY > 0 && toCenter.Y <= fromCenter.Y {
			return false // Queremos ir para baixo, mas elemento está acima
		}
		if dirY < 0 && toCenter.Y >= fromCenter.Y {
			return false // Queremos ir para cima, mas elemento está abaixo
		}
	}

	return true
}

// calculateDistance calcula a distância entre elementos considerando a direção
func (sn *SpatialNavigation) calculateDistance(fromCenter clay.Vector2, toBoundingBox clay.BoundingBox, dirX, dirY float32) float32 {
	toCenter := sn.getCenter(toBoundingBox)

	// Distância Euclidiana básica
	dx := toCenter.X - fromCenter.X
	dy := toCenter.Y - fromCenter.Y

	distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	// Dar preferência para elementos mais alinhados na direção principal
	if dirX != 0 {
		// Para movimento horizontal, penalizar diferenças verticais
		alignmentPenalty := float32(math.Abs(float64(dy))) * 0.5
		distance += alignmentPenalty
	}

	if dirY != 0 {
		// Para movimento vertical, penalizar diferenças horizontais
		alignmentPenalty := float32(math.Abs(float64(dx))) * 0.5
		distance += alignmentPenalty
	}

	return distance
}

// getCenter retorna o centro de um bounding box
func (sn *SpatialNavigation) getCenter(box clay.BoundingBox) clay.Vector2 {
	return clay.Vector2{
		X: box.X + box.Width/2,
		Y: box.Y + box.Height/2,
	}
}

// setFocus define o foco em um elemento específico
func (sn *SpatialNavigation) setFocus(elementID string) bool {
	// Remover foco anterior
	if sn.currentFocus != "" {
		if widget, exists := sn.focusables[sn.currentFocus]; exists {
			widget.OnFocusChanged(false)
		}
	}

	// Definir novo foco
	if widget, exists := sn.focusables[elementID]; exists && widget.CanFocus() {
		sn.currentFocus = elementID
		widget.OnFocusChanged(true)
		log.Printf("SpatialNavigation: Focus changed to '%s'", elementID)
		return true
	}

	return false
}

// focusFirst foca no primeiro elemento disponível
func (sn *SpatialNavigation) focusFirst() bool {
	if len(sn.elements) > 0 {
		return sn.setFocus(sn.elements[0].ID)
	}
	return false
}

// getCurrentElement retorna o elemento atualmente focado
func (sn *SpatialNavigation) getCurrentElement() *ElementPosition {
	if sn.currentFocus == "" {
		return nil
	}

	for i := range sn.elements {
		if sn.elements[i].ID == sn.currentFocus {
			return &sn.elements[i]
		}
	}

	return nil
}

// GetCurrentFocus retorna o ID do elemento atualmente focado
func (sn *SpatialNavigation) GetCurrentFocus() string {
	return sn.currentFocus
}

// GetCurrentWidget retorna o widget atualmente focado
func (sn *SpatialNavigation) GetCurrentWidget() Focusable {
	if sn.currentFocus == "" {
		return nil
	}
	return sn.focusables[sn.currentFocus]
}

// Clear limpa todos os elementos e foco
func (sn *SpatialNavigation) Clear() {
	sn.elements = sn.elements[:0]
	sn.currentFocus = ""
	sn.focusables = make(map[string]Focusable)
}

// SetEnabled habilita/desabilita a navegação espacial
func (sn *SpatialNavigation) SetEnabled(enabled bool) {
	sn.enabled = enabled
	if !enabled {
		sn.currentFocus = ""
	}
}

// IsEnabled retorna se a navegação espacial está habilitada
func (sn *SpatialNavigation) IsEnabled() bool {
	return sn.enabled
}

// GetElementCount retorna o número de elementos focáveis
func (sn *SpatialNavigation) GetElementCount() int {
	return len(sn.elements)
}

// DebugPrintElements imprime informações de debug sobre os elementos
func (sn *SpatialNavigation) DebugPrintElements() {
	log.Printf("SpatialNavigation Debug: %d elements", len(sn.elements))
	for i, element := range sn.elements {
		focused := element.ID == sn.currentFocus
		log.Printf("  [%d] %s: (%.1f,%.1f) %.1fx%.1f %s",
			i, element.ID,
			element.BoundingBox.X, element.BoundingBox.Y,
			element.BoundingBox.Width, element.BoundingBox.Height,
			map[bool]string{true: "[FOCUSED]", false: ""}[focused])
	}
}

package ui

// InputDirection representa direções de navegação
type InputDirection int

const (
	DirectionUp InputDirection = iota
	DirectionDown
	DirectionLeft
	DirectionRight
	DirectionConfirm
	DirectionBack
)

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
	HandleInput(direction InputDirection) bool
}

// FocusGroup representa um grupo de widgets focáveis
type FocusGroup struct {
	ID           string
	Focusables   []Focusable
	currentIndex int
	enabled      bool
}

// NewFocusGroup cria um novo grupo de foco
func NewFocusGroup(id string) *FocusGroup {
	return &FocusGroup{
		ID:           id,
		Focusables:   make([]Focusable, 0),
		currentIndex: -1,
		enabled:      true,
	}
}

// AddFocusable adiciona um widget focável ao grupo
func (fg *FocusGroup) AddFocusable(focusable Focusable) {
	fg.Focusables = append(fg.Focusables, focusable)

	// Se é o primeiro focusable, torna-se o foco inicial
	if len(fg.Focusables) == 1 && focusable.CanFocus() {
		fg.currentIndex = 0
	}
}

// GetCurrentFocusable retorna o widget atualmente focado no grupo
func (fg *FocusGroup) GetCurrentFocusable() Focusable {
	if fg.currentIndex >= 0 && fg.currentIndex < len(fg.Focusables) {
		return fg.Focusables[fg.currentIndex]
	}
	return nil
}

// MoveFocus move o foco dentro do grupo
func (fg *FocusGroup) MoveFocus(direction InputDirection) bool {
	if !fg.enabled || len(fg.Focusables) == 0 {
		return false
	}

	// Se não há foco atual, focar no primeiro item
	if fg.currentIndex == -1 {
		return fg.FocusFirst()
	}

	current := fg.GetCurrentFocusable()
	if current != nil && current.HandleInput(direction) {
		// O widget atual consumiu o input
		return true
	}

	// Navegar entre widgets do grupo
	switch direction {
	case DirectionUp:
		return fg.MovePrevious()
	case DirectionDown:
		return fg.MoveNext()
	default:
		return false
	}
}

// MoveNext move para o próximo widget focável
func (fg *FocusGroup) MoveNext() bool {
	if len(fg.Focusables) == 0 {
		return false
	}

	startIndex := fg.currentIndex
	for {
		fg.currentIndex = (fg.currentIndex + 1) % len(fg.Focusables)

		if fg.Focusables[fg.currentIndex].CanFocus() {
			fg.updateFocus()
			return true
		}

		// Evitar loop infinito
		if fg.currentIndex == startIndex {
			break
		}
	}
	return false
}

// MovePrevious move para o widget focável anterior
func (fg *FocusGroup) MovePrevious() bool {
	if len(fg.Focusables) == 0 {
		return false
	}

	startIndex := fg.currentIndex
	for {
		fg.currentIndex--
		if fg.currentIndex < 0 {
			fg.currentIndex = len(fg.Focusables) - 1
		}

		if fg.Focusables[fg.currentIndex].CanFocus() {
			fg.updateFocus()
			return true
		}

		// Evitar loop infinito
		if fg.currentIndex == startIndex {
			break
		}
	}
	return false
}

// FocusFirst foca no primeiro widget focável
func (fg *FocusGroup) FocusFirst() bool {
	for i, focusable := range fg.Focusables {
		if focusable.CanFocus() {
			fg.currentIndex = i
			fg.updateFocus()
			return true
		}
	}
	return false
}

// SetEnabled habilita/desabilita o grupo
func (fg *FocusGroup) SetEnabled(enabled bool) {
	if fg.enabled == enabled {
		return
	}

	fg.enabled = enabled
	if !enabled {
		// Perder foco quando desabilitado
		fg.ClearFocus()
	}
}

// ClearFocus remove o foco de todos os widgets do grupo
func (fg *FocusGroup) ClearFocus() {
	for _, focusable := range fg.Focusables {
		if focusable.IsFocused() {
			focusable.OnFocusChanged(false)
		}
	}
}

// updateFocus atualiza o estado de foco de todos os widgets
func (fg *FocusGroup) updateFocus() {
	for i, focusable := range fg.Focusables {
		focused := i == fg.currentIndex && fg.enabled
		if focusable.IsFocused() != focused {
			focusable.OnFocusChanged(focused)
		}
	}
}

// FocusManager gerencia múltiplos grupos de foco
type FocusManager struct {
	groups            []*FocusGroup
	currentGroupIndex int
	enabled           bool
}

// NewFocusManager cria um novo gerenciador de foco
func NewFocusManager() *FocusManager {
	return &FocusManager{
		groups:            make([]*FocusGroup, 0),
		currentGroupIndex: -1,
		enabled:           true,
	}
}

// AddGroup adiciona um grupo de foco
func (fm *FocusManager) AddGroup(group *FocusGroup) {
	fm.groups = append(fm.groups, group)

	// Se é o primeiro grupo, torna-se ativo
	if len(fm.groups) == 1 {
		fm.currentGroupIndex = 0
		fm.updateGroupFocus()
	}
}

// GetCurrentGroup retorna o grupo atualmente ativo
func (fm *FocusManager) GetCurrentGroup() *FocusGroup {
	if fm.currentGroupIndex >= 0 && fm.currentGroupIndex < len(fm.groups) {
		return fm.groups[fm.currentGroupIndex]
	}
	return nil
}

// HandleInput processa entrada direcional
func (fm *FocusManager) HandleInput(direction InputDirection) bool {
	if !fm.enabled {
		return false
	}

	currentGroup := fm.GetCurrentGroup()
	if currentGroup == nil {
		return false
	}

	// Tentar processar no grupo atual primeiro
	if currentGroup.MoveFocus(direction) {
		return true
	}

	// Se não foi processado, tentar mudar de grupo
	switch direction {
	case DirectionLeft:
		return fm.MoveToPreviousGroup()
	case DirectionRight:
		return fm.MoveToNextGroup()
	default:
		return false
	}
}

// MoveToNextGroup move para o próximo grupo
func (fm *FocusManager) MoveToNextGroup() bool {
	if len(fm.groups) <= 1 {
		return false
	}

	startIndex := fm.currentGroupIndex
	for {
		fm.currentGroupIndex = (fm.currentGroupIndex + 1) % len(fm.groups)

		if fm.groups[fm.currentGroupIndex].enabled && len(fm.groups[fm.currentGroupIndex].Focusables) > 0 {
			fm.updateGroupFocus()
			return true
		}

		// Evitar loop infinito
		if fm.currentGroupIndex == startIndex {
			break
		}
	}
	return false
}

// MoveToPreviousGroup move para o grupo anterior
func (fm *FocusManager) MoveToPreviousGroup() bool {
	if len(fm.groups) <= 1 {
		return false
	}

	startIndex := fm.currentGroupIndex

	for {
		fm.currentGroupIndex--
		if fm.currentGroupIndex < 0 {
			fm.currentGroupIndex = len(fm.groups) - 1
		}

		group := fm.groups[fm.currentGroupIndex]

		if group.enabled && len(group.Focusables) > 0 {
			fm.updateGroupFocus()
			return true
		}

		// Evitar loop infinito
		if fm.currentGroupIndex == startIndex {
			break
		}
	}
	return false
}

// SetEnabled habilita/desabilita o gerenciador
func (fm *FocusManager) SetEnabled(enabled bool) {
	fm.enabled = enabled
	if !enabled {
		fm.ClearAllFocus()
	} else {
		fm.updateGroupFocus()
	}
}

// ClearAllFocus remove foco de todos os grupos
func (fm *FocusManager) ClearAllFocus() {
	for _, group := range fm.groups {
		group.ClearFocus()
	}
}

// updateGroupFocus atualiza qual grupo está ativo
func (fm *FocusManager) updateGroupFocus() {
	for i, group := range fm.groups {
		if i == fm.currentGroupIndex && fm.enabled {
			// Focar no grupo atual
			if group.currentIndex == -1 {
				group.FocusFirst()
			} else {
				group.updateFocus()
			}
		} else {
			// Remover foco dos outros grupos, mas mantê-los habilitados
			group.ClearFocus()
		}
	}
}

// GetCurrentFocusable retorna o widget atualmente focado
func (fm *FocusManager) GetCurrentFocusable() Focusable {
	currentGroup := fm.GetCurrentGroup()
	if currentGroup != nil {
		return currentGroup.GetCurrentFocusable()
	}
	return nil
}

// FocusableScreen interface para telas que suportam foco
type FocusableScreen interface {
	GetFocusManager() *FocusManager
	InitializeFocus()
}

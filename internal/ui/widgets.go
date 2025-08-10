package ui

import (
	"fmt"
	"log"

	"github.com/TotallyGamerJet/clay"
)

// CheckboxListItem é um item da lista com tipo genérico para Value
type CheckboxListItem[T any] struct {
	Label    string
	Value    T
	Selected bool
}

type CheckboxListConfig struct {
	Sizing          clay.Sizing
	Padding         clay.Padding
	ChildGap        uint16
	BackgroundColor clay.Color
	MaxHeight       float32
	ScrollOffset    int
	CheckboxSize    float32
	ItemHeight      float32
}

// CheckboxList é um componente de lista com checkboxes (versão genérica)
type CheckboxList[T any] struct {
	ID           string
	Items        []CheckboxListItem[T]
	Config       CheckboxListConfig
	ScrollOffset int
	VisibleStart int
	VisibleEnd   int
	// Campos para navegação e foco
	FocusedIndex int  // Índice do item em foco (-1 se nenhum)
	HasFocus     bool // Se o checkbox list tem foco global
	// Campo para viewport dinâmico
	CurrentViewportHeight float32 // Altura atual do viewport para cálculo de scroll
}

// NewCheckboxList cria uma nova lista com checkboxes
func NewCheckboxList[T any](id string, items []CheckboxListItem[T], config CheckboxListConfig) *CheckboxList[T] {
	return &CheckboxList[T]{
		ID:           id,
		Items:        items,
		Config:       config,
		ScrollOffset: 0,
		VisibleStart: 0,
		VisibleEnd:   0,
		FocusedIndex: -1,    // Nenhum item em foco inicialmente
		HasFocus:     false, // Sem foco inicialmente
	}
}

// GetSelectedItems retorna os itens selecionados
func (cl *CheckboxList[T]) GetSelectedItems() []CheckboxListItem[T] {
	var selected []CheckboxListItem[T]
	for _, item := range cl.Items {
		if item.Selected {
			selected = append(selected, item)
		}
	}
	return selected
}

// GetSelectedValues retorna apenas os valores dos itens selecionados
func (cl *CheckboxList[T]) GetSelectedValues() []T {
	var values []T
	for _, item := range cl.Items {
		if item.Selected {
			values = append(values, item.Value)
		}
	}
	return values
}

// ToggleItem alterna o estado de seleção de um item
func (cl *CheckboxList[T]) ToggleItem(index int) {
	if index >= 0 && index < len(cl.Items) {
		cl.Items[index].Selected = !cl.Items[index].Selected
	}
}

// ScrollUp rola a lista para cima
func (cl *CheckboxList[T]) ScrollUp() {
	if cl.ScrollOffset > 0 {
		cl.ScrollOffset--
	}
}

// ScrollDown rola a lista para baixo
func (cl *CheckboxList[T]) ScrollDown() {
	maxVisibleItems := cl.GetMaxVisibleItems()
	if cl.ScrollOffset+maxVisibleItems < len(cl.Items) {
		cl.ScrollOffset++
	}
}

// GetMaxVisibleItems calcula itens visíveis dinamicamente baseado no viewport atual
func (cl *CheckboxList[T]) GetMaxVisibleItems() int {
	if cl.CurrentViewportHeight > 0 && cl.Config.ItemHeight > 0 {
		totalPadding := float32(cl.Config.Padding.Top + cl.Config.Padding.Bottom)
		availableHeight := cl.CurrentViewportHeight - totalPadding
		return int(availableHeight / cl.Config.ItemHeight)
	}
	return 10 // Fallback padrão quando viewport não está definido
}

// SetVisibleItemsCount define número de itens visíveis sem usar altura fixa
func (cl *CheckboxList[T]) SetVisibleItemsCount(count int) {
	// Recalcula MaxHeight baseado no número de itens desejados
	cl.Config.MaxHeight = float32(count) * cl.Config.ItemHeight
}

// GetVisibleRange retorna o range atual de itens visíveis
func (cl *CheckboxList[T]) GetVisibleRange() (start, end int) {
	maxVisible := cl.GetMaxVisibleItems()
	start = cl.ScrollOffset
	end = min(start+maxVisible, len(cl.Items))
	return start, end
}

// EnsureItemVisible garante que um item específico esteja visível
func (cl *CheckboxList[T]) EnsureItemVisible(itemIndex int) {
	if itemIndex < 0 || itemIndex >= len(cl.Items) {
		return
	}

	maxVisible := cl.GetMaxVisibleItems()
	start, end := cl.GetVisibleRange()

	if itemIndex < start {
		cl.ScrollOffset = itemIndex
	} else if itemIndex >= end {
		cl.ScrollOffset = max(0, itemIndex-maxVisible+1)
	}
}

// SetFocused define se o checkbox list tem foco
func (cl *CheckboxList[T]) SetFocused(focused bool) {
	cl.HasFocus = focused
	if focused && cl.FocusedIndex == -1 && len(cl.Items) > 0 {
		cl.FocusedIndex = cl.ScrollOffset
	}
}

// MoveFocusUp move o foco para o item anterior
func (cl *CheckboxList[T]) MoveFocusUp() bool {
	if !cl.HasFocus || len(cl.Items) == 0 {
		return false
	}

	if cl.FocusedIndex > 0 {
		cl.FocusedIndex--
		// Ajustar scroll se item focado saiu da área visível (acima)
		if cl.FocusedIndex < cl.ScrollOffset {
			cl.ScrollOffset = cl.FocusedIndex
		}
		return true
	}
	return false
}

// MoveFocusDown move o foco para o próximo item
func (cl *CheckboxList[T]) MoveFocusDown() bool {
	if !cl.HasFocus || len(cl.Items) == 0 {
		return false
	}

	if cl.FocusedIndex < len(cl.Items)-1 {
		cl.FocusedIndex++
		// Calcular itens máximos usando método dinâmico
		maxVisibleItems := cl.GetMaxVisibleItems()
		if maxVisibleItems == 0 {
			maxVisibleItems = 10 // Fallback para evitar divisão por zero
		}

		// Se o item focado saiu da área visível (abaixo), fazer scroll
		if cl.FocusedIndex >= cl.ScrollOffset+maxVisibleItems {
			cl.ScrollOffset = cl.FocusedIndex - maxVisibleItems + 1
		}
		return true
	}
	return false
}

// ToggleFocusedItem alterna a seleção do item em foco
func (cl *CheckboxList[T]) ToggleFocusedItem() {
	if !cl.HasFocus || cl.FocusedIndex < 0 || cl.FocusedIndex >= len(cl.Items) {
		return
	}
	cl.ToggleItem(cl.FocusedIndex)
}

// calculateViewportMetrics calcula as métricas do viewport para scroll virtual
func (cl *CheckboxList[T]) calculateViewportMetrics(parentHeight float32) (maxVisibleItems int, actualVisibleStart int, actualVisibleEnd int) {
	// Armazenar altura do viewport para uso em navegação
	cl.CurrentViewportHeight = parentHeight

	// Calcular altura disponível para itens (descontando padding do CheckboxList)
	totalPadding := float32(cl.Config.Padding.Top + cl.Config.Padding.Bottom)
	availableHeight := parentHeight - totalPadding

	// Calcular quantos itens cabem no viewport disponível
	maxVisibleItems = max(int(availableHeight/cl.Config.ItemHeight), 1)

	if len(cl.Items) <= maxVisibleItems {
		// Todos os itens cabem no viewport - mostrar todos
		actualVisibleStart = 0
		actualVisibleEnd = len(cl.Items)
		cl.ScrollOffset = 0 // Reset scroll quando todos cabem
	} else {
		// Scroll virtual necessário - viewport menor que total de itens
		actualVisibleStart = cl.ScrollOffset
		actualVisibleEnd = min(cl.ScrollOffset+maxVisibleItems, len(cl.Items))

		// Log quando scroll está ativo
		if cl.ScrollOffset > 0 || actualVisibleEnd < len(cl.Items) {
			log.Printf("Scroll: showing items %d-%d of %d total",
				actualVisibleStart, actualVisibleEnd-1, len(cl.Items))
		}

		// Auto-scroll para manter item focado sempre visível no viewport
		actualVisibleStart, actualVisibleEnd = cl.adjustScrollForFocus(maxVisibleItems, actualVisibleStart, actualVisibleEnd)

		// Ajustar scroll se necessário para evitar espaço vazio
		actualVisibleStart, actualVisibleEnd = cl.adjustScrollToAvoidEmptySpace(maxVisibleItems, actualVisibleStart, actualVisibleEnd)
	}

	return maxVisibleItems, actualVisibleStart, actualVisibleEnd
}

// adjustScrollForFocus ajusta o scroll para manter o item focado visível
func (cl *CheckboxList[T]) adjustScrollForFocus(maxVisibleItems, actualVisibleStart, actualVisibleEnd int) (int, int) {
	if cl.HasFocus && cl.FocusedIndex >= 0 {
		if cl.FocusedIndex < actualVisibleStart {
			// Item focado está acima do viewport - scroll para cima
			cl.ScrollOffset = cl.FocusedIndex
			actualVisibleStart = cl.ScrollOffset
			actualVisibleEnd = min(actualVisibleStart+maxVisibleItems, len(cl.Items))
		} else if cl.FocusedIndex >= actualVisibleEnd {
			// Item focado está abaixo do viewport - scroll para baixo
			cl.ScrollOffset = max(cl.FocusedIndex-maxVisibleItems+1, 0)
			actualVisibleStart = cl.ScrollOffset
			actualVisibleEnd = min(actualVisibleStart+maxVisibleItems, len(cl.Items))
		}
	}
	return actualVisibleStart, actualVisibleEnd
}

// adjustScrollToAvoidEmptySpace ajusta o scroll para evitar espaço vazio no final
func (cl *CheckboxList[T]) adjustScrollToAvoidEmptySpace(maxVisibleItems, actualVisibleStart, actualVisibleEnd int) (int, int) {
	if actualVisibleEnd < len(cl.Items) && actualVisibleEnd-actualVisibleStart < maxVisibleItems {
		actualVisibleStart = max(0, len(cl.Items)-maxVisibleItems)
		actualVisibleEnd = len(cl.Items)
		cl.ScrollOffset = actualVisibleStart
	}
	return actualVisibleStart, actualVisibleEnd
}

// getItemBackgroundColor determina a cor de fundo de um item baseado no seu estado
func (cl *CheckboxList[T]) getItemBackgroundColor(index int) clay.Color {
	if cl.HasFocus && cl.FocusedIndex == index {
		return clay.Color{R: 60, G: 120, B: 200, A: 255} // Azul para item focado
	} else if cl.Items[index].Selected {
		return clay.Color{R: 40, G: 120, B: 80, A: 255} // Verde para item selecionado
	}
	return clay.Color{R: 0, G: 0, B: 0, A: 0} // Transparente para item normal
}

// getLabelColor determina a cor do texto do label baseado no estado do item
func (cl *CheckboxList[T]) getLabelColor(index int) clay.Color {
	if cl.HasFocus && cl.FocusedIndex == index {
		return clay.Color{R: 255, G: 255, B: 255, A: 255} // Branco para item focado
	} else if cl.Items[index].Selected {
		return clay.Color{R: 180, G: 255, B: 200, A: 255} // Verde claro para item selecionado
	}
	return clay.Color{R: 220, G: 230, B: 245, A: 255} // Cinza claro para item normal
}

// renderCheckbox renderiza o checkbox de um item
func (cl *CheckboxList[T]) renderCheckbox(claySystem *ClayLayoutSystem, itemIndex int) {
	item := cl.Items[itemIndex]
	checkboxID := fmt.Sprintf("%s-checkbox-%d", cl.ID, itemIndex)

	checkboxColor := clay.Color{R: 60, G: 70, B: 85, A: 255} // Cinza escuro moderno
	if item.Selected {
		checkboxColor = clay.Color{R: 30, G: 180, B: 120, A: 255} // Verde moderno vibrante
	}

	clay.UI()(clay.ElementDeclaration{
		Id: clay.ID(checkboxID),
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  clay.SizingFixed(cl.Config.CheckboxSize),
				Height: clay.SizingFixed(cl.Config.CheckboxSize),
			},
			ChildAlignment: clay.ChildAlignment{
				X: clay.ALIGN_X_CENTER,
				Y: clay.ALIGN_Y_CENTER,
			},
		},
		CornerRadius:    clay.CornerRadiusAll(4), // Checkbox ligeiramente arredondado
		BackgroundColor: checkboxColor,
	}, func() {
		if item.Selected {
			claySystem.CreateText("◾", TextConfig{
				FontSize:  25,
				TextColor: clay.Color{R: 255, G: 255, B: 255, A: 255}, // Branco puro
			})
		}
	})
}

// renderLabel renderiza o label de um item
func (cl *CheckboxList[T]) renderLabel(claySystem *ClayLayoutSystem, itemIndex int) {
	item := cl.Items[itemIndex]
	labelColor := cl.getLabelColor(itemIndex)
	labelContainerID := fmt.Sprintf("%s-label-%d", cl.ID, itemIndex)

	clay.UI()(clay.ElementDeclaration{
		Id: clay.ID(labelContainerID),
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  clay.SizingPercent(1.0),
				Height: clay.SizingPercent(1.0),
			},
			LayoutDirection: clay.TOP_TO_BOTTOM,
			ChildAlignment: clay.ChildAlignment{
				Y: clay.ALIGN_Y_CENTER,
			},
		},
	}, func() {
		claySystem.CreateText(item.Label, TextConfig{
			FontSize:  15,
			TextColor: labelColor,
		})
	})
}

// renderItem renderiza um único item da lista
func (cl *CheckboxList[T]) renderItem(claySystem *ClayLayoutSystem, itemIndex int) {
	itemID := fmt.Sprintf("%s-item-%d", cl.ID, itemIndex)
	itemBgColor := cl.getItemBackgroundColor(itemIndex)

	// Checkbox item container
	clay.UI()(clay.ElementDeclaration{
		Id: clay.ID(itemID),
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width: clay.SizingPercent(1.0),
			},
			Padding:         clay.PaddingAll(12),
			ChildGap:        12,
			LayoutDirection: clay.LEFT_TO_RIGHT,
		},
		CornerRadius:    clay.CornerRadiusAll(10), // Bordas arredondadas nos itens
		BackgroundColor: itemBgColor,
	}, func() {
		cl.renderCheckbox(claySystem, itemIndex)
		cl.renderLabel(claySystem, itemIndex)
	})
}

// updateVisiblePositions atualiza as posições visíveis usando o método correto
func (cl *CheckboxList[T]) updateVisiblePositions() {
	maxItems := cl.GetMaxVisibleItems()
	cl.VisibleStart = cl.ScrollOffset
	cl.VisibleEnd = min(cl.ScrollOffset+maxItems, len(cl.Items))
}

// RenderCheckboxList renderiza uma lista com viewport dinâmico baseado na altura do container pai
func (cl *CheckboxList[T]) RenderCheckboxList(claySystem *ClayLayoutSystem, parentHeight float32) {
	_, actualVisibleStart, actualVisibleEnd := cl.calculateViewportMetrics(parentHeight)

	// Container principal da lista usando toda altura disponível do pai (viewport)
	clay.UI()(clay.ElementDeclaration{
		Id: clay.ID(cl.ID),
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width: cl.Config.Sizing.Width,
				// Height: clay.SizingPercent(1.0),
			},
			Padding:         cl.Config.Padding,
			ChildGap:        cl.Config.ChildGap,
			LayoutDirection: clay.TOP_TO_BOTTOM,
		},
		CornerRadius:    clay.CornerRadiusAll(12),
		BackgroundColor: cl.Config.BackgroundColor,
	}, func() {
		for i := actualVisibleStart; i < actualVisibleEnd; i++ {
			cl.renderItem(claySystem, i)
		}
	})

	cl.updateVisiblePositions()

	log.Printf("Checkbox list created successfully: %s", cl.ID)
}

// DefaultCheckboxListConfig retorna configuração para viewport dinâmico
func DefaultCheckboxListConfig() CheckboxListConfig {
	return CheckboxListConfig{
		Sizing: clay.Sizing{
			Width: clay.SizingPercent(1.0),
		},
		Padding:         clay.PaddingAll(10),
		ChildGap:        5,
		BackgroundColor: clay.Color{R: 25, G: 30, B: 40, A: 240},
		ScrollOffset:    0,
		CheckboxSize:    25,
		ItemHeight:      45,
	}
}

// ViewportCheckboxListConfig configuração otimizada para viewport dinâmico
func ViewportCheckboxListConfig() CheckboxListConfig {
	config := DefaultCheckboxListConfig()
	config.ItemHeight = 40   // Itens um pouco menores para mais itens visíveis
	config.CheckboxSize = 22 // Checkbox proporcionalmente menor
	config.ChildGap = 3      // Gap menor para aproveitar melhor o espaço
	return config
}

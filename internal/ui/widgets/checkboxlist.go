package widgets

import (
	"fmt"
	"log"
	"retroart-sdl2/internal/input"
	"retroart-sdl2/internal/theme"

	"github.com/TotallyGamerJet/clay"
)

// CheckboxListItem é um item da lista com tipo genérico para Value
type CheckboxListItem[T any] struct {
	Label    string
	Value    T
	Selected bool
}

type CheckboxListConfig = theme.CheckboxListStyle

// CheckboxList é um componente de lista com checkboxes (versão genérica)
type CheckboxList[T any] struct {
	ID           string
	Items        []CheckboxListItem[T]
	Config       CheckboxListConfig
	ScrollOffset int
	VisibleStart int
	VisibleEnd   int
	FocusedIndex int
	HasFocus     bool
	Height       float32
}

// NewCheckboxList cria uma nova lista usando o design system
func NewCheckboxList[T any](id string, items []CheckboxListItem[T]) *CheckboxList[T] {
	return &CheckboxList[T]{
		ID:           id,
		Items:        items,
		Config:       theme.GetCheckboxListStyle(),
		ScrollOffset: 0,
		FocusedIndex: -1, // Nenhum item em foco inicialmente
		HasFocus:     false,
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
	if cl.Height > 0 && cl.Config.ItemHeight > 0 {
		totalPadding := float32(cl.Config.Padding.Top + cl.Config.Padding.Bottom)
		availableHeight := cl.Height - totalPadding
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
		// Calculate max items using dynamic method
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

// calculateHeight calcula as métricas do viewport para scroll virtual
func (cl *CheckboxList[T]) calculateHeight(parentHeight float32) (maxVisibleItems int, actualVisibleStart int, actualVisibleEnd int) {
	// Armazenar altura do viewport para uso em navegação
	cl.Height = parentHeight

	// Calculate available height for items (excluding CheckboxList padding)
	totalPadding := float32(cl.Config.Padding.Top + cl.Config.Padding.Bottom)
	availableHeight := parentHeight - totalPadding

	// Calculate how many items fit in the available viewport
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
		return cl.Config.ItemFocusedBg
	} else if cl.Items[index].Selected {
		return cl.Config.ItemSelectedBg
	}
	return cl.Config.ItemNormalBg
}

// getLabelColor determina a cor do texto do label baseado no estado do item
func (cl *CheckboxList[T]) getLabelColor(index int) clay.Color {
	if cl.HasFocus && cl.FocusedIndex == index {
		return cl.Config.ItemFocusedText
	} else if cl.Items[index].Selected {
		return cl.Config.ItemSelectedText
	}
	return cl.Config.ItemNormalText
}

// renderCheckbox renderiza o checkbox de um item
func (cl *CheckboxList[T]) renderCheckbox(itemIndex int) {
	item := cl.Items[itemIndex]
	checkboxID := fmt.Sprintf("%s-checkbox-%d", cl.ID, itemIndex)

	checkboxColor := cl.Config.CheckboxNormal
	if item.Selected {
		checkboxColor = cl.Config.CheckboxSelected
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
			clay.Text("◣", &clay.TextElementConfig{
				FontSize:  14,
				TextColor: cl.Config.CheckboxMark,
			})
		}
	})
}

// renderLabel renderiza o label de um item
func (cl *CheckboxList[T]) renderLabel(itemIndex int) {
	item := cl.Items[itemIndex]
	labelColor := cl.getLabelColor(itemIndex)
	labelContainerID := fmt.Sprintf("%s-label-%d", cl.ID, itemIndex)

	clay.UI()(clay.ElementDeclaration{
		Id: clay.ID(labelContainerID),
		Layout: clay.LayoutConfig{
			LayoutDirection: clay.TOP_TO_BOTTOM,
			ChildAlignment: clay.ChildAlignment{
				Y: clay.ALIGN_Y_CENTER,
			},
		},
	}, func() {
		clay.Text(item.Label, &clay.TextElementConfig{
			FontSize:  14,
			TextColor: labelColor,
		})
	})
}

// renderItem renderiza um único item da lista
func (cl *CheckboxList[T]) renderItem(itemIndex int) {
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
		cl.renderCheckbox(itemIndex)
		cl.renderLabel(itemIndex)
	})
}

// updateVisiblePositions atualiza as posições visíveis usando o método correto
func (cl *CheckboxList[T]) updateVisiblePositions() {
	maxItems := cl.GetMaxVisibleItems()
	cl.VisibleStart = cl.ScrollOffset
	cl.VisibleEnd = min(cl.ScrollOffset+maxItems, len(cl.Items))
}

// Render renderiza uma lista com viewport dinâmico baseado na altura do container pai
func (cl *CheckboxList[T]) Render(height float32) {
	_, actualVisibleStart, actualVisibleEnd := cl.calculateHeight(height)

	// Container principal da lista
	clay.UI()(clay.ElementDeclaration{
		Id: clay.ID(cl.ID),
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width: clay.SizingPercent(1.0),
				// Height: clay.SizingPercent(0.9),
			},
			Padding:         cl.Config.Padding,
			ChildGap:        cl.Config.ChildGap,
			LayoutDirection: clay.TOP_TO_BOTTOM,
		},
		CornerRadius:    clay.CornerRadiusAll(12),
		BackgroundColor: cl.Config.BackgroundColor,
	}, func() {
		for i := actualVisibleStart; i < actualVisibleEnd; i++ {
			cl.renderItem(i)
		}
	})

	cl.updateVisiblePositions()

	log.Printf("Checkbox list created successfully: %s", cl.ID)
}

// Interface Focusable implementation
func (cl *CheckboxList[T]) GetID() string {
	return cl.ID
}

func (cl *CheckboxList[T]) IsFocused() bool {
	return cl.HasFocus
}

func (cl *CheckboxList[T]) OnFocusChanged(focused bool) {
	cl.HasFocus = focused
	if focused && cl.FocusedIndex == -1 && len(cl.Items) > 0 {
		cl.FocusedIndex = cl.ScrollOffset
	}
}

func (cl *CheckboxList[T]) CanFocus() bool {
	return len(cl.Items) > 0
}

func (cl *CheckboxList[T]) HandleInput(inputType input.InputType) bool {
	if !cl.HasFocus {
		return false
	}

	switch inputType {
	case input.InputUp:
		return cl.MoveFocusUp()
	case input.InputDown:
		return cl.MoveFocusDown()
	case input.InputConfirm:
		cl.ToggleFocusedItem()
		return true
	default:
		return false
	}
}

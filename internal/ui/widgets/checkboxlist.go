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
	Width        clay.SizingAxis
	Height       clay.SizingAxis
	listHeight   float32
	itemHeight   float32
}

// NewCheckboxList cria uma nova lista usando o design system
func NewCheckboxList[T any](id string, width, height clay.SizingAxis, items []CheckboxListItem[T]) *CheckboxList[T] {
	return &CheckboxList[T]{
		ID:           id,
		Items:        items,
		Config:       theme.GetCheckboxListStyle(),
		ScrollOffset: 0,
		Width:        width,
		Height:       height,
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

// GetMaxVisibleItems calcula itens visíveis dinamicamente baseado no viewport atual
func (cl *CheckboxList[T]) GetMaxVisibleItems() int {
	if cl.listHeight > 0 && cl.itemHeight > 0 {
		r := int(cl.listHeight / (cl.itemHeight + float32(cl.Config.ChildGap)))
		return r
	}
	return 12
}

func (cl *CheckboxList[T]) GetVisibleItemsRange() (start, end int) {
	maxVisible := cl.GetMaxVisibleItems()
	start = cl.ScrollOffset
	end = min(start+maxVisible, len(cl.Items))
	fmt.Printf("Visible items range: start=%d, end=%d, total items=%d\n", start, end, len(cl.Items))
	return start, end
}

// ScrollUp move o foco para o item anterior
func (cl *CheckboxList[T]) ScrollUp() bool {
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

// ScrollDown move o foco para o próximo item
func (cl *CheckboxList[T]) ScrollDown() bool {
	if !cl.HasFocus || len(cl.Items) == 0 {
		return false
	}

	if cl.FocusedIndex < len(cl.Items)-1 {
		cl.FocusedIndex++
		// Calculate max items using dynamic method
		maxVisibleItems := cl.GetMaxVisibleItems()

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
		Id:          clay.ID(checkboxID),
		AspectRatio: clay.AspectRatioElementConfig{AspectRatio: 1.0},
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  clay.SizingFixed(cl.Config.CheckboxSize),
				Height: clay.SizingFixed(cl.Config.CheckboxSize),
			},
			ChildAlignment: clay.ChildAlignment{
				X: clay.ALIGN_X_CENTER,
				Y: clay.ALIGN_Y_CENTER,
			},
			Padding: cl.Config.Padding,
		},
		CornerRadius:    clay.CornerRadiusAll(cl.Config.CornerRadius), // Checkbox ligeiramente arredondado
		BackgroundColor: checkboxColor,
	}, func() {
		if item.Selected {
			Text("◣", 16, cl.Config.CheckboxMark)
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
				X: clay.ALIGN_X_CENTER,
				Y: clay.ALIGN_Y_CENTER,
			},
		},
	}, func() {
		Text(item.Label, 14, labelColor)
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
				Width: clay.SizingGrow(0),
			},
			Padding:         cl.Config.Padding,
			ChildGap:        cl.Config.ChildGap,
			LayoutDirection: clay.LEFT_TO_RIGHT,
			ChildAlignment: clay.ChildAlignment{
				X: clay.ALIGN_X_LEFT,
				Y: clay.ALIGN_Y_CENTER,
			},
		},
		CornerRadius:    clay.CornerRadiusAll(cl.Config.CornerRadius), // Bordas arredondadas nos itens
		BackgroundColor: itemBgColor,
	}, func() {
		parentData := clay.GetElementData(clay.ID(itemID))
		cl.itemHeight = parentData.BoundingBox.Height
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
func (cl *CheckboxList[T]) Render() {
	clay.UI()(clay.ElementDeclaration{
		Id: clay.ID(cl.ID),
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  cl.Width,
				Height: cl.Height,
			},
			ChildAlignment: clay.ChildAlignment{
				X: clay.ALIGN_X_CENTER,
				Y: clay.ALIGN_Y_CENTER,
			},
		},
		CornerRadius:    clay.CornerRadiusAll(cl.Config.CornerRadius),
		BackgroundColor: cl.Config.BackgroundColor,
	}, func() {
		clay.UI()(clay.ElementDeclaration{
			Id: clay.ID("item-holder"),
			Layout: clay.LayoutConfig{
				Sizing: clay.Sizing{
					Width: clay.SizingGrow(0),
				},
				Padding:         cl.Config.Padding,
				ChildGap:        cl.Config.ChildGap,
				LayoutDirection: clay.TOP_TO_BOTTOM,
			},
		}, func() {
			parentData := clay.GetElementData(clay.ID(cl.ID))
			cl.listHeight = parentData.BoundingBox.Height
			start, end := cl.GetVisibleItemsRange()
			for i := start; i < end; i++ {
				cl.renderItem(i)
			}
		})
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
		return cl.ScrollUp()
	case input.InputDown:
		return cl.ScrollDown()
	case input.InputConfirm:
		cl.ToggleFocusedItem()
		return true
	default:
		return false
	}
}

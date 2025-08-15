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
	ID               string
	Items            []CheckboxListItem[T]
	Config           CheckboxListConfig
	ScrollOffset     int
	VisibleStart     int
	VisibleEnd       int
	FocusedIndex     int
	HasFocus         bool
	Width            clay.SizingAxis
	Height           clay.SizingAxis
	listHeight       float32
	itemHeight       float32
	scrollDownHeight float32
	scrollUpHeight   float32
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

// Render renders the CheckboxList widget, including its layout, background, and child items.
// It sets up the main container with the specified sizing, padding, corner radius, and background color,
// then creates a vertically arranged holder for the list items. The method calculates the visible range
// of items, renders scroll controls if necessary, and displays each visible item in the list.
// A log message is printed upon successful creation of the checkbox list.
func (cl *CheckboxList[T]) Render() {
	clay.UI()(clay.ElementDeclaration{
		Id: clay.ID(cl.ID),
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  cl.Width,
				Height: cl.Height,
			},
			Padding: cl.Config.Padding,
			ChildAlignment: clay.ChildAlignment{
				X: clay.ALIGN_X_CENTER,
				Y: clay.ALIGN_Y_CENTER,
			},
		},
		CornerRadius:    clay.CornerRadiusAll(cl.Config.CornerRadius),
		BackgroundColor: cl.Config.BackgroundColor,
	}, func() {
		parentData := clay.GetElementData(clay.ID(cl.ID))
		cl.listHeight = parentData.BoundingBox.Height
		start, end := cl.getVisibleItemsRange()

		clay.UI()(clay.ElementDeclaration{
			Id: clay.ID("items-holder"),
			Layout: clay.LayoutConfig{
				Sizing: clay.Sizing{
					Width: clay.SizingGrow(0),
				},
				ChildGap:        cl.Config.ChildGap,
				LayoutDirection: clay.TOP_TO_BOTTOM,
			},
		}, func() {
			for i := start; i < end; i++ {
				cl.renderItem(i)
			}
		})
		cl.renderScrollBar(start, end)
	})

	log.Printf("CheckboxList: render id=%s items=%d", cl.ID, len(cl.Items))
}

func (cl *CheckboxList[T]) renderScrollBar(start, end int) {
	clay.UI()(clay.ElementDeclaration{
		Id: clay.ID(cl.ID + "-scroll"),
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  clay.SizingFixed(20),
				Height: clay.SizingGrow(0),
			},
			ChildAlignment: clay.ChildAlignment{
				X: clay.ALIGN_X_CENTER,
			},
			LayoutDirection: clay.TOP_TO_BOTTOM,
		},
	}, func() {
		cl.renderScrollUp(start)
		cl.renderScrollDown(end)
	})
}

func (cl *CheckboxList[T]) renderScrollUp(start int) {
	scrollUpId := clay.ID(cl.ID + "-scroll-up")
	clay.UI()(clay.ElementDeclaration{
		Id: scrollUpId,
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  clay.SizingGrow(0),
				Height: clay.SizingGrow(0),
			},
			ChildAlignment: clay.ChildAlignment{
				X: clay.ALIGN_X_RIGHT,
				Y: clay.ALIGN_Y_TOP,
			},
		},
		CornerRadius: clay.CornerRadiusAll(cl.Config.Checkbox.CornerRadius),
	}, func() {
		cl.scrollUpHeight = clay.GetElementData(scrollUpId).BoundingBox.Height
		indicatorColor := cl.Config.Checkbox.Color.Normal
		if start > 0 {
			indicatorColor = cl.Config.Checkbox.Color.Mark
		}
		Text(cl.Config.Checkbox.ScrollIndicator.UpSymbol, cl.Config.Checkbox.ScrollIndicator.Size, indicatorColor)
	})
}

func (cl *CheckboxList[T]) renderScrollDown(end int) {
	scrollDownId := clay.ID(cl.ID + "-scroll-down")
	clay.UI()(clay.ElementDeclaration{
		Id: scrollDownId,
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  clay.SizingGrow(0),
				Height: clay.SizingGrow(0),
			},
			ChildAlignment: clay.ChildAlignment{
				X: clay.ALIGN_X_RIGHT,
				Y: clay.ALIGN_Y_BOTTOM,
			},
		},
		CornerRadius: clay.CornerRadiusAll(cl.Config.Checkbox.CornerRadius),
	}, func() {
		cl.scrollDownHeight = clay.GetElementData(scrollDownId).BoundingBox.Height
		var indicatorColor = cl.Config.Checkbox.Color.Normal
		if end < len(cl.Items) {
			indicatorColor = cl.Config.Checkbox.Color.Mark
		}
		Text(cl.Config.Checkbox.ScrollIndicator.DownSymbol, cl.Config.Checkbox.ScrollIndicator.Size, indicatorColor)
	})
}

func (cl *CheckboxList[T]) renderItem(itemIndex int) {
	itemID := fmt.Sprintf("%s-item-%d", cl.ID, itemIndex)

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
		BackgroundColor: cl.getItemBackgroundColor(itemIndex),
	}, func() {
		parentData := clay.GetElementData(clay.ID(itemID))
		cl.itemHeight = parentData.BoundingBox.Height
		cl.renderCheckbox(itemIndex)
		cl.renderLabel(itemIndex)
	})
}

func (cl *CheckboxList[T]) renderCheckbox(itemIndex int) {
	item := cl.Items[itemIndex]
	checkboxID := fmt.Sprintf("%s-checkbox-%d", cl.ID, itemIndex)

	checkboxColor := cl.Config.Checkbox.Color.Normal
	if item.Selected {
		checkboxColor = cl.Config.Checkbox.Color.Selected
	}

	clay.UI()(clay.ElementDeclaration{
		Id: clay.ID(checkboxID),
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  clay.SizingFixed(cl.Config.Checkbox.Size),
				Height: clay.SizingFixed(cl.Config.Checkbox.Size),
			},
			ChildAlignment: clay.ChildAlignment{
				X: clay.ALIGN_X_CENTER,
				Y: clay.ALIGN_Y_CENTER,
			},
			Padding: cl.Config.Padding,
		},
		CornerRadius:    clay.CornerRadiusAll(cl.Config.Checkbox.CornerRadius),
		BackgroundColor: checkboxColor,
	}, func() {
		if item.Selected {
			Text(
				cl.Config.Checkbox.Mark.Symbol,
				cl.Config.Checkbox.Mark.Size,
				cl.Config.Checkbox.Color.Mark,
			)
		}
	})
}

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
		Text(item.Label, cl.Config.FontSize, labelColor)
	})
}

func (cl *CheckboxList[T]) getItemBackgroundColor(index int) clay.Color {
	if cl.HasFocus && cl.FocusedIndex == index {
		return cl.Config.ItemFocusedBg
	} else if cl.Items[index].Selected {
		return cl.Config.ItemSelectedBg
	}
	return cl.Config.ItemNormalBg
}

func (cl *CheckboxList[T]) getLabelColor(index int) clay.Color {
	if cl.HasFocus && cl.FocusedIndex == index {
		return cl.Config.ItemFocusedText
	} else if cl.Items[index].Selected {
		return cl.Config.ItemSelectedText
	}
	return cl.Config.ItemNormalText
}

func (cl *CheckboxList[T]) getMaxVisibleItems() int {
	// r := int((cl.listHeight - (cl.scrollDownHeight + cl.scrollUpHeight + float32(cl.Config.Padding.Top) + float32(cl.Config.Padding.Bottom))) / (cl.itemHeight + float32(cl.Config.ChildGap)))
	r := int((cl.listHeight) / (cl.itemHeight + float32(cl.Config.ChildGap)))
	if r > 0 {
		return r
	}
	return 10
}

func (cl *CheckboxList[T]) getVisibleItemsRange() (start, end int) {
	maxVisible := cl.getMaxVisibleItems() // -1 para considerar o item focado
	start = cl.ScrollOffset
	end = min(start+maxVisible, len(cl.Items))
	return start, end
}

// GetSelectedItems returns a slice of CheckboxListItem[T] containing all items
// from the CheckboxList that are currently selected. Only items with the Selected
// field set to true are included in the returned slice.
func (cl *CheckboxList[T]) GetSelectedItems() []CheckboxListItem[T] {
	var selected []CheckboxListItem[T]
	for _, item := range cl.Items {
		if item.Selected {
			selected = append(selected, item)
		}
	}
	return selected
}

// GetSelectedValues returns a slice containing the values of all items in the CheckboxList
// that are currently selected. The returned slice contains values of type T corresponding
// to each selected item.
func (cl *CheckboxList[T]) GetSelectedValues() []T {
	var values []T
	for _, item := range cl.Items {
		if item.Selected {
			values = append(values, item.Value)
		}
	}
	return values
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

// ScrollDown moves the focus down by one item in the CheckboxList.
// If the list has focus and is not empty, it increments the FocusedIndex,
// and if the newly focused item is outside the visible area, it adjusts
// the ScrollOffset to bring the item into view. Returns true if the focus
// was moved, or false if already at the last item or the list is unfocused/empty.
func (cl *CheckboxList[T]) ScrollDown() bool {
	if !cl.HasFocus || len(cl.Items) == 0 {
		return false
	}

	if cl.FocusedIndex < len(cl.Items)-1 {
		cl.FocusedIndex++
		// Calculate max items using dynamic method
		maxVisibleItems := cl.getMaxVisibleItems()

		// Se o item focado saiu da área visível (abaixo), fazer scroll
		if cl.FocusedIndex >= cl.ScrollOffset+maxVisibleItems {
			cl.ScrollOffset = cl.FocusedIndex - maxVisibleItems + 1
		}
		return true
	}
	return false
}

func (cl *CheckboxList[T]) toggleFocusedItem() {
	if !cl.HasFocus || cl.FocusedIndex < 0 || cl.FocusedIndex >= len(cl.Items) {
		return
	}
	if cl.FocusedIndex >= 0 && cl.FocusedIndex < len(cl.Items) {
		cl.Items[cl.FocusedIndex].Selected = !cl.Items[cl.FocusedIndex].Selected
	}
}

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
		cl.toggleFocusedItem()
		return true
	default:
		return false
	}
}

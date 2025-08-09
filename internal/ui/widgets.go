package ui

import (
	"fmt"
	"log"

	"github.com/TotallyGamerJet/clay"
)

// Estruturas para CheckboxList (versão genérica)

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
	maxItems := int(cl.Config.MaxHeight / cl.Config.ItemHeight)
	if cl.ScrollOffset+maxItems < len(cl.Items) {
		cl.ScrollOffset++
	}
}

// SetFocused define se o checkbox list tem foco
func (cl *CheckboxList[T]) SetFocused(focused bool) {
	cl.HasFocus = focused
	if focused && cl.FocusedIndex == -1 && len(cl.Items) > 0 {
		// Se ganhou foco e nenhum item estava focado, focar no primeiro visível
		cl.FocusedIndex = cl.ScrollOffset
	}
}

// MoveFocusUp move o foco para o item anterior
func (cl *CheckboxList[T]) MoveFocusUp() {
	if !cl.HasFocus || len(cl.Items) == 0 {
		return
	}

	if cl.FocusedIndex > 0 {
		cl.FocusedIndex--
		// Se o item focado saiu da área visível, fazer scroll
		if cl.FocusedIndex < cl.ScrollOffset {
			cl.ScrollOffset = cl.FocusedIndex
		}
	}
}

// MoveFocusDown move o foco para o próximo item
func (cl *CheckboxList[T]) MoveFocusDown() {
	if !cl.HasFocus || len(cl.Items) == 0 {
		return
	}

	if cl.FocusedIndex < len(cl.Items)-1 {
		cl.FocusedIndex++
		// Se o item focado saiu da área visível, fazer scroll
		maxItems := int(cl.Config.MaxHeight / cl.Config.ItemHeight)
		if cl.FocusedIndex >= cl.ScrollOffset+maxItems {
			cl.ScrollOffset = cl.FocusedIndex - maxItems + 1
		}
	}
}

// ToggleFocusedItem alterna a seleção do item em foco
func (cl *CheckboxList[T]) ToggleFocusedItem() {
	if !cl.HasFocus || cl.FocusedIndex < 0 || cl.FocusedIndex >= len(cl.Items) {
		return
	}
	cl.ToggleItem(cl.FocusedIndex)
}

// CreateCheckboxList renderiza uma lista com checkboxes
func (cls *ClayLayoutSystem) CreateCheckboxList(checkboxList interface{}) {
	if !cls.enabled || !cls.isActive {
		log.Printf("Clay not enabled or not active, skipping CreateCheckboxList")
		return
	}

	// Usa type assertion para acessar propriedades comuns
	switch cbl := checkboxList.(type) {
	case *CheckboxList[string]:
		renderCheckboxListItems(cls, cbl.ID, cbl.Items, cbl.Config, cbl.HasFocus, cbl.FocusedIndex, cbl.VisibleStart, cbl.VisibleEnd)
		// Atualizar posições visíveis
		maxVisibleItems := int(cbl.Config.MaxHeight / cbl.Config.ItemHeight)
		cbl.VisibleStart = cbl.ScrollOffset
		cbl.VisibleEnd = cbl.ScrollOffset + maxVisibleItems
		if cbl.VisibleEnd > len(cbl.Items) {
			cbl.VisibleEnd = len(cbl.Items)
		}

	case *CheckboxList[int]:
		renderCheckboxListItems(cls, cbl.ID, cbl.Items, cbl.Config, cbl.HasFocus, cbl.FocusedIndex, cbl.VisibleStart, cbl.VisibleEnd)
		// Atualizar posições visíveis
		maxVisibleItems := int(cbl.Config.MaxHeight / cbl.Config.ItemHeight)
		cbl.VisibleStart = cbl.ScrollOffset
		cbl.VisibleEnd = cbl.ScrollOffset + maxVisibleItems
		if cbl.VisibleEnd > len(cbl.Items) {
			cbl.VisibleEnd = len(cbl.Items)
		}

	default:
		log.Printf("Unsupported CheckboxList type: %T", checkboxList)
		return
	}
}

// Função auxiliar genérica para renderizar itens do checkbox list
func renderCheckboxListItems[T any](
	cls *ClayLayoutSystem,
	id string,
	items []CheckboxListItem[T],
	config CheckboxListConfig,
	hasFocus bool,
	focusedIndex int,
	visibleStart int,
	visibleEnd int,
) {
	log.Printf("Creating checkbox list: %s", id)

	// Atualizar posições visíveis baseado no scroll
	maxVisibleItems := int(config.MaxHeight / config.ItemHeight)
	actualVisibleStart := 0
	actualVisibleEnd := len(items)

	if len(items) > maxVisibleItems {
		actualVisibleStart = visibleStart
		actualVisibleEnd = visibleStart + maxVisibleItems
		if actualVisibleEnd > len(items) {
			actualVisibleEnd = len(items)
		}
	}

	// Container principal da lista
	clay.UI()(clay.ElementDeclaration{
		Id: clay.ID(id),
		Layout: clay.LayoutConfig{
			Sizing:          config.Sizing,
			Padding:         config.Padding,
			ChildGap:        config.ChildGap,
			LayoutDirection: clay.TOP_TO_BOTTOM,
		},
		CornerRadius:    clay.CornerRadiusAll(12), // Bordas mais arredondadas
		BackgroundColor: config.BackgroundColor,
	}, func() {
		// Renderizar apenas os itens visíveis
		for i := actualVisibleStart; i < actualVisibleEnd; i++ {
			item := items[i]
			itemID := fmt.Sprintf("%s-item-%d", id, i)

			// Determinar cor de fundo do item (destaque se em foco)
			itemBgColor := clay.Color{R: 0, G: 0, B: 0, A: 0} // Transparente por padrão
			if hasFocus && focusedIndex == i {
				// Item em foco - gradiente azul moderno
				itemBgColor = clay.Color{R: 60, G: 120, B: 200, A: 180}
			} else if item.Selected {
				// Item selecionado - fundo verde suave
				itemBgColor = clay.Color{R: 40, G: 120, B: 80, A: 100}
			}

			// Container do item
			clay.UI()(clay.ElementDeclaration{
				Id: clay.ID(itemID),
				Layout: clay.LayoutConfig{
					Sizing: clay.Sizing{
						Width:  clay.SizingGrow(0),
						Height: clay.SizingFixed(config.ItemHeight),
					},
					Padding:         clay.Padding{Left: 15, Right: 15, Top: 8, Bottom: 8},
					ChildGap:        12,
					LayoutDirection: clay.LEFT_TO_RIGHT,
					ChildAlignment: clay.ChildAlignment{
						Y: clay.ALIGN_Y_CENTER,
					},
				},
				CornerRadius:    clay.CornerRadiusAll(8), // Bordas arredondadas nos itens
				BackgroundColor: itemBgColor,
			}, func() {
				// Checkbox
				checkboxID := fmt.Sprintf("%s-checkbox-%d", id, i)
				checkboxColor := clay.Color{R: 60, G: 70, B: 85, A: 255} // Cinza escuro moderno
				if item.Selected {
					checkboxColor = clay.Color{R: 30, G: 180, B: 120, A: 255} // Verde moderno vibrante
				}

				clay.UI()(clay.ElementDeclaration{
					Id: clay.ID(checkboxID),
					Layout: clay.LayoutConfig{
						Sizing: clay.Sizing{
							Width:  clay.SizingFixed(config.CheckboxSize),
							Height: clay.SizingFixed(config.CheckboxSize),
						},
					},
					CornerRadius:    clay.CornerRadiusAll(4), // Checkbox ligeiramente arredondado
					BackgroundColor: checkboxColor,
				}, func() {
					// Marca de seleção (checkmark) se selecionado
					if item.Selected {
						cls.CreateText("✓", TextConfig{
							FontSize:  14,
							TextColor: clay.Color{R: 255, G: 255, B: 255, A: 255}, // Branco puro
						})
					}
				})

				// Label do item
				labelColor := clay.Color{R: 220, G: 230, B: 245, A: 255} // Branco azulado suave por padrão
				if hasFocus && focusedIndex == i {
					// Texto mais brilhante quando em foco
					labelColor = clay.Color{R: 255, G: 255, B: 255, A: 255} // Branco puro
				} else if item.Selected {
					// Texto verde claro para itens selecionados
					labelColor = clay.Color{R: 180, G: 255, B: 200, A: 255}
				}

				cls.CreateText(item.Label, TextConfig{
					FontSize:  15, // Fonte ligeiramente maior
					TextColor: labelColor,
				})
			})
		}
	})

	log.Printf("Checkbox list created successfully: %s", id)
}

// DefaultCheckboxListConfig retorna uma configuração padrão para CheckboxList
func DefaultCheckboxListConfig() CheckboxListConfig {
	return CheckboxListConfig{
		Sizing: clay.Sizing{
			Width:  clay.SizingFixed(300),
			Height: clay.SizingFixed(250),
		},
		Padding:         clay.PaddingAll(18),
		ChildGap:        8,
		BackgroundColor: clay.Color{R: 25, G: 30, B: 40, A: 240}, // Fundo escuro moderno com transparência
		MaxHeight:       250,
		ScrollOffset:    0,
		CheckboxSize:    18,
		ItemHeight:      35,
	}
}

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

// RenderCheckboxList renderiza uma lista com checkboxes e retorna se foi criada com sucesso
func (cl *CheckboxList[T]) RenderCheckboxList(claySystem *ClayLayoutSystem) {
	log.Printf("Creating checkbox list: %s", cl.ID)

	// Atualizar posições visíveis baseado no scroll
	maxVisibleItems := int(cl.Config.MaxHeight / cl.Config.ItemHeight)
	actualVisibleStart := 0
	actualVisibleEnd := len(cl.Items)

	if len(cl.Items) > maxVisibleItems {
		actualVisibleStart = cl.VisibleStart
		actualVisibleEnd = min(cl.VisibleStart+maxVisibleItems, len(cl.Items))
	}

	// Container principal da lista com dimensões fixas baseadas no container pai
	clay.UI()(clay.ElementDeclaration{
		Id: clay.ID(cl.ID),
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  cl.Config.Sizing.Width,
				Height: cl.Config.Sizing.Height,
			},
			Padding:         cl.Config.Padding,
			ChildGap:        cl.Config.ChildGap,
			LayoutDirection: clay.TOP_TO_BOTTOM,
		},
		CornerRadius:    clay.CornerRadiusAll(12), // Bordas mais arredondadas
		BackgroundColor: cl.Config.BackgroundColor,
	}, func() {
		// Renderizar apenas os itens visíveis
		for i := actualVisibleStart; i < actualVisibleEnd; i++ {
			item := cl.Items[i]
			itemID := fmt.Sprintf("%s-item-%d", cl.ID, i)

			// Determinar cor de fundo do item (destaque se em foco)
			itemBgColor := clay.Color{R: 0, G: 0, B: 0, A: 0} // Transparente por padrão
			if cl.HasFocus && cl.FocusedIndex == i {
				// Item em foco - gradiente azul moderno
				itemBgColor = clay.Color{R: 60, G: 120, B: 200, A: 255} // Alpha máximo
				log.Printf("Item %d (focused): BgColor R=%.0f G=%.0f B=%.0f A=%.0f", i, itemBgColor.R, itemBgColor.G, itemBgColor.B, itemBgColor.A)
			} else if item.Selected {
				// Item selecionado - fundo verde suave
				itemBgColor = clay.Color{R: 40, G: 120, B: 80, A: 255} // Alpha máximo
				log.Printf("Item %d (selected): BgColor R=%.0f G=%.0f B=%.0f A=%.0f", i, itemBgColor.R, itemBgColor.G, itemBgColor.B, itemBgColor.A)
			} else {
				log.Printf("Item %d (normal): BgColor R=%.0f G=%.0f B=%.0f A=%.0f", i, itemBgColor.R, itemBgColor.G, itemBgColor.B, itemBgColor.A)
			}

			// Container do item com dimensões fixas adequadas
			clay.UI()(clay.ElementDeclaration{
				Id: clay.ID(itemID),
				Layout: clay.LayoutConfig{
					Sizing: clay.Sizing{
						Width: clay.SizingPercent(1.0),
					},
					// Padding:         clay.Padding{Left: 15, Right: 15, Top: 12, Bottom: 12},
					Padding:         clay.PaddingAll(12),
					ChildGap:        12,
					LayoutDirection: clay.LEFT_TO_RIGHT,
				},
				CornerRadius:    clay.CornerRadiusAll(10), // Bordas arredondadas nos itens
				BackgroundColor: itemBgColor,
			}, func() {
				// Checkbox
				checkboxID := fmt.Sprintf("%s-checkbox-%d", cl.ID, i)
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
					// Marca de seleção (checkmark) se selecionado
					if item.Selected {
						// Container adicional para melhor controle de posicionamento
						claySystem.CreateText("◾", TextConfig{
							FontSize:  25,
							TextColor: clay.Color{R: 255, G: 255, B: 255, A: 255}, // Branco puro
						})
					}
				})

				// Label do item - container para controlar largura e quebra de texto
				labelColor := clay.Color{R: 220, G: 230, B: 245, A: 255} // Branco azulado suave por padrão
				if cl.HasFocus && cl.FocusedIndex == i {
					// Texto mais brilhante quando em foco
					labelColor = clay.Color{R: 255, G: 255, B: 255, A: 255} // Branco puro
				} else if item.Selected {
					// Texto verde claro para itens selecionados
					labelColor = clay.Color{R: 180, G: 255, B: 200, A: 255}
				}

				// Container para o texto com largura fixa no layout horizontal
				labelContainerID := fmt.Sprintf("%s-label-%d", cl.ID, i)
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
						FontSize:  15, // Fonte ligeiramente maior
						TextColor: labelColor,
					})
				})
			})
		}
	})

	// Atualizar posições visíveis
	maxItems := int(cl.Config.MaxHeight / cl.Config.ItemHeight)
	cl.VisibleStart = cl.ScrollOffset
	cl.VisibleEnd = min(cl.ScrollOffset+maxItems, len(cl.Items))

	log.Printf("Checkbox list created successfully: %s", cl.ID)
}

// DefaultCheckboxListConfig retorna configuração otimizada para layout automático seguindo padrões Clay
func DefaultCheckboxListConfig() CheckboxListConfig {
	return CheckboxListConfig{
		Sizing: clay.Sizing{
			Width: clay.SizingPercent(1.0),
			// Height: clay.SizingPercent(1.0),
		},
		Padding:         clay.PaddingAll(10),
		ChildGap:        5, // Espaçamento maior entre itens para acomodar texto multi-linha
		BackgroundColor: clay.Color{R: 25, G: 30, B: 40, A: 240},
		MaxHeight:       600, // Altura máxima aumentada
		ScrollOffset:    0,
		CheckboxSize:    25,
		ItemHeight:      45, // Altura maior para acomodar texto longo
	}
}

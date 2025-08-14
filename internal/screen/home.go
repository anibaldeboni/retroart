package screen

import (
	"log"
	"os"

	"github.com/TotallyGamerJet/clay"

	"retroart-sdl2/internal/core"
	"retroart-sdl2/internal/input"
	"retroart-sdl2/internal/theme"
	"retroart-sdl2/internal/ui"
	"retroart-sdl2/internal/ui/widgets"
)

type Home struct {
	navigator    Navigator // Use Navigator interface instead of concrete Manager
	buttons      []*widgets.Button
	checkboxList *widgets.CheckboxList[string]
}

func NewHome() *Home {
	home := &Home{}

	home.initializeWidgets()
	home.InitializeFocus()

	return home
}

// initializeWidgets cria todos os widgets focáveis
func (h *Home) initializeWidgets() {
	h.buttons = []*widgets.Button{
		widgets.NewButton("next-button", "Second screen", theme.StylePrimary, func() {
			if h.navigator != nil {
				h.navigator.NavigateTo("second")
			}
		}),
		widgets.NewButton("exit-button", "Exit", theme.StyleDanger, func() {
			log.Println("Exit button pressed")
			os.Exit(0)
		}),
		widgets.NewButton("test-selected-button", "Show Selected", theme.StyleSecondary, func() {
			selectedItems := h.checkboxList.GetSelectedItems()
			log.Printf("Selected games: %v", selectedItems)
		}),
	}

	// Criar dados de teste para o checkbox list
	testItems := []widgets.CheckboxListItem[string]{
		{Label: "Arcade", Value: "game1", Selected: false},
		{Label: "Gameboy", Value: "game2", Selected: true},
		{Label: "Gameboy color", Value: "game3", Selected: false},
		{Label: "Gameboy Advance", Value: "game4", Selected: false},
		{Label: "Nintendo Entertainment System", Value: "game5", Selected: true},
		{Label: "Super Nintendo", Value: "game6", Selected: false},
		{Label: "Master System", Value: "game7", Selected: false},
		{Label: "Mega Drive", Value: "game8", Selected: false},
		{Label: "Nintendo 64", Value: "game9", Selected: false},
		{Label: "Sega Saturn", Value: "game10", Selected: true},
		{Label: "Atari 2600", Value: "game11", Selected: false},
		{Label: "Game & Watch", Value: "game12", Selected: false},
		{Label: "CPS II", Value: "game13", Selected: false},
		{Label: "NeoGeo", Value: "game14", Selected: true},
		{Label: "GameGear", Value: "game15", Selected: false},
		{Label: "PlayStation", Value: "game16", Selected: false},
	}

	h.checkboxList = widgets.NewCheckboxList("consoles-checkbox-list", testItems)
}

// InitializeFocus configura os widgets no sistema de navegação espacial
func (h *Home) InitializeFocus() {
	layout := ui.GetLayout()
	if layout != nil {
		layout.RegisterFocusable(h.checkboxList)

		for _, button := range h.buttons {
			layout.RegisterFocusable(button)
		}
	}

	log.Println("Spatial navigation system initialized for Home")
}

// Implementação da interface Screen

func (h *Home) Update() {
	// Lógica de atualização se necessária
}

// Render - interface Screen (wrapper para o método Clay)
func (h *Home) Render() {
	mainContainerStyle := theme.GetMainContainerStyle()
	containerStyle := theme.GetContentContainerStyle()
	spacing := theme.GetSpacing()

	clay.UI()(clay.ElementDeclaration{
		Id: clay.ID("main-container"),
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  clay.SizingGrow(core.WINDOW_WIDTH),
				Height: clay.SizingGrow(core.WINDOW_HEIGHT),
			},
			Padding:         clay.Padding{Left: spacing.LG, Right: spacing.LG, Top: spacing.LG, Bottom: spacing.LG},
			ChildGap:        spacing.MD,
			LayoutDirection: clay.LEFT_TO_RIGHT,
		},
		BackgroundColor: mainContainerStyle.BackgroundColor,
	}, func() {
		// Container para lista de checkboxes (lado esquerdo)
		clay.UI()(clay.ElementDeclaration{
			Id: clay.ID("left-container"),
			Layout: clay.LayoutConfig{
				Sizing: clay.Sizing{
					Width:  clay.SizingPercent(0.35),
					Height: clay.SizingPercent(1.0),
				},
				Padding:         containerStyle.Padding,
				ChildGap:        spacing.SM,
				LayoutDirection: clay.TOP_TO_BOTTOM,
			},
			CornerRadius:    clay.CornerRadiusAll(containerStyle.CornerRadius),
			BackgroundColor: containerStyle.BackgroundColor,
		}, func() {
			// Título
			clay.UI()(clay.ElementDeclaration{
				Id: clay.ID("list-title"),
				Layout: clay.LayoutConfig{
					ChildAlignment: clay.ChildAlignment{
						X: clay.ALIGN_X_CENTER,
					},
				},
			}, func() {
				ds := theme.DefaultDesignSystem()
				widgets.TextLarge("Games list", ds.Colors.TextPrimary)
			})

			h.checkboxList.Render(core.WINDOW_HEIGHT - 220)
		})

		// Container para botões (lado direito)
		clay.UI()(clay.ElementDeclaration{
			Id: clay.ID("right-container"),
			Layout: clay.LayoutConfig{
				Sizing: clay.Sizing{
					Width:  clay.SizingPercent(0.65),
					Height: clay.SizingPercent(1.0),
				},
				Padding:         containerStyle.Padding,
				ChildGap:        spacing.MD,
				LayoutDirection: clay.TOP_TO_BOTTOM,
				ChildAlignment: clay.ChildAlignment{
					X: clay.ALIGN_X_CENTER,
				},
			},
			CornerRadius:    clay.CornerRadiusAll(containerStyle.CornerRadius),
			BackgroundColor: containerStyle.BackgroundColor,
		}, func() {
			// Título
			clay.UI()(clay.ElementDeclaration{
				Id: clay.ID("controls-title"),
				Layout: clay.LayoutConfig{
					ChildAlignment: clay.ChildAlignment{
						X: clay.ALIGN_X_CENTER,
					},
				},
			}, func() {
				ds := theme.DefaultDesignSystem()
				widgets.TextXLarge("Controls", ds.Colors.TextPrimary)
			})

			// Container para botões
			clay.UI()(clay.ElementDeclaration{
				Id: clay.ID("buttons-container"),
				Layout: clay.LayoutConfig{
					Padding:         clay.Padding{Left: spacing.SM, Right: spacing.SM, Top: spacing.SM, Bottom: spacing.SM},
					ChildGap:        spacing.MD,
					LayoutDirection: clay.TOP_TO_BOTTOM,
					ChildAlignment: clay.ChildAlignment{
						X: clay.ALIGN_X_CENTER,
					},
				},
			}, func() {
				// Renderizar todos os botões focáveis
				for _, button := range h.buttons {
					button.Render()
				}
			})

			// Informações do sistema de foco (debug)
			clay.UI()(clay.ElementDeclaration{
				Id: clay.ID("focus-debug"),
				Layout: clay.LayoutConfig{
					Padding: clay.PaddingAll(5),
					ChildAlignment: clay.ChildAlignment{
						X: clay.ALIGN_X_CENTER,
					},
				},
			}, func() {
				layout := ui.GetLayout()
				var currentFocus string
				var currentWidget ui.Focusable

				if layout != nil && layout.GetSpatialNavigation() != nil {
					currentFocus = layout.GetSpatialNavigation().GetCurrentFocus()
					currentWidget = layout.GetSpatialNavigation().GetCurrentWidget()
				}

				focusInfo := "Spatial Navigation"
				if currentFocus != "" {
					focusInfo += " | Focus: " + currentFocus
				}
				if currentWidget != nil {
					focusInfo += " | Widget: " + currentWidget.GetID()
				}

				ds := theme.DefaultDesignSystem()
				widgets.TextSmall(focusInfo, ds.Colors.TextSecondary)
			})
		})
	})
}

// HandleInput - refatorado para usar input.InputType diretamente
func (h *Home) HandleInput(inputType input.InputType) {
	switch inputType {
	case input.InputBack:
		if h.navigator != nil {
			h.navigator.GoBack() // Use Navigator's GoBack functionality
		}
		return
	default:
		layout := ui.GetLayout()
		if layout != nil {
			processed := layout.HandleSpatialInput(inputType)
			if processed {
				log.Printf("Input processed by focus system: %v", inputType)
			}
		}
	}
}

// OnEnter - chamado quando a tela se torna ativa
func (h *Home) OnEnter(navigator Navigator) {
	h.navigator = navigator // Store navigator reference
	log.Println("HomeV2 screen entered")
}

// OnExit - chamado quando a tela sai de foco
func (h *Home) OnExit() {
	log.Println("HomeV2 screen exited")
}

package screen

import (
	"log"
	"os"

	"github.com/TotallyGamerJet/clay"
	"github.com/veandco/go-sdl2/sdl"

	"retroart-sdl2/internal/core"
	"retroart-sdl2/internal/ui"
)

// Home é a versão refatorada da Home que usa o sistema de foco unificado
type Home struct {
	*BaseScreen
	screenMgr *Manager

	// Grupos focáveis
	buttonGroup *ui.FocusGroup
	listGroup   *ui.FocusGroup

	// Widgets focáveis
	buttons      []*ui.Button
	checkboxList *ui.FocusableCheckboxList[string]
}

func NewHome(screenMgr *Manager) *Home {
	home := &Home{
		BaseScreen: NewBaseScreen("home"),
		screenMgr:  screenMgr,
	}

	// Inicializar widgets focáveis
	home.initializeWidgets()

	// Configurar sistema de foco
	home.InitializeFocus()

	return home
}

// initializeWidgets cria todos os widgets focáveis
func (h *Home) initializeWidgets() {
	// Criar botões focáveis
	h.buttons = []*ui.Button{
		ui.NewButton("next-button", "Próxima Tela", ui.PrimaryButtonConfig(), func() {
			h.screenMgr.SetCurrentScreen("second")
		}),
		ui.NewButton("exit-button", "Sair", ui.DangerButtonConfig(), func() {
			log.Println("Exit button pressed")
			os.Exit(0)
		}),
		ui.NewButton("test-selected-button", "Mostrar Selecionados", ui.SecondaryButtonConfig(), func() {
			selectedItems := h.checkboxList.GetSelectedItems()
			log.Printf("Selected games: %v", selectedItems)
		}),
	}

	// Criar dados de teste para o checkbox list
	testItems := []ui.CheckboxListItem[string]{
		{Label: "Jogo de Ação Super Aventura", Value: "game1", Selected: false},
		{Label: "RPG Épico", Value: "game2", Selected: true},
		{Label: "Plataforma Retrô", Value: "game3", Selected: false},
		{Label: "Corrida de Velocidade", Value: "game4", Selected: false},
		{Label: "Puzzle Inteligente", Value: "game5", Selected: true},
		{Label: "Tiro em Primeira Pessoa", Value: "game6", Selected: false},
		{Label: "Estratégia em Tempo Real", Value: "game7", Selected: false},
		{Label: "Simulador de Vida", Value: "game8", Selected: false},
		{Label: "Aventura Point-and-Click", Value: "game9", Selected: false},
		{Label: "Luta Arcade Clássica", Value: "game10", Selected: true},
		{Label: "Música e Ritmo", Value: "game11", Selected: false},
		{Label: "Terror Psicológico", Value: "game12", Selected: false},
		{Label: "Alex Kidd in the miracle world", Value: "game13", Selected: false},
		{Label: "Street Fighter", Value: "game14", Selected: true},
		{Label: "Need for speed", Value: "game15", Selected: false},
		{Label: "BLACK", Value: "game16", Selected: false},
	}

	h.checkboxList = ui.NewFocusableCheckboxList("consoles-checkbox-list", testItems, ui.DefaultCheckboxListConfig())
}

// InitializeFocus configura os grupos de foco
func (h *Home) InitializeFocus() {
	// Criar grupo do checkbox list (primeiro, será o grupo ativo inicial)
	h.listGroup = ui.NewFocusGroup("checkbox-list")
	h.listGroup.AddFocusable(h.checkboxList)

	// Criar grupo de botões
	h.buttonGroup = ui.NewFocusGroup("buttons")
	for _, button := range h.buttons {
		h.buttonGroup.AddFocusable(button)
	}

	// Adicionar grupos ao gerenciador de foco (ordem importa para foco inicial)
	h.AddFocusGroup(h.listGroup)
	h.AddFocusGroup(h.buttonGroup)

	log.Println("Focus system initialized for HomeV2")
}

// Implementação da interface Screen

func (h *Home) Update() {
	// Lógica de atualização se necessária
}

// Render - interface Screen (wrapper para o método Clay)
func (h *Home) Render() {
	// Layout principal horizontal
	clay.UI()(clay.ElementDeclaration{
		Id: clay.ID("main-container"),
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  clay.SizingGrow(core.WINDOW_WIDTH),
				Height: clay.SizingGrow(core.WINDOW_HEIGHT),
			},
			Padding:         clay.PaddingAll(20),
			ChildGap:        15,
			LayoutDirection: clay.LEFT_TO_RIGHT,
		},
		BackgroundColor: clay.Color{R: 40, G: 42, B: 54, A: 255},
	}, func() {
		// Container para lista de checkboxes (lado esquerdo)
		clay.UI()(clay.ElementDeclaration{
			Id: clay.ID("left-container"),
			Layout: clay.LayoutConfig{
				Sizing: clay.Sizing{
					Width:  clay.SizingPercent(0.35),
					Height: clay.SizingPercent(1.0),
				},
				Padding:         clay.PaddingAll(15),
				ChildGap:        10,
				LayoutDirection: clay.TOP_TO_BOTTOM,
			},
			CornerRadius:    clay.CornerRadiusAll(12),
			BackgroundColor: clay.Color{R: 60, G: 63, B: 75, A: 180},
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
				clay.Text("Lista de Jogos", &clay.TextElementConfig{
					FontSize:  20,
					TextColor: clay.Color{R: 255, G: 255, B: 255, A: 255},
				})
			})

			// Renderizar checkbox list focável
			h.checkboxList.Render(core.WINDOW_HEIGHT - 170)
		})

		// Container para botões (lado direito)
		clay.UI()(clay.ElementDeclaration{
			Id: clay.ID("right-container"),
			Layout: clay.LayoutConfig{
				Sizing: clay.Sizing{
					Width:  clay.SizingPercent(0.65),
					Height: clay.SizingPercent(1.0),
				},
				Padding:         clay.PaddingAll(15),
				ChildGap:        15,
				LayoutDirection: clay.TOP_TO_BOTTOM,
				ChildAlignment: clay.ChildAlignment{
					X: clay.ALIGN_X_CENTER,
				},
			},
			CornerRadius:    clay.CornerRadiusAll(12),
			BackgroundColor: clay.Color{R: 60, G: 63, B: 75, A: 180},
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
				clay.Text("Controls", &clay.TextElementConfig{
					FontSize:  24,
					TextColor: clay.Color{R: 255, G: 255, B: 255, A: 255},
				})
			})

			// Container para botões
			clay.UI()(clay.ElementDeclaration{
				Id: clay.ID("buttons-container"),
				Layout: clay.LayoutConfig{
					Padding:         clay.PaddingAll(10),
					ChildGap:        15,
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
				currentGroup := h.GetCurrentGroup()
				if currentGroup != nil {
					currentFocusable := h.GetCurrentFocusable()
					focusInfo := "Grupo: " + currentGroup.ID
					if currentFocusable != nil {
						focusInfo += " | Widget: " + currentFocusable.GetID()
					}

					clay.Text(focusInfo, &clay.TextElementConfig{
						FontSize:  12,
						TextColor: clay.Color{R: 200, G: 200, B: 200, A: 255},
					})
				}
			})
		})
	})
}

// HandleInput - compatibilidade com screen.InputType
func (h *Home) HandleInput(input InputType) {
	var direction ui.InputDirection
	switch input {
	case InputUp:
		direction = ui.DirectionUp
	case InputDown:
		direction = ui.DirectionDown
	case InputLeft:
		direction = ui.DirectionLeft
	case InputRight:
		direction = ui.DirectionRight
	case InputConfirm:
		direction = ui.DirectionConfirm
	case InputBack:
		h.screenMgr.SetCurrentScreen("home") // voltar ou ação específica
		return
	default:
		return
	}

	processed := h.BaseScreen.HandleInput(direction)
	if processed {
		log.Printf("Input processed by focus system: %v", input)
	}
}

// OnEnter - chamado quando a tela se torna ativa
func (h *Home) OnEnter() {
	log.Println("HomeV2 screen entered")
}

// OnExit - chamado quando a tela sai de foco
func (h *Home) OnExit() {
	log.Println("HomeV2 screen exited")
}

// HandleKeyDown processa teclas pressionadas (método para uso direto)
func (h *Home) HandleKeyDown(key sdl.Scancode) {
	var direction ui.InputDirection
	handled := false

	switch key {
	case sdl.SCANCODE_UP:
		direction = ui.DirectionUp
		handled = true
	case sdl.SCANCODE_DOWN:
		direction = ui.DirectionDown
		handled = true
	case sdl.SCANCODE_LEFT:
		direction = ui.DirectionLeft
		handled = true
	case sdl.SCANCODE_RIGHT:
		direction = ui.DirectionRight
		handled = true
	case sdl.SCANCODE_RETURN, sdl.SCANCODE_SPACE:
		direction = ui.DirectionConfirm
		handled = true
	}

	if handled {
		processed := h.BaseScreen.HandleInput(direction)
		if processed {
			log.Printf("Input processed: %v", direction)
		} else {
			log.Printf("Input not processed: %v", direction)
		}
	}
}

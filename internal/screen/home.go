package screen

import (
	"log"

	"github.com/TotallyGamerJet/clay"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"

	"retroart-sdl2/internal/core"
	"retroart-sdl2/internal/ui"
)

type Home struct {
	screenMgr     *Manager
	claySystem    *ui.ClayLayoutSystem
	selectedIndex int
	buttonCount   int
	checkboxList  *ui.CheckboxList[string] // Usando string como tipo genérico
	focusMode     FocusMode                // Controla se estamos focando botões ou checkbox list
}

// Modos de foco
type FocusMode int

const (
	FocusButtons      FocusMode = iota // Foco nos botões (padrão)
	FocusCheckboxList                  // Foco no checkbox list
)

func NewHome(screenMgr *Manager, renderer *sdl.Renderer, font *ttf.Font) *Home {
	screen := &Home{
		screenMgr:     screenMgr,
		selectedIndex: 0,
		buttonCount:   3,            // próxima tela, sair, teste selecionados
		focusMode:     FocusButtons, // Inicia focando nos botões
	}

	// Criar sistema Clay para esta tela
	screen.claySystem = ui.NewClayLayoutSystem(renderer, font)

	// Criar dados de teste para o checkbox list - mais itens para testar sizing dinâmico
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

	screen.checkboxList = ui.NewCheckboxList("consoles-checkbox-list", testItems, ui.DefaultCheckboxListConfig())

	return screen
}

func (hs *Home) Update() {
	// Lógica de atualização se necessária
}

func (hs *Home) Render(renderer *sdl.Renderer) {
	// Iniciar layout Clay
	hs.claySystem.BeginLayout()

	// Layout principal horizontal com dimensões fixas no raiz para estabelecer contexto
	clay.UI()(clay.ElementDeclaration{
		Id: clay.ID("main-container"),
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  clay.SizingGrow(core.WINDOW_WIDTH),  // Largura da janela menos padding (1280 - 40)
				Height: clay.SizingGrow(core.WINDOW_HEIGHT), // Altura da janela menos padding (720 - 40)
			},
			Padding:         clay.PaddingAll(20),
			ChildGap:        20,
			LayoutDirection: clay.LEFT_TO_RIGHT, // Layout horizontal
		},
		BackgroundColor: clay.Color{R: 32, G: 34, B: 37, A: 255}, // Fundo escuro principal
	}, func() {
		clay.UI()(clay.ElementDeclaration{
			Id: clay.ID("checkbox-container"),
			Layout: clay.LayoutConfig{
				Sizing: clay.Sizing{
					Width:  clay.SizingPercent(0.4),
					Height: clay.SizingPercent(1.0),
				},
				Padding:         clay.PaddingAll(10),
				LayoutDirection: clay.TOP_TO_BOTTOM,
			},
			CornerRadius:    clay.CornerRadiusAll(12),
			BackgroundColor: clay.Color{R: 45, G: 50, B: 65, A: 255},
		}, func() {
			containerHeight := float32(600)
			hs.checkboxList.RenderCheckboxList(hs.claySystem, containerHeight)
		})

		clay.UI()(clay.ElementDeclaration{
			Id: clay.ID("content-container"),
			Layout: clay.LayoutConfig{
				Sizing: clay.Sizing{
					Width:  clay.SizingPercent(0.6), // 60% de 1200 menos gap (1200*0.6 - 20)
					Height: clay.SizingPercent(1.0), // Altura disponível (680-40 de padding)
				},
				Padding:         clay.PaddingAll(20),
				ChildGap:        15,
				LayoutDirection: clay.TOP_TO_BOTTOM,
				ChildAlignment: clay.ChildAlignment{
					X: clay.ALIGN_X_CENTER, // Centralizar botões horizontalmente
				},
			},
			CornerRadius:    clay.CornerRadiusAll(12),
			BackgroundColor: clay.Color{R: 45, G: 50, B: 65, A: 255}, // Fundo do painel direito
		}, func() {

			// Título
			clay.UI()(clay.ElementDeclaration{
				Id: clay.ID("title-text"),
				Layout: clay.LayoutConfig{
					Sizing: clay.Sizing{
						Width:  clay.SizingPercent(1.0),
						Height: clay.SizingFit(0, 50),
					},
					ChildAlignment: clay.ChildAlignment{
						X: clay.ALIGN_X_CENTER,
						Y: clay.ALIGN_Y_CENTER,
					},
				},
			}, func() {
				hs.claySystem.CreateText("RetroArt", ui.TextConfig{
					FontSize:  48,
					TextColor: clay.Color{R: 100, G: 200, B: 255, A: 255}, // Azul moderno vibrante
				})
			})

			// Subtítulo
			clay.UI()(clay.ElementDeclaration{
				Id: clay.ID("subtitle-text"),
				Layout: clay.LayoutConfig{
					// Sizing: clay.Sizing{
					// 	Width:  clay.SizingGrow(1),
					// 	Height: clay.SizingFit(0, 0),
					// },
					ChildAlignment: clay.ChildAlignment{
						X: clay.ALIGN_X_CENTER,
					},
				},
			}, func() {
				hs.claySystem.CreateText("Aplicação Gráfica para Trimui Smart Pro", ui.TextConfig{
					FontSize:  16,
					TextColor: clay.Color{R: 180, G: 190, B: 210, A: 255}, // Cinza azulado suave
				})
			})

			// Container para botões com layout centralizado
			clay.UI()(clay.ElementDeclaration{
				Id: clay.ID("buttons-container"),
				Layout: clay.LayoutConfig{
					Sizing: clay.Sizing{
						Width:  clay.SizingGrow(1),
						Height: clay.SizingFit(0, 0),
					},
					Padding:         clay.PaddingAll(15),
					ChildGap:        15,
					LayoutDirection: clay.TOP_TO_BOTTOM,
					ChildAlignment: clay.ChildAlignment{
						X: clay.ALIGN_X_CENTER, // Centralizar botões
					},
				},
			}, func() {

				// Botão próxima tela - Estilo moderno azul
				nextButtonConfig := ui.DefaultButtonConfig()
				if hs.selectedIndex == 0 && hs.focusMode == FocusButtons {
					nextButtonConfig.BackgroundColor = clay.Color{R: 30, G: 150, B: 255, A: 255} // Azul vibrante focado
					nextButtonConfig.TextColor = clay.Color{R: 255, G: 255, B: 255, A: 255}      // Texto branco
				} else {
					nextButtonConfig.BackgroundColor = clay.Color{R: 50, G: 120, B: 200, A: 200} // Azul suave normal
					nextButtonConfig.TextColor = clay.Color{R: 220, G: 230, B: 255, A: 255}      // Texto azul claro
				}
				hs.claySystem.CreateButton("next-button", "Próxima Tela", nextButtonConfig, func() {
					hs.screenMgr.SetCurrentScreen("second")
				})

				// Botão sair - Estilo moderno vermelho
				exitButtonConfig := ui.DefaultButtonConfig()
				if hs.selectedIndex == 1 && hs.focusMode == FocusButtons {
					exitButtonConfig.BackgroundColor = clay.Color{R: 255, G: 80, B: 80, A: 255} // Vermelho vibrante focado
					exitButtonConfig.TextColor = clay.Color{R: 255, G: 255, B: 255, A: 255}     // Texto branco
				} else {
					exitButtonConfig.BackgroundColor = clay.Color{R: 200, G: 60, B: 60, A: 180} // Vermelho suave normal
					exitButtonConfig.TextColor = clay.Color{R: 255, G: 200, B: 200, A: 255}     // Texto vermelho claro
				}
				hs.claySystem.CreateButton("exit-button", "Sair", exitButtonConfig, func() {
					log.Println("Exit button pressed")
				})

				// Botão de teste - Estilo moderno roxo/violeta
				testButtonConfig := ui.DefaultButtonConfig()
				if hs.selectedIndex == 2 && hs.focusMode == FocusButtons {
					testButtonConfig.BackgroundColor = clay.Color{R: 150, G: 80, B: 255, A: 255} // Roxo vibrante focado
					testButtonConfig.TextColor = clay.Color{R: 255, G: 255, B: 255, A: 255}      // Texto branco
				} else {
					testButtonConfig.BackgroundColor = clay.Color{R: 120, G: 60, B: 200, A: 180} // Roxo suave normal
					testButtonConfig.TextColor = clay.Color{R: 220, G: 200, B: 255, A: 255}      // Texto roxo claro
				}
				hs.claySystem.CreateButton("test-selected-button", "Mostrar Selecionados", testButtonConfig, func() {
					selected := hs.checkboxList.GetSelectedItems()
					log.Printf("=== ELEMENTOS SELECIONADOS ===")
					for i, item := range selected {
						log.Printf("Item %d: Label='%s', Value='%v'", i+1, item.Label, item.Value)
					}
					log.Printf("=== FIM DA LISTA ===")

					// Também testar GetSelectedValues
					values := hs.checkboxList.GetSelectedValues()
					log.Printf("Valores selecionados: %v", values)
				})
			})
		})
	})

	// Finalizar e renderizar
	hs.claySystem.Render()
}

func (hs *Home) HandleInput(input InputType) {
	switch hs.focusMode {
	case FocusButtons:
		hs.handleButtonInput(input)
	case FocusCheckboxList:
		hs.handleCheckboxListInput(input)
	}
}

// handleButtonInput lida com entrada quando focando nos botões
func (hs *Home) handleButtonInput(input InputType) {
	switch input {
	case InputUp:
		if hs.selectedIndex > 0 {
			hs.selectedIndex--
		}
	case InputDown:
		if hs.selectedIndex < hs.buttonCount-1 {
			hs.selectedIndex++
		}
	case InputLeft:
		// Mover foco para o checkbox list
		hs.focusMode = FocusCheckboxList
		hs.checkboxList.SetFocused(true)
		log.Println("Foco movido para checkbox list")
	case InputConfirm:
		hs.executeCurrentButton()
	}
}

// handleCheckboxListInput lida com entrada quando focando no checkbox list
func (hs *Home) handleCheckboxListInput(input InputType) {
	switch input {
	case InputUp:
		hs.checkboxList.MoveFocusUp()
	case InputDown:
		hs.checkboxList.MoveFocusDown()
	case InputConfirm:
		hs.checkboxList.ToggleFocusedItem()
		log.Printf("Item toggled. Selected items: %d", len(hs.checkboxList.GetSelectedItems()))
	case InputRight:
		// Voltar foco para os botões
		hs.focusMode = FocusButtons
		hs.checkboxList.SetFocused(false)
		log.Println("Foco movido para botões")
	}
}

func (hs *Home) executeCurrentButton() {
	switch hs.selectedIndex {
	case 0: // Próxima tela
		hs.screenMgr.SetCurrentScreen("second")
	case 1: // Sair
		log.Println("Exit button pressed from Home screen")
	case 2: // Mostrar selecionados
		selected := hs.checkboxList.GetSelectedItems()
		log.Printf("=== ELEMENTOS SELECIONADOS ===")
		for i, item := range selected {
			log.Printf("Item %d: Label='%s', Value='%v'", i+1, item.Label, item.Value)
		}
		log.Printf("=== FIM DA LISTA ===")

		// Também testar GetSelectedValues
		values := hs.checkboxList.GetSelectedValues()
		log.Printf("Valores selecionados: %v", values)
	}
}

func (hs *Home) OnEnter() {
	log.Println("Entering Home screen")
	hs.focusMode = FocusButtons
	hs.checkboxList.SetFocused(false)
}

func (hs *Home) OnExit() {
	log.Println("Exiting Home screen")
}

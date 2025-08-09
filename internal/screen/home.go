package screen

import (
	"log"

	"github.com/TotallyGamerJet/clay"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"

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

	// Criar dados de teste para o checkbox list
	testItems := []ui.CheckboxListItem[string]{
		{Label: "Item 1", Value: "value1", Selected: false},
		{Label: "Item 2", Value: "value2", Selected: true},
		{Label: "Item 3", Value: "value3", Selected: false},
		{Label: "Item 4", Value: "value4", Selected: false},
		{Label: "Item 5", Value: "value5", Selected: true},
		{Label: "Item 6", Value: "value6", Selected: false},
		{Label: "Item 7", Value: "value7", Selected: false},
		{Label: "Item 8", Value: "value8", Selected: false},
	}

	screen.checkboxList = ui.NewCheckboxList("test-checkbox-list", testItems, ui.CheckboxListConfig{
		Sizing: clay.Sizing{
			Width:  clay.SizingFixed(320),
			Height: clay.SizingFixed(350),
		},
		Padding:         clay.PaddingAll(20),
		ChildGap:        8,
		BackgroundColor: clay.Color{R: 30, G: 35, B: 45, A: 240}, // Fundo escuro translúcido moderno
		ItemHeight:      35,
		MaxHeight:       350,
		CheckboxSize:    18,
	})

	return screen
}

func (hs *Home) Update() {
	// Lógica de atualização se necessária
}

func (hs *Home) Render(renderer *sdl.Renderer) {
	// Iniciar layout Clay
	hs.claySystem.BeginLayout()

	// Container principal
	hs.claySystem.CreateContainer("main", ui.ContainerConfig{
		Sizing: clay.Sizing{
			Width:  clay.SizingGrow(0),
			Height: clay.SizingGrow(0),
		},
		Padding:         clay.PaddingAll(40),
		ChildGap:        40,
		LayoutDirection: clay.LEFT_TO_RIGHT,                      // Layout horizontal
		BackgroundColor: clay.Color{R: 20, G: 25, B: 35, A: 255}, // Fundo principal mais escuro e elegante
	}, func() {
		// Checkbox list à esquerda
		hs.claySystem.CreateCheckboxList(hs.checkboxList)

		// Container para conteúdo principal
		hs.claySystem.CreateContainer("content", ui.ContainerConfig{
			Sizing: clay.Sizing{
				Width:  clay.SizingGrow(0),
				Height: clay.SizingGrow(0),
			},
			Padding:         clay.PaddingAll(30),
			ChildGap:        25,
			LayoutDirection: clay.TOP_TO_BOTTOM,
			BackgroundColor: clay.Color{R: 35, G: 40, B: 50, A: 200}, // Fundo com transparência sutil
		}, func() {
			// Título
			hs.claySystem.CreateText("RetroArt", ui.TextConfig{
				FontSize:  32,
				TextColor: clay.Color{R: 100, G: 200, B: 255, A: 255}, // Azul moderno vibrante
			})

			// Subtítulo
			hs.claySystem.CreateText("Aplicação Gráfica para Trimui Smart Pro", ui.TextConfig{
				FontSize:  16,
				TextColor: clay.Color{R: 180, G: 190, B: 210, A: 255}, // Cinza azulado suave
			})

			// Container para botões
			hs.claySystem.CreateContainer("buttons", ui.ContainerConfig{
				Sizing: clay.Sizing{
					Width:  clay.SizingFixed(320),
					Height: clay.SizingFit(0, 1000),
				},
				Padding:         clay.PaddingAll(15),
				ChildGap:        15,
				LayoutDirection: clay.TOP_TO_BOTTOM,
				BackgroundColor: clay.Color{R: 0, G: 0, B: 0, A: 0}, // Transparente
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
					// Note: Seria melhor ter uma referência à app ou um callback
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

package screen

import (
	"log"
	"os"

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
	buttons       []ButtonDefinition       // Lista de botões com seus comportamentos
	checkboxList  *ui.CheckboxList[string] // Usando string como tipo genérico
	focusMode     FocusMode                // Controla se estamos focando botões ou checkbox list
}

// Modos de foco
type FocusMode int

const (
	FocusButtons      FocusMode = iota // Foco nos botões (padrão)
	FocusCheckboxList                  // Foco no checkbox list
)

// ButtonDefinition define um botão com seu comportamento
type ButtonDefinition struct {
	ID      string
	Label   string
	Config  ui.StatefulButtonConfig
	OnClick func()
}

func NewHome(screenMgr *Manager, renderer *sdl.Renderer, font *ttf.Font) *Home {
	screen := &Home{
		screenMgr:     screenMgr,
		selectedIndex: 0,
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

	// Inicializar botões com seus comportamentos
	screen.initializeButtons()

	return screen
}

// initializeButtons define os botões e seus comportamentos
func (hs *Home) initializeButtons() {
	hs.buttons = []ButtonDefinition{
		{
			ID:     "next-button",
			Label:  "Próxima Tela",
			Config: ui.PrimaryButtonConfig(),
			OnClick: func() {
				hs.screenMgr.SetCurrentScreen("second")
			},
		},
		{
			ID:     "exit-button",
			Label:  "Sair",
			Config: ui.DangerButtonConfig(),
			OnClick: func() {
				log.Println("Exit button pressed")
				os.Exit(0)
			},
		},
		{
			ID:     "test-selected-button",
			Label:  "Mostrar Selecionados",
			Config: ui.SecondaryButtonConfig(),
			OnClick: func() {
				selected := hs.checkboxList.GetSelectedItems()
				log.Printf("=== ELEMENTOS SELECIONADOS ===")
				for i, item := range selected {
					log.Printf("Item %d: Label='%s', Value='%v'", i+1, item.Label, item.Value)
				}
				log.Printf("=== FIM DA LISTA ===")

				// Também testar GetSelectedValues
				values := hs.checkboxList.GetSelectedValues()
				log.Printf("Valores selecionados: %v", values)
			},
		},
	}
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
					Width:  clay.SizingPercent(0.6),
					Height: clay.SizingPercent(1.0),
				},
				Padding:         clay.PaddingAll(20),
				ChildGap:        15,
				LayoutDirection: clay.TOP_TO_BOTTOM,
				ChildAlignment: clay.ChildAlignment{
					X: clay.ALIGN_X_CENTER,
				},
			},
			CornerRadius:    clay.CornerRadiusAll(12),
			BackgroundColor: clay.Color{R: 45, G: 50, B: 65, A: 255},
		}, func() {

			// Título
			clay.UI()(clay.ElementDeclaration{
				Id: clay.ID("title-text"),
				Layout: clay.LayoutConfig{
					Sizing: clay.Sizing{
						Width:  clay.SizingPercent(1.0),
						Height: clay.SizingFit(0, 100),
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
					Padding:         clay.PaddingAll(15),
					ChildGap:        15,
					LayoutDirection: clay.TOP_TO_BOTTOM,
					ChildAlignment: clay.ChildAlignment{
						X: clay.ALIGN_X_CENTER, // Centralizar botões
					},
				},
			}, func() {

				// Renderizar botões dinamicamente usando CreateStatefulButton
				for i, button := range hs.buttons {
					isFocused := hs.selectedIndex == i && hs.focusMode == FocusButtons
					hs.claySystem.CreateStatefulButton(button.ID, button.Label, button.Config, isFocused, button.OnClick)
				}
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
		if hs.selectedIndex < len(hs.buttons)-1 {
			hs.selectedIndex++
		}
	case InputLeft:
		// Mover foco para o checkbox list
		hs.focusMode = FocusCheckboxList
		hs.checkboxList.SetFocused(true)
		log.Println("Foco movido para checkbox list")
	case InputConfirm:
		// Executar a função OnClick do botão selecionado
		if hs.selectedIndex >= 0 && hs.selectedIndex < len(hs.buttons) {
			hs.buttons[hs.selectedIndex].OnClick()
		}
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

func (hs *Home) OnEnter() {
	log.Println("Entering Home screen")
	hs.focusMode = FocusButtons
	hs.checkboxList.SetFocused(false)
}

func (hs *Home) OnExit() {
	log.Println("Exiting Home screen")
}

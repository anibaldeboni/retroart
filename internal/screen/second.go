package screen

import (
	"log"
	"retroart-sdl2/internal/core"
	"retroart-sdl2/internal/ui"

	"github.com/TotallyGamerJet/clay"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Second struct {
	*BaseScreen
	screenMgr  *Manager
	claySystem *ui.ClayLayoutSystem

	// Widgets focáveis
	buttons []*ui.FocusableButton

	// Grupos de foco
	buttonGroup *ui.FocusGroup
}

func NewSecond(screenMgr *Manager, renderer *sdl.Renderer, font *ttf.Font) *Second {
	screen := &Second{
		BaseScreen: NewBaseScreen("second-screen"),
		screenMgr:  screenMgr,
	}

	// Criar sistema Clay para esta tela
	screen.claySystem = ui.NewClayLayoutSystem(renderer, font)

	screen.initializeWidgets()
	screen.InitializeFocus()

	return screen
}

func (ss *Second) initializeWidgets() {
	// Criar botões focáveis
	ss.buttons = []*ui.FocusableButton{
		ui.NewFocusableButton("back-btn", "Voltar", ui.PrimaryButtonConfig(), func() {
			ss.screenMgr.SetCurrentScreen("home")
		}),
		ui.NewFocusableButton("options-btn", "Opções", ui.SecondaryButtonConfig(), func() {
			// Ação para opções (pode ser implementada futuramente)
		}),
		ui.NewFocusableButton("exit-btn", "Sair", ui.DangerButtonConfig(), func() {
			// Ação para sair
		}),
	}
}

func (ss *Second) InitializeFocus() {
	// Criar grupo de botões
	ss.buttonGroup = ui.NewFocusGroup("second-buttons")
	for _, btn := range ss.buttons {
		ss.buttonGroup.AddFocusable(btn)
	}

	// Adicionar grupos ao gerenciador de foco
	ss.AddFocusGroup(ss.buttonGroup)
}

func (ss *Second) Update() {
	// Lógica de atualização se necessária
}

func (ss *Second) Render(renderer *sdl.Renderer) {
	ss.RenderWithClay()
}

// RenderWithClay - método específico para renderização Clay
func (ss *Second) RenderWithClay() {
	ss.claySystem.BeginLayout()

	// Layout principal centralizado
	clay.UI()(clay.ElementDeclaration{
		Id: clay.ID("main-container"),
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  clay.SizingGrow(core.WINDOW_WIDTH),
				Height: clay.SizingGrow(core.WINDOW_HEIGHT),
			},
			Padding:         clay.PaddingAll(20),
			ChildGap:        30,
			LayoutDirection: clay.TOP_TO_BOTTOM,
			ChildAlignment: clay.ChildAlignment{
				X: clay.ALIGN_X_CENTER,
				Y: clay.ALIGN_Y_CENTER,
			},
		},
		BackgroundColor: clay.Color{R: 40, G: 42, B: 54, A: 255},
	}, func() {
		// Container para conteúdo principal
		clay.UI()(clay.ElementDeclaration{
			Id: clay.ID("content-container"),
			Layout: clay.LayoutConfig{
				Sizing: clay.Sizing{
					Width:  clay.SizingPercent(0.8),
					Height: clay.SizingFit(0, 0),
				},
				Padding:         clay.PaddingAll(30),
				ChildGap:        20,
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
				Id: clay.ID("title"),
				Layout: clay.LayoutConfig{
					ChildAlignment: clay.ChildAlignment{
						X: clay.ALIGN_X_CENTER,
					},
				},
			}, func() {
				ss.claySystem.CreateText("Segunda Tela", ui.TextConfig{
					FontSize:  28,
					TextColor: clay.Color{R: 255, G: 255, B: 255, A: 255},
				})
			})

			// Container para textos de conteúdo
			clay.UI()(clay.ElementDeclaration{
				Id: clay.ID("text-content"),
				Layout: clay.LayoutConfig{
					ChildGap:        15,
					LayoutDirection: clay.TOP_TO_BOTTOM,
					ChildAlignment: clay.ChildAlignment{
						X: clay.ALIGN_X_CENTER,
					},
				},
			}, func() {
				ss.claySystem.CreateText("Esta é a segunda tela da aplicação.", ui.TextConfig{
					FontSize:  18,
					TextColor: clay.Color{R: 230, G: 230, B: 230, A: 255},
				})

				ss.claySystem.CreateText("Aqui você pode adicionar qualquer conteúdo desejado.", ui.TextConfig{
					FontSize:  16,
					TextColor: clay.Color{R: 200, G: 200, B: 200, A: 255},
				})

				ss.claySystem.CreateText("Esta estrutura permite fácil expansão.", ui.TextConfig{
					FontSize:  16,
					TextColor: clay.Color{R: 200, G: 200, B: 200, A: 255},
				})

				ss.claySystem.CreateText("O Clay permite layouts flexíveis e responsivos.", ui.TextConfig{
					FontSize:  16,
					TextColor: clay.Color{R: 200, G: 200, B: 200, A: 255},
				})
			})

			// Container para botões
			clay.UI()(clay.ElementDeclaration{
				Id: clay.ID("buttons-container"),
				Layout: clay.LayoutConfig{
					Padding:         clay.PaddingAll(10),
					ChildGap:        15,
					LayoutDirection: clay.LEFT_TO_RIGHT,
					ChildAlignment: clay.ChildAlignment{
						X: clay.ALIGN_X_CENTER,
					},
				},
			}, func() {
				// Renderizar todos os botões focáveis
				for _, button := range ss.buttons {
					button.Render(ss.claySystem)
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
				currentGroup := ss.GetCurrentGroup()
				if currentGroup != nil {
					currentFocusable := ss.GetCurrentFocusable()
					focusInfo := "Grupo: " + currentGroup.ID
					if currentFocusable != nil {
						focusInfo += " | Widget: " + currentFocusable.GetID()
					}

					ss.claySystem.CreateText(focusInfo, ui.TextConfig{
						FontSize:  12,
						TextColor: clay.Color{R: 200, G: 200, B: 200, A: 255},
					})
				}
			})
		})
	})

	// Finalizar e renderizar
	ss.claySystem.Render()
}

func (ss *Second) HandleInput(input InputType) {
	// Mapear InputType para InputDirection
	var direction ui.InputDirection
	handled := false

	switch input {
	case InputUp:
		direction = ui.DirectionUp
		handled = true
	case InputDown:
		direction = ui.DirectionDown
		handled = true
	case InputLeft:
		direction = ui.DirectionLeft
		handled = true
	case InputRight:
		direction = ui.DirectionRight
		handled = true
	case InputConfirm:
		direction = ui.DirectionConfirm
		handled = true
	case InputBack:
		ss.screenMgr.SetCurrentScreen("home")
		return
	}

	if handled {
		ss.BaseScreen.HandleInput(direction)
	}
}

func (ss *Second) OnEnter() {
	log.Println("Entering Second screen")
}

func (ss *Second) OnExit() {
	log.Println("Exiting Second screen")
}

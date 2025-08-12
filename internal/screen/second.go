package screen

import (
	"log"
	"retroart-sdl2/internal/core"
	"retroart-sdl2/internal/input"
	"retroart-sdl2/internal/ui"

	"github.com/TotallyGamerJet/clay"
)

type Second struct {
	navigator Navigator // Use Navigator interface instead of concrete Manager
	buttons   []*ui.Button
}

func NewSecond() *Second {
	screen := &Second{
		// navigator will be set in OnEnter
	}

	screen.initializeWidgets()
	screen.InitializeFocus()

	return screen
}

func (ss *Second) initializeWidgets() {
	// Criar botões focáveis
	ss.buttons = []*ui.Button{
		ui.NewButton("back-btn", "Back", ui.PrimaryButtonConfig(), func() {
			if ss.navigator != nil {
				ss.navigator.GoBack()
			}
		}),
		ui.NewButton("options-btn", "Options", ui.SecondaryButtonConfig(), func() {
			// Ação para opções (pode ser implementada futuramente)
		}),
		ui.NewButton("exit-btn", "Exit", ui.DangerButtonConfig(), func() {
			// Ação para sair
		}),
	}
}

func (ss *Second) InitializeFocus() {
	layout := ui.GetLayout()
	if layout != nil {
		for _, btn := range ss.buttons {
			layout.RegisterFocusable(btn)
		}
	}
}

func (ss *Second) Update() {
	// Lógica de atualização se necessária
}

func (ss *Second) Render() {
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
				clay.Text("Second screen", &clay.TextElementConfig{
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
				clay.Text("This is the second application screen.", &clay.TextElementConfig{
					FontSize:  18,
					TextColor: clay.Color{R: 230, G: 230, B: 230, A: 255},
				})

				clay.Text("Here you can add any desired content.", &clay.TextElementConfig{
					FontSize:  16,
					TextColor: clay.Color{R: 200, G: 200, B: 200, A: 255},
				})

				clay.Text("This structure allows for easy expansion.", &clay.TextElementConfig{
					FontSize:  16,
					TextColor: clay.Color{R: 200, G: 200, B: 200, A: 255},
				})

				clay.Text("Clay allows for flexible and responsive layouts.", &clay.TextElementConfig{
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
				for _, button := range ss.buttons {
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

				clay.Text(focusInfo, &clay.TextElementConfig{
					FontSize:  12,
					TextColor: clay.Color{R: 200, G: 200, B: 200, A: 255},
				})
			})
		})
	})
}

func (ss *Second) HandleInput(inputType input.InputType) {
	var handled bool

	switch inputType {
	case input.InputBack:
		if ss.navigator != nil {
			ss.navigator.GoBack()
		}
		return
	default:
		layout := ui.GetLayout()
		if layout != nil {
			handled = layout.HandleSpatialInput(inputType)
		}
	}

	if !handled {
		log.Printf("Second: Input %d not handled", inputType)
	}
}

func (ss *Second) OnEnter(navigator Navigator) {
	ss.navigator = navigator // Store navigator reference
	log.Println("Entering Second screen")
}

func (ss *Second) OnExit() {
	log.Println("Exiting Second screen")
}

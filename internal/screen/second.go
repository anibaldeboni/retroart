package screen

import (
	"log"
	"retroart-sdl2/internal/core"
	"retroart-sdl2/internal/input"
	"retroart-sdl2/internal/theme"
	"retroart-sdl2/internal/ui"
	"retroart-sdl2/internal/ui/widgets"

	"github.com/TotallyGamerJet/clay"
)

type Second struct {
	navigator Navigator // Use Navigator interface instead of concrete Manager
	buttons   []*widgets.Button
}

func NewSecond() *Second {
	screen := &Second{}

	screen.initializeWidgets()
	screen.InitializeFocus()

	return screen
}

func (ss *Second) initializeWidgets() {
	ss.buttons = []*widgets.Button{
		widgets.NewButton("back-btn", "Back", theme.StylePrimary, func() {
			if ss.navigator != nil {
				ss.navigator.GoBack()
			}
		}),
		widgets.NewButton("options-btn", "Options", theme.StyleSecondary, func() {
			// Ação para opções (pode ser implementada futuramente)
		}),
		widgets.NewButton("exit-btn", "Exit", theme.StyleDanger, func() {
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
	mainStyle := theme.GetMainContainerStyle()
	contentStyle := theme.GetContentContainerStyle()
	colors := theme.GetColors()
	spacing := theme.GetSpacing()
	typography := theme.GetTypography()

	clay.UI()(clay.ElementDeclaration{
		Id: clay.ID("main-container"),
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  clay.SizingGrow(core.WINDOW_WIDTH),
				Height: clay.SizingGrow(core.WINDOW_HEIGHT),
			},
			Padding:         clay.Padding{Left: spacing.LG, Right: spacing.LG, Top: spacing.LG, Bottom: spacing.LG},
			ChildGap:        spacing.XL,
			LayoutDirection: clay.TOP_TO_BOTTOM,
			ChildAlignment: clay.ChildAlignment{
				X: clay.ALIGN_X_CENTER,
				Y: clay.ALIGN_Y_CENTER,
			},
		},
		BackgroundColor: mainStyle.BackgroundColor,
	}, func() {
		// Container para conteúdo principal
		clay.UI()(clay.ElementDeclaration{
			Id: clay.ID("content-container"),
			Layout: clay.LayoutConfig{
				Sizing: clay.Sizing{
					Width:  clay.SizingPercent(0.8),
					Height: clay.SizingFit(0, 0),
				},
				Padding:         contentStyle.Padding,
				ChildGap:        spacing.LG,
				LayoutDirection: clay.TOP_TO_BOTTOM,
				ChildAlignment: clay.ChildAlignment{
					X: clay.ALIGN_X_CENTER,
				},
			},
			CornerRadius:    clay.CornerRadiusAll(contentStyle.CornerRadius),
			BackgroundColor: contentStyle.BackgroundColor,
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
				clay.Text("Second screen", theme.CreateTextConfig(typography.XLarge, colors.TextPrimary))
			})

			// Container para textos de conteúdo
			clay.UI()(clay.ElementDeclaration{
				Id: clay.ID("text-content"),
				Layout: clay.LayoutConfig{
					ChildGap:        spacing.MD,
					LayoutDirection: clay.TOP_TO_BOTTOM,
					ChildAlignment: clay.ChildAlignment{
						X: clay.ALIGN_X_CENTER,
					},
				},
			}, func() {
				clay.Text("This is the second application screen.", theme.CreateTextConfig(typography.Large, colors.TextSecondary))

				clay.Text("Here you can add any desired content.", theme.CreateTextConfig(typography.Base, colors.TextMuted))

				clay.Text("This structure allows for easy expansion.", theme.CreateTextConfig(typography.Base, colors.TextMuted))

				clay.Text("Clay allows for flexible and responsive layouts.", theme.CreateTextConfig(typography.Base, colors.TextMuted))
			})

			// Container para botões
			clay.UI()(clay.ElementDeclaration{
				Id: clay.ID("buttons-container"),
				Layout: clay.LayoutConfig{
					Padding:         clay.Padding{Left: spacing.SM, Right: spacing.SM, Top: spacing.SM, Bottom: spacing.SM},
					ChildGap:        spacing.MD,
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

				clay.Text(focusInfo, theme.CreateTextConfig(typography.XSmall, colors.TextMuted))
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

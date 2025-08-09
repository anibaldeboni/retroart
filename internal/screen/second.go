package screen

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"

	"retroart-sdl2/internal/ui"
)

type Second struct {
	screenMgr  *Manager
	claySystem *ui.ClayLayoutSystem
}

func NewSecond(screenMgr *Manager, renderer *sdl.Renderer, font *ttf.Font) *Second {
	screen := &Second{
		screenMgr: screenMgr,
	}

	// Criar sistema Clay para esta tela
	screen.claySystem = ui.NewClayLayoutSystem(renderer, font)

	return screen
}

func (ss *Second) Update() {
	// Lógica de atualização se necessária
}

func (ss *Second) Render(renderer *sdl.Renderer) {
	// Iniciar layout Clay
	ss.claySystem.BeginLayout()

	// Container principal
	ss.claySystem.CreateContainer("main", ui.ContainerConfig{
		Sizing:          ui.DefaultContainerConfig().Sizing,
		Padding:         ui.DefaultContainerConfig().Padding,
		ChildGap:        30,
		LayoutDirection: ui.DefaultContainerConfig().LayoutDirection,
		BackgroundColor: ui.DefaultContainerConfig().BackgroundColor,
	}, func() {
		// Título
		ss.claySystem.CreateText("Segunda Tela", ui.TextConfig{
			FontSize:  28,
			TextColor: ui.DefaultTextConfig().TextColor,
		})

		// Container para conteúdo
		ss.claySystem.CreateContainer("content", ui.DefaultContainerConfig(), func() {
			// Textos de conteúdo
			ss.claySystem.CreateText("Esta é a segunda tela da aplicação.", ui.TextConfig{
				FontSize:  18,
				TextColor: ui.DefaultTextConfig().TextColor,
			})

			ss.claySystem.CreateText("Aqui você pode adicionar qualquer conteúdo desejado.", ui.TextConfig{
				FontSize:  16,
				TextColor: ui.DefaultTextConfig().TextColor,
			})

			ss.claySystem.CreateText("Esta estrutura permite fácil expansão.", ui.TextConfig{
				FontSize:  16,
				TextColor: ui.DefaultTextConfig().TextColor,
			})

			ss.claySystem.CreateText("O Clay permite layouts flexíveis e responsivos.", ui.TextConfig{
				FontSize:  16,
				TextColor: ui.DefaultTextConfig().TextColor,
			})
		})

		// Botão de voltar
		backButtonConfig := ui.DefaultButtonConfig()

		ss.claySystem.CreateButton("back-button", "Voltar", backButtonConfig, func() {
			ss.screenMgr.SetCurrentScreen("home")
		})
	})

	// Finalizar e renderizar
	ss.claySystem.Render()
}

func (ss *Second) HandleInput(input InputType) {
	switch input {
	case InputBack, InputConfirm:
		ss.screenMgr.SetCurrentScreen("home")
	}
}

func (ss *Second) OnEnter() {
	log.Println("Entering Second screen")
}

func (ss *Second) OnExit() {
	log.Println("Exiting Second screen")
}

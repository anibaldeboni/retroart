package app

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"

	"retroart-sdl2/internal/core"
	"retroart-sdl2/internal/input"
	"retroart-sdl2/internal/screen"
	"retroart-sdl2/internal/ui"
)

type App struct {
	window    *sdl.Window
	renderer  *sdl.Renderer
	running   bool
	screenMgr *screen.Manager
}

func New() *App {
	return &App{}
}

func (app *App) Init() error {
	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_JOYSTICK | sdl.INIT_GAMECONTROLLER); err != nil {
		return fmt.Errorf("error initializing SDL: %v", err)
	}

	if err := ttf.Init(); err != nil {
		return fmt.Errorf("error initializing TTF: %v", err)
	}

	window, err := sdl.CreateWindow(
		"RetroArt",
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		core.WINDOW_WIDTH,
		core.WINDOW_HEIGHT,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		return fmt.Errorf("error creating window: %v", err)
	}
	app.window = window

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		return fmt.Errorf("error creating renderer: %v", err)
	}
	app.renderer = renderer
	layout := ui.NewLayout(renderer)

	app.screenMgr = screen.NewManager(layout)

	app.screenMgr.AddScreen("home", screen.NewHome())
	app.screenMgr.AddScreen("second", screen.NewSecond())

	app.screenMgr.SetCurrentScreen("home")

	app.running = true
	return nil
}

func (app *App) Run() {
	targetFrameTime := uint64(1000 / core.FPS) // ms por frame
	inputCh := input.Initialize()

	for app.running {
		frameStart := sdl.GetTicks64()

		app.handleEvents(inputCh)
		app.update()
		app.render()

		// Controle de FPS
		frameTime := sdl.GetTicks64() - frameStart
		if frameTime < targetFrameTime {
			sdl.Delay(uint32(targetFrameTime - frameTime))
		}
	}
}

func (app *App) handleEvents(inputCh <-chan input.InputEvent) {
	select {
	case inputEvent := <-inputCh:
		if inputEvent.Pressed {
			// Passar diretamente o InputType sem conversão desnecessária
			app.screenMgr.HandleInput(inputEvent.Type)
		}
	default:
	}

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			app.running = false
		}
	}
}

func (app *App) update() {
	app.screenMgr.Update()
}

func (app *App) render() {
	app.renderer.SetDrawColor(0, 0, 0, 255)
	app.renderer.Clear()
	app.screenMgr.Render()
	app.renderer.Present()
}

func (app *App) Cleanup() {
	if app.renderer != nil {
		app.renderer.Destroy()
	}
	if app.window != nil {
		app.window.Destroy()
	}
	ttf.Quit()
	sdl.Quit()
}

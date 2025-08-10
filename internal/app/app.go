package app

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"

	"retroart-sdl2/internal/core"
	"retroart-sdl2/internal/screen"
	"retroart-sdl2/internal/ui"
)

type App struct {
	window    *sdl.Window
	renderer  *sdl.Renderer
	font      *ttf.Font
	running   bool
	screenMgr *screen.Manager
}

func New() *App {
	return &App{}
}

func (app *App) Init() error {
	// Inicializar SDL
	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_JOYSTICK | sdl.INIT_GAMECONTROLLER); err != nil {
		return fmt.Errorf("erro ao inicializar SDL: %v", err)
	}

	// Inicializar TTF
	if err := ttf.Init(); err != nil {
		return fmt.Errorf("erro ao inicializar TTF: %v", err)
	}

	// Criar janela
	window, err := sdl.CreateWindow(
		"[CLAY + SDL2] RetroArt - Renderização Ativa!",
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		core.WINDOW_WIDTH,
		core.WINDOW_HEIGHT,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar janela: %v", err)
	}
	app.window = window

	// Criar renderer
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		return fmt.Errorf("erro ao criar renderer: %v", err)
	}
	app.renderer = renderer

	// Carregar fonte
	font, err := ttf.OpenFont("assets/DejaVuSansCondensed.ttf", 16)
	if err != nil {
		// Se não conseguir carregar a fonte personalizada, usar uma fonte padrão do sistema
		fmt.Printf("Aviso: Não foi possível carregar fonte personalizada: %v\n", err)
		// Tentar carregar uma fonte do sistema (macOS)
		font, err = ttf.OpenFont("/System/Library/Fonts/Helvetica.ttc", 16)
		if err != nil {
			return fmt.Errorf("erro ao carregar fonte do sistema: %v", err)
		}
	}
	app.font = font

	// Inicializar Clay globalmente
	if err := ui.InitializeClayGlobally(); err != nil {
		return fmt.Errorf("erro ao inicializar Clay: %v", err)
	}

	// Configurar fonte global para medição de texto
	ui.SetGlobalFont(font)

	// Inicializar gerenciador de telas
	app.screenMgr = screen.NewManager(app.renderer, app.font)

	// Adicionar telas
	app.screenMgr.AddScreen("home", screen.NewHome(app.screenMgr, app.renderer, app.font))
	app.screenMgr.AddScreen("second", screen.NewSecond(app.screenMgr, app.renderer, app.font))

	// Definir tela inicial
	app.screenMgr.SetCurrentScreen("home")

	app.running = true
	return nil
}

func (app *App) Run() {
	targetFrameTime := uint64(1000 / core.FPS) // ms por frame

	for app.running {
		frameStart := sdl.GetTicks64()

		// Processar eventos
		app.handleEvents()

		// Atualizar
		app.update()

		// Renderizar
		app.render()

		// Controle de FPS
		frameTime := sdl.GetTicks64() - frameStart
		if frameTime < targetFrameTime {
			sdl.Delay(uint32(targetFrameTime - frameTime))
		}
	}
}

func (app *App) handleEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			app.running = false
		case *sdl.KeyboardEvent:
			if e.Type == sdl.KEYDOWN {
				app.screenMgr.HandleInput(e.Keysym.Sym)
			}
		}
	}
}

func (app *App) update() {
	app.screenMgr.Update()
}

func (app *App) render() {
	// Limpar tela
	app.renderer.SetDrawColor(0, 0, 0, 255)
	app.renderer.Clear()

	// Renderizar tela atual
	app.screenMgr.Render(app.renderer)

	// Apresentar frame
	app.renderer.Present()
}

func (app *App) Cleanup() {
	if app.font != nil {
		app.font.Close()
	}
	if app.renderer != nil {
		app.renderer.Destroy()
	}
	if app.window != nil {
		app.window.Destroy()
	}
	ttf.Quit()
	sdl.Quit()
}

func (app *App) SetRunning(running bool) {
	app.running = running
}

package ui

import (
	"errors"
	"fmt"
	"log"
	"retroart-sdl2/internal/input"
	"retroart-sdl2/internal/theme"
	"sync"
	"unsafe"

	"github.com/TotallyGamerJet/clay"
	claysdl2 "github.com/TotallyGamerJet/clay/renderers/sdl2"
	"github.com/veandco/go-sdl2/sdl"

	"retroart-sdl2/internal/core"
)

// Layout gerencia o sistema de layout Clay como um Singleton
type Layout struct {
	renderer         *sdl.Renderer
	clayContext      *clay.Context
	clayArena        clay.Arena
	arenaResetOffset uint64
	fontSystem       *theme.FontSystem
	spatialNav       *SpatialNavigation
}

var (
	instance *Layout
	once     sync.Once
)

// GetLayout retorna a instância singleton do ClayLayoutSystem
func GetLayout() *Layout {
	return instance
}

// NewLayout cria ou retorna a instância singleton do sistema de layout Clay
func NewLayout(renderer *sdl.Renderer, fontSystem *theme.FontSystem) (*Layout, error) {
	var initError error

	once.Do(func() {
		instance = &Layout{
			renderer:   renderer,
			fontSystem: fontSystem,
			spatialNav: NewSpatialNavigation(),
		}

		if err := instance.initializeClay(); err != nil {
			initError = fmt.Errorf("failed to initialize Clay: %w", err)
			return
		}

		err := instance.configureMeasureTextFunction()
		if err != nil {
			initError = fmt.Errorf("failed to configure text measurement function: %w", err)
			return
		}

		log.Println("Layout created successfully")
	})

	if initError != nil {
		return nil, initError
	}

	if instance != nil && instance.renderer != renderer {
		instance.renderer = renderer
		log.Println("Layout renderer updated")
	}

	return instance, nil
}

// initializeClay inicializa o sistema Clay internamente
func (l *Layout) initializeClay() error {
	// Use MinMemorySize to calculate the correct arena size
	arenaSize := clay.MinMemorySize()
	memory := make([]byte, arenaSize)
	l.clayArena = clay.CreateArenaWithCapacityAndMemory(memory)

	log.Printf("Creating arena with size: %d bytes", arenaSize)

	// Layout dimensions
	dimensions := clay.Dimensions{
		Width:  float32(core.WINDOW_WIDTH),
		Height: float32(core.WINDOW_HEIGHT),
	}

	// Initialize Clay
	l.clayContext = clay.Initialize(l.clayArena, dimensions, clay.ErrorHandler{})
	if l.clayContext == nil {
		return errors.New("Clay.Initialize returned nil context")
	}

	// Check if current context was set correctly
	currentContext := clay.GetCurrentContext()
	if currentContext == nil {
		return errors.New("Clay current context is nil after initialization")
	}

	log.Printf("Current context address: %p, clay context address: %p", currentContext, l.clayContext)

	// Capture arena reset offset after initialization
	l.arenaResetOffset = l.clayArena.NextAllocation

	log.Printf("Clay initialized successfully, arenaResetOffset: %d", l.arenaResetOffset)
	return nil
}

// configureMeasureTextFunction configures Clay's text measurement after fonts are initialized
func (l *Layout) configureMeasureTextFunction() error {
	clayFontsPtr := l.fontSystem.GetFonts()
	if clayFontsPtr == nil {
		log.Printf("Error: Cannot configure text measurement function - no fonts available")
		return errors.New("no fonts available")
	}

	clay.SetMeasureTextFunction(claysdl2.MeasureText, unsafe.Pointer(clayFontsPtr))
	log.Printf("Configured text measurement function with stable pointer to %d fonts", len(*clayFontsPtr))
	return nil
}

// Render executes the layout rendering process for the Layout instance.
// It first ensures the rendering context is valid, prepares the layout arena,
// and then begins the layout process using the clay package. If a screenRenderFunc
// is provided, it is called to perform custom rendering logic. After ending the layout,
// it updates spatial navigation and renders the resulting commands to SDL.
//
// screenRenderFunc: An optional function that performs additional rendering logic
//
//	during the layout process.
func (l *Layout) Render(screenRenderFunc func()) {
	if !l.ensureValidContext() {
		return
	}

	l.prepareArena()

	clay.BeginLayout()
	if screenRenderFunc != nil {
		screenRenderFunc()
	}
	commands := clay.EndLayout()

	l.updateSpatialNavigation(commands)
	l.renderToSDL(commands)
}

// ensureValidContext verifica e configura um contexto Clay válido
func (l *Layout) ensureValidContext() bool {
	currentContext := clay.GetCurrentContext()
	if currentContext != nil {
		return true
	}

	if l.clayContext == nil {
		log.Println("Clay context is nil, aborting Render")
		return false
	}

	clay.SetCurrentContext(l.clayContext)
	if clay.GetCurrentContext() == nil {
		log.Println("Failed to set current context, aborting Render")
		return false
	}

	return true
}

// prepareArena prepara a arena Clay para um novo ciclo de renderização
func (l *Layout) prepareArena() {
	l.clayArena.NextAllocation = l.arenaResetOffset
}

// updateSpatialNavigation atualiza o sistema de navegação espacial com os commands de renderização
func (l *Layout) updateSpatialNavigation(commands clay.RenderCommandArray) {
	if l.spatialNav != nil {
		l.spatialNav.UpdateLayout(commands)
	}
}

// renderToSDL renderiza os commands Clay usando SDL2
func (l *Layout) renderToSDL(commands clay.RenderCommandArray) {
	clayFonts := l.fontSystem.GetFonts()
	if err := claysdl2.ClayRender(l.renderer, commands, *clayFonts); err != nil {
		log.Printf("Error rendering Clay commands: %v", err)
	}
}

// GetSpatialNavigation retorna o sistema de navegação espacial
func (l *Layout) GetSpatialNavigation() *SpatialNavigation {
	return l.spatialNav
}

// RegisterFocusable registra um widget focável no sistema de navegação espacial
func (l *Layout) RegisterFocusable(widget Focusable) {
	if l.spatialNav != nil {
		l.spatialNav.RegisterFocusable(widget)
	}
}

// UnregisterFocusable remove um widget focável do sistema de navegação espacial
func (l *Layout) UnregisterFocusable(id string) {
	if l.spatialNav != nil {
		l.spatialNav.UnregisterFocusable(id)
	}
}

// HandleSpatialInput processa input de navegação espacial
func (l *Layout) HandleSpatialInput(inputType input.InputType) bool {
	if l.spatialNav != nil {
		return l.spatialNav.HandleInput(inputType)
	}
	return false
}

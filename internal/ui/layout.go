package ui

import (
	"errors"
	"fmt"
	"log"
	"retroart-sdl2/internal/input"
	"sync"
	"unsafe"

	"github.com/TotallyGamerJet/clay"
	claysdl2 "github.com/TotallyGamerJet/clay/renderers/sdl2"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"

	"retroart-sdl2/internal/core"
)

// Layout gerencia o sistema de layout Clay como um Singleton
type Layout struct {
	renderer         *sdl.Renderer
	clayContext      *clay.Context
	clayArena        clay.Arena
	arenaResetOffset uint64
	fontCache        map[uint16]*ttf.Font
	clayFonts        []claysdl2.Font // Official renderer font format
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
func NewLayout(renderer *sdl.Renderer) (*Layout, error) {
	var initError error

	once.Do(func() {
		instance = &Layout{
			renderer:   renderer,
			fontCache:  make(map[uint16]*ttf.Font),
			spatialNav: NewSpatialNavigation(),
		}

		// Inicializar Clay e FontSystem
		if err := instance.initializeClay(); err != nil {
			initError = fmt.Errorf("failed to initialize Clay: %w", err)
			return
		}

		if err := instance.initializeFontSystem(); err != nil {
			initError = fmt.Errorf("failed to initialize font system: %w", err)
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

	// Configure text measurement function using official Clay SDL2 renderer
	clay.SetMeasureTextFunction(claysdl2.MeasureText, unsafe.Pointer(&l.clayFonts))

	// Capture arena reset offset after initialization
	l.arenaResetOffset = l.clayArena.NextAllocation

	log.Printf("Clay initialized successfully, arenaResetOffset: %d", l.arenaResetOffset)
	return nil
}

// initializeFontSystem inicializa o sistema de fontes internamente
func (l *Layout) initializeFontSystem() error {
	// Carregar fontes com diferentes tamanhos
	fontSizes := []uint16{10, 12, 14, 16, 18, 20, 24, 28, 32}

	// Initialize clayFonts slice for official renderer
	l.clayFonts = make([]claysdl2.Font, 0, len(fontSizes))

	for _, size := range fontSizes {
		font, err := l.loadFontWithSize(size)
		if err != nil {
			log.Printf("Warning: Could not load font size %d: %v", size, err)
			continue
		}
		l.fontCache[size] = font

		// Add to clayFonts for official renderer (use array index as FontId)
		l.clayFonts = append(l.clayFonts, claysdl2.Font{
			FontId: uint32(len(l.clayFonts)),
			Font:   font,
		})

		log.Printf("Loaded font size: %d (FontId: %d)", size, len(l.clayFonts)-1)
	}

	if len(l.fontCache) == 0 {
		return errors.New("failed to load any fonts")
	}

	log.Printf("Font system initialized with %d font sizes", len(l.fontCache))
	return nil
}

// loadFontWithSize carrega uma fonte com tamanho específico
func (l *Layout) loadFontWithSize(size uint16) (*ttf.Font, error) {
	// Lista de possíveis caminhos de fontes no sistema
	fontPaths := []string{
		"assets/DejaVuSansCondensed.ttf",
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
		"/usr/share/fonts/TTF/DejaVuSans.ttf",
		"/System/Library/Fonts/Helvetica.ttc",
		"/usr/share/fonts/liberation/LiberationSans-Regular.ttf",
	}

	for _, fontPath := range fontPaths {
		font, err := ttf.OpenFont(fontPath, int(size))
		if err == nil {
			log.Printf("Successfully loaded font from: %s (size %d)", fontPath, size)
			return font, nil
		}
	}

	return nil, fmt.Errorf("could not load any font for size %d", size)
}

// GetFontForSize retorna a fonte para um tamanho específico
func (l *Layout) GetFontForSize(size uint16) *ttf.Font {
	if l.fontCache == nil {
		return nil
	}

	if font, exists := l.fontCache[size]; exists {
		return font
	}

	var closestSize uint16
	var minDiff uint16 = 1000

	for cachedSize := range l.fontCache {
		var diff uint16
		if cachedSize > size {
			diff = cachedSize - size
		} else {
			diff = size - cachedSize
		}

		if diff < minDiff {
			minDiff = diff
			closestSize = cachedSize
		}
	}

	return l.fontCache[closestSize]
}

// Render executa o ciclo completo de renderização Clay
func (l *Layout) Render(screenRenderFunc func()) {
	// Verificar se o contexto atual existe
	currentContext := clay.GetCurrentContext()
	if currentContext == nil {
		log.Println("Clay current context is nil, trying to set clay context")
		if l.clayContext != nil {
			clay.SetCurrentContext(l.clayContext)
			currentContext = clay.GetCurrentContext()
			if currentContext == nil {
				log.Println("Failed to set current context, aborting Render")
				return
			}
		} else {
			log.Println("Clay context is nil, aborting Render")
			return
		}
	}

	// Resetar arena para o offset correto como no clay.h:2155
	l.clayArena.NextAllocation = l.arenaResetOffset
	log.Printf("Arena reset to offset: %d", l.arenaResetOffset)

	clay.BeginLayout()
	log.Println("clay.BeginLayout() completed")

	// Executar a função de renderização da tela
	if screenRenderFunc != nil {
		screenRenderFunc()
	}

	commands := clay.EndLayout()
	log.Printf("clay.EndLayout() completed, got %d commands", commands.Length)

	// Atualizar navegação espacial com os comandos de renderização
	if l.spatialNav != nil {
		l.spatialNav.UpdateLayout(commands)
	}

	log.Printf("ClayRender called with %d commands", commands.Length)
	err := claysdl2.ClayRender(l.renderer, commands, l.clayFonts)
	if err != nil {
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

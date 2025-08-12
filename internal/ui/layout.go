package ui

import (
	"errors"
	"fmt"
	"log"
	"retroart-sdl2/internal/input"
	"sync"
	"unsafe"

	"github.com/TotallyGamerJet/clay"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"

	"retroart-sdl2/internal/core"
)

// Layout gerencia o sistema de layout Clay como um Singleton
type Layout struct {
	renderer         *sdl.Renderer
	enabled          bool
	isActive         bool
	clayContext      *clay.Context
	clayArena        clay.Arena
	arenaResetOffset uint64
	fontCache        map[uint16]*ttf.Font
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
func NewLayout(renderer *sdl.Renderer) *Layout {
	once.Do(func() {
		instance = &Layout{
			renderer:   renderer,
			enabled:    false,
			isActive:   false,
			fontCache:  make(map[uint16]*ttf.Font),
			spatialNav: NewSpatialNavigation(),
		}

		// Inicializar Clay e FontSystem
		if err := instance.initializeClay(); err != nil {
			log.Printf("Failed to initialize Clay: %v", err)
			return
		}

		if err := instance.initializeFontSystem(); err != nil {
			log.Printf("Failed to initialize font system: %v", err)
			return
		}

		instance.enabled = true
		log.Println("ClayLayoutSystem singleton created successfully")
	})

	// Se já existe uma instância, apenas atualizar o renderer se necessário
	if instance != nil && instance.renderer != renderer {
		instance.renderer = renderer
		log.Println("ClayLayoutSystem renderer updated")
	}

	return instance
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

	// Configure text measurement function
	// Create dummy userData to avoid nil conversion
	dummyUserData := make([]byte, 1)
	clay.SetMeasureTextFunction(l.measureTextWithFont, unsafe.Pointer(&dummyUserData[0]))

	// Capture arena reset offset after initialization
	l.arenaResetOffset = l.clayArena.NextAllocation

	log.Printf("Clay initialized successfully, arenaResetOffset: %d", l.arenaResetOffset)
	return nil
}

// initializeFontSystem inicializa o sistema de fontes internamente
func (l *Layout) initializeFontSystem() error {
	// Carregar fontes com diferentes tamanhos
	fontSizes := []uint16{10, 12, 14, 16, 18, 20, 24, 28, 32}

	for _, size := range fontSizes {
		font, err := l.loadFontWithSize(size)
		if err != nil {
			log.Printf("Warning: Could not load font size %d: %v", size, err)
			continue
		}
		l.fontCache[size] = font
		log.Printf("Loaded font size: %d", size)
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

// measureTextWithFont retorna as dimensões do texto renderizado com uma fonte específica
func (l *Layout) measureTextWithFont(text clay.StringSlice, config *clay.TextElementConfig, userData unsafe.Pointer) clay.Dimensions {
	// Converter clay.StringSlice para string
	textString := text.String()

	// Se string vazia, retornar dimensões zero
	if textString == "" {
		return clay.Dimensions{Width: 0, Height: 0}
	}

	// Obter fonte do tamanho especificado
	font := l.GetFontForSize(config.FontSize)
	if font == nil {
		// Fallback para tamanhos padrão
		log.Printf("No font found for size %d, using fallback", config.FontSize)
		return clay.Dimensions{Width: float32(len(textString) * 8), Height: 16}
	}

	// Tentar medir o texto completo primeiro
	w, h, err := font.SizeUTF8(textString)
	if err != nil {
		log.Printf("Error measuring text '%s': %v", textString, err)
		return l.calculateFallbackDimensions(textString, config.FontSize)
	}

	// Se retornou largura zero, pode ser problema de caracteres Unicode não suportados
	if w == 0 && len(textString) > 0 {
		log.Printf("TTF_SizeUTF8 returned zero width for text '%s', trying character-by-character", textString)
		return l.measureTextCharByChar(textString, font, config.FontSize)
	}

	return clay.Dimensions{Width: float32(w), Height: float32(h)}
}

// measureTextCharByChar mede texto caractere por caractere para contornar problemas Unicode
func (l *Layout) measureTextCharByChar(text string, font *ttf.Font, fontSize uint16) clay.Dimensions {
	totalWidth := 0
	maxHeight := 0

	// Percorre diretamente a string para lidar corretamente com UTF-8
	for _, r := range text {
		charStr := string(r)

		// Tentar medir o caractere individual
		w, h, err := font.SizeUTF8(charStr)
		if err != nil {
			// Se falhar, usar fallback baseado no tamanho da fonte
			w = int(fontSize) / 2 // Aproximação baseada no tamanho da fonte
			h = int(fontSize)
		}

		// Se ainda retornar largura zero, usar fallback
		if w == 0 {
			w = int(fontSize) / 2
		}
		if h == 0 {
			h = int(fontSize)
		}

		totalWidth += w
		if h > maxHeight {
			maxHeight = h
		}
	}

	return clay.Dimensions{Width: float32(totalWidth), Height: float32(maxHeight)}
}

// calculateFallbackDimensions calcula dimensões de fallback baseadas no tamanho da fonte
func (l *Layout) calculateFallbackDimensions(text string, fontSize uint16) clay.Dimensions {
	// Aproximação: largura média de caractere é ~60% do tamanho da fonte
	// Altura é aproximadamente o tamanho da fonte
	runeCount := len([]rune(text)) // Conta corretamente caracteres UTF-8
	avgCharWidth := float32(fontSize) * 0.6
	totalWidth := avgCharWidth * float32(runeCount)
	height := float32(fontSize)

	log.Printf("Using fallback dimensions for text '%s': W=%.2f H=%.2f", text, totalWidth, height)
	return clay.Dimensions{Width: totalWidth, Height: height}
}

// BeginLayout inicia um novo frame de layout seguindo o padrão do videodemo.go
func (l *Layout) BeginLayout() {
	if !l.enabled {
		log.Println("Clay not enabled, skipping BeginLayout")
		return
	}

	// Verificar se o contexto atual existe
	currentContext := clay.GetCurrentContext()
	if currentContext == nil {
		log.Println("Clay current context is nil, trying to set clay context")
		if l.clayContext != nil {
			clay.SetCurrentContext(l.clayContext)
			currentContext = clay.GetCurrentContext()
			if currentContext == nil {
				log.Println("Failed to set current context, aborting BeginLayout")
				return
			}
		} else {
			log.Println("Clay context is nil, aborting BeginLayout")
			return
		}
	}

	// Resetar arena para o offset correto como no clay.h:2155
	l.clayArena.NextAllocation = l.arenaResetOffset
	log.Printf("Arena reset to offset: %d", l.arenaResetOffset)

	l.isActive = true
	clay.BeginLayout()
	log.Println("clay.BeginLayout() completed")
}

// EndLayout finaliza o layout e retorna os comandos de renderização
func (l *Layout) EndLayout() clay.RenderCommandArray {
	if !l.enabled || !l.isActive {
		log.Println("Clay not enabled or not active, returning empty RenderCommandArray")
		return clay.RenderCommandArray{}
	}
	l.isActive = false

	commands := clay.EndLayout()
	log.Printf("clay.EndLayout() completed, got %d commands", commands.Length)

	// Atualizar navegação espacial com os comandos de renderização
	if l.spatialNav != nil {
		l.spatialNav.UpdateLayout(commands)
	}

	return commands
}

// Render executa o ciclo completo de renderização
func (l *Layout) Render() {
	if !l.enabled {
		return
	}

	commands := l.EndLayout()

	err := l.RenderClayCommands(commands)
	if err != nil {
		log.Printf("Error rendering Clay commands: %v", err)
	}
}

// RenderClayCommands renderiza os comandos Clay no renderer SDL2
func (l *Layout) RenderClayCommands(commands clay.RenderCommandArray) error {
	if !l.enabled {
		return nil
	}

	log.Printf("RenderClayCommands called with %d commands", commands.Length)

	// Iterar pelos comandos de renderização
	for i := int32(0); i < commands.Length; i++ {
		command := clay.RenderCommandArray_Get(&commands, i)

		switch command.CommandType {
		case clay.RENDER_COMMAND_TYPE_RECTANGLE:
			config := &command.RenderData.Rectangle
			log.Printf("ClayRectangleCommand %d: BoundingBox X=%.2f Y=%.2f W=%.2f H=%.2f, Color R=%.0f G=%.0f B=%.0f A=%.0f",
				i, command.BoundingBox.X, command.BoundingBox.Y, command.BoundingBox.Width, command.BoundingBox.Height,
				config.BackgroundColor.R, config.BackgroundColor.G, config.BackgroundColor.B, config.BackgroundColor.A)
			err := l.renderRectangle(command)
			if err != nil {
				log.Printf("Error rendering rectangle: %v", err)
				return err
			}
		case clay.RENDER_COMMAND_TYPE_TEXT:
			err := l.renderText(command)
			if err != nil {
				log.Printf("Error rendering text: %v", err)
				return err
			}
		case clay.RENDER_COMMAND_TYPE_BORDER:
			log.Printf("Rendering border command")
			err := l.renderBorder(command)
			if err != nil {
				return err
			}
		case clay.RENDER_COMMAND_TYPE_SCISSOR_START:
			err := l.renderScissorStart(command)
			if err != nil {
				return err
			}
		case clay.RENDER_COMMAND_TYPE_SCISSOR_END:
			err := l.renderScissorEnd()
			if err != nil {
				return err
			}
		default:
			log.Printf("Command not implemented: %v", command.CommandType)
		}
	}

	return nil
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

package ui

import (
	"errors"
	"log"

	"unsafe"

	"github.com/TotallyGamerJet/clay"
	"github.com/veandco/go-sdl2/ttf"

	"retroart-sdl2/internal/core"
)

// Variables para controlar inicialização global
var globalClayContext *clay.Context
var globalClayArena clay.Arena
var arenaResetOffset uint64
var globalFont *ttf.Font // Adicionar referência global da fonte

// measureTextWithFont é uma função de medição de texto que usa a fonte TTF real quando disponível
func measureTextWithFont(text clay.StringSlice, config *clay.TextElementConfig, userData unsafe.Pointer) clay.Dimensions {
	if text.Length == 0 {
		return clay.Dimensions{
			Width:  1.0,  // Evitar zero width
			Height: 16.0, // Altura padrão
		}
	}

	textString := text.String()
	if textString == "" {
		return clay.Dimensions{
			Width:  1.0,
			Height: 16.0,
		}
	}

	// Tentar usar a fonte TTF global se disponível
	if globalFont != nil {
		// Usar TTF.SizeUTF8 para medição real do texto
		w, h, err := globalFont.SizeUTF8(textString)
		if err == nil && w > 0 && h > 0 {
			return clay.Dimensions{
				Width:  float32(w),
				Height: float32(h),
			}
		}
	}

	// Fallback para estimativa se TTF não disponível
	fontSize := float32(16) // tamanho padrão
	if config != nil && config.FontSize > 0 {
		fontSize = float32(config.FontSize)
	}

	// Calculate dimensions based on font size
	// Use a more realistic estimation for character width
	charWidth := fontSize * 0.6 // Approximate ratio for monospace fonts
	textLength := float32(text.Length)

	width := textLength * charWidth
	height := fontSize * 1.2 // Add space for ascenders/descenders

	// Ensure dimensions are at least 1 pixel
	if width < 1.0 {
		width = 1.0
	}
	if height < 1.0 {
		height = 1.0
	}

	return clay.Dimensions{
		Width:  width,
		Height: height,
	}
}

// InitializeClayGlobally inicializa Clay uma única vez globalmente
func InitializeClayGlobally() error {
	if globalClayContext != nil {
		log.Println("Clay already initialized globally")
		return nil
	}

	// Use MinMemorySize to calculate the correct arena size
	arenaSize := clay.MinMemorySize()
	memory := make([]byte, arenaSize)
	globalClayArena = clay.CreateArenaWithCapacityAndMemory(memory)

	log.Printf("Creating arena with size: %d bytes", arenaSize)

	// Layout dimensions
	dimensions := clay.Dimensions{
		Width:  float32(core.WINDOW_WIDTH),
		Height: float32(core.WINDOW_HEIGHT),
	}

	// Initialize Clay globally
	globalClayContext = clay.Initialize(globalClayArena, dimensions, clay.ErrorHandler{})
	if globalClayContext == nil {
		return errors.New("Clay.Initialize returned nil context")
	}

	// Check if current context was set correctly
	currentContext := clay.GetCurrentContext()
	if currentContext == nil {
		return errors.New("Clay current context is nil after initialization")
	}

	log.Printf("Current context address: %p, global context address: %p", currentContext, globalClayContext)

	// Configure text measurement function
	// Create dummy userData to avoid nil conversion
	dummyUserData := make([]byte, 1)
	clay.SetMeasureTextFunction(measureTextWithFont, unsafe.Pointer(&dummyUserData[0]))

	// Capture arena reset offset after initialization
	arenaResetOffset = globalClayArena.NextAllocation

	log.Printf("Clay initialized globally successfully, arenaResetOffset: %d", arenaResetOffset)
	return nil
}

// SetGlobalFont configura a fonte global para medição de texto
func SetGlobalFont(font *ttf.Font) {
	globalFont = font
	log.Println("Global font set for text measurement")
}

package ui

import (
	"errors"
	"fmt"
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

// Sistema de fontes
var fontCache map[uint16]*ttf.Font // Cache de fontes por tamanho

// measureTextWithFont retorna as dimensões do texto renderizado com uma fonte específica
func measureTextWithFont(text clay.StringSlice, config *clay.TextElementConfig, userData unsafe.Pointer) clay.Dimensions {
	// Converter clay.StringSlice para string
	textString := text.String()

	// Se string vazia, retornar dimensões zero
	if textString == "" {
		return clay.Dimensions{Width: 0, Height: 0}
	}

	// Obter fonte do tamanho especificado
	font := GetFontForSize(config.FontSize)
	if font == nil {
		// Fallback para tamanhos padrão
		log.Printf("No font found for size %d, using fallback", config.FontSize)
		return clay.Dimensions{Width: float32(len(textString) * 8), Height: 16}
	}

	// Tentar medir o texto completo primeiro
	w, h, err := font.SizeUTF8(textString)
	if err != nil {
		log.Printf("Error measuring text '%s': %v", textString, err)
		return calculateFallbackDimensions(textString, config.FontSize)
	}

	// Se retornou largura zero, pode ser problema de caracteres Unicode não suportados
	if w == 0 && len(textString) > 0 {
		log.Printf("TTF_SizeUTF8 returned zero width for text '%s', trying character-by-character", textString)
		return measureTextCharByChar(textString, font, config.FontSize)
	}

	return clay.Dimensions{Width: float32(w), Height: float32(h)}
}

// measureTextCharByChar mede texto caractere por caractere para contornar problemas Unicode
func measureTextCharByChar(text string, font *ttf.Font, fontSize uint16) clay.Dimensions {
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
func calculateFallbackDimensions(text string, fontSize uint16) clay.Dimensions {
	// Aproximação: largura média de caractere é ~60% do tamanho da fonte
	// Altura é aproximadamente o tamanho da fonte
	runeCount := len([]rune(text)) // Conta corretamente caracteres UTF-8
	avgCharWidth := float32(fontSize) * 0.6
	totalWidth := avgCharWidth * float32(runeCount)
	height := float32(fontSize)

	log.Printf("Using fallback dimensions for text '%s': W=%.2f H=%.2f", text, totalWidth, height)
	return clay.Dimensions{Width: totalWidth, Height: height}
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

// InitializeFontSystem inicializa o sistema de fontes
func InitializeFontSystem() error {
	fontCache = make(map[uint16]*ttf.Font)

	// Carregar fontes com diferentes tamanhos
	fontSizes := []uint16{10, 12, 14, 16, 18, 20, 24, 28, 32}

	for _, size := range fontSizes {
		font, err := loadFontWithSize(size)
		if err != nil {
			log.Printf("Warning: Could not load font size %d: %v", size, err)
			continue
		}
		fontCache[size] = font
		log.Printf("Loaded font size: %d", size)
	}

	if len(fontCache) == 0 {
		return errors.New("failed to load any fonts")
	}

	log.Printf("Font system initialized with %d font sizes", len(fontCache))
	return nil
}

// loadFontWithSize carrega uma fonte com tamanho específico
func loadFontWithSize(size uint16) (*ttf.Font, error) {
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
func GetFontForSize(size uint16) *ttf.Font {
	if fontCache == nil {
		return nil
	}

	// Primeiro, tentar o tamanho exato
	if font, exists := fontCache[size]; exists {
		return font
	}

	// Se não existe, procurar o tamanho mais próximo
	var closestSize uint16
	var minDiff uint16 = 1000

	for cachedSize := range fontCache {
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

	return fontCache[closestSize]
}

package ui

import (
	"fmt"
	"log"

	"unsafe"

	"github.com/TotallyGamerJet/clay"

	"retroart-sdl2/internal/core"
)

// Variables para controlar inicialização global
var globalClayContext *clay.Context
var globalClayArena clay.Arena
var arenaResetOffset uint64

// simpleMeasureText é uma função de medição de texto básica
func simpleMeasureText(text clay.StringSlice, config *clay.TextElementConfig, userData unsafe.Pointer) clay.Dimensions {
	// Função de medição de texto simples - retorna dimensões aproximadas
	// Em um caso real, você usaria TTF para medir o texto real
	charWidth := float32(8)   // largura aproximada de um caractere
	charHeight := float32(16) // altura aproximada de um caractere

	if config != nil && config.FontSize > 0 {
		// Usar o tamanho da fonte se disponível
		charHeight = float32(config.FontSize)
		charWidth = float32(config.FontSize) * 0.6 // relação aproximada largura/altura
	}

	textLength := float32(text.Length)
	return clay.Dimensions{
		Width:  textLength * charWidth,
		Height: charHeight,
	}
}

// InitializeClayGlobally inicializa Clay uma única vez globalmente
func InitializeClayGlobally() error {
	if globalClayContext != nil {
		log.Println("Clay already initialized globally")
		return nil
	}

	// Usar MinMemorySize para calcular o tamanho correto da arena
	arenaSize := clay.MinMemorySize()
	memory := make([]byte, arenaSize)
	globalClayArena = clay.CreateArenaWithCapacityAndMemory(memory)

	log.Printf("Creating arena with size: %d bytes", arenaSize)

	// Dimensões do layout
	dimensions := clay.Dimensions{
		Width:  float32(core.WINDOW_WIDTH),
		Height: float32(core.WINDOW_HEIGHT),
	}

	// Inicializar Clay globalmente
	globalClayContext = clay.Initialize(globalClayArena, dimensions, clay.ErrorHandler{})
	if globalClayContext == nil {
		return fmt.Errorf("Clay.Initialize returned nil context")
	}

	// Verificar se o contexto atual foi definido corretamente
	currentContext := clay.GetCurrentContext()
	if currentContext == nil {
		return fmt.Errorf("Clay current context is nil after initialization")
	}

	log.Printf("Current context address: %p, global context address: %p", currentContext, globalClayContext)

	// Configurar função de medição de texto simples
	// Criar um userData dummy para evitar conversão nil
	dummyUserData := make([]byte, 1)
	clay.SetMeasureTextFunction(simpleMeasureText, unsafe.Pointer(&dummyUserData[0]))

	// Capturar o offset de reset da arena após inicialização
	arenaResetOffset = globalClayArena.NextAllocation

	log.Printf("Clay initialized globally successfully, arenaResetOffset: %d", arenaResetOffset)
	return nil
}

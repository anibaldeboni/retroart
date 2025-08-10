package ui

import (
	"log"

	"github.com/TotallyGamerJet/clay"
)

// Estruturas de configuração simplificadas

type ContainerConfig struct {
	Sizing          clay.Sizing
	Padding         clay.Padding
	ChildGap        uint16
	LayoutDirection clay.LayoutDirection
	BackgroundColor clay.Color
}

type TextConfig struct {
	FontSize  uint16
	TextColor clay.Color
}

// ButtonConfig contém configurações para diferentes estados do botão
type ButtonConfig struct {
	Sizing       clay.Sizing
	Padding      clay.Padding
	TextSize     uint16
	CornerRadius float32

	Normal  ButtonState
	Focused ButtonState
}

// ButtonState define a aparência de um estado específico do botão
type ButtonState struct {
	BackgroundColor clay.Color
	TextColor       clay.Color
}

// Funções helper para criar configurações comuns

// DefaultContainerConfig retorna uma configuração padrão para containers
func DefaultContainerConfig() ContainerConfig {
	return ContainerConfig{
		Sizing: clay.Sizing{
			Width:  clay.SizingPercent(1.0), // Cresce para ocupar espaço disponível no pai
			Height: clay.SizingPercent(1.0), // Ajusta altura baseado no conteúdo
		},
		Padding:         clay.PaddingAll(10),
		ChildGap:        8,
		LayoutDirection: clay.TOP_TO_BOTTOM,
		BackgroundColor: clay.Color{R: 40, G: 42, B: 54, A: 255}, // Cor de fundo padrão
	}
}

// DefaultButtonConfig retorna uma configuração padrão para botões com estado
func DefaultButtonConfig() ButtonConfig {
	return ButtonConfig{
		Sizing: clay.Sizing{
			Width:  clay.SizingFixed(220),
			Height: clay.SizingFixed(45),
		},
		Padding:      clay.Padding{Left: 20, Right: 20, Top: 12, Bottom: 12},
		TextSize:     16,
		CornerRadius: 12,
		Normal: ButtonState{
			BackgroundColor: clay.Color{R: 70, G: 130, B: 220, A: 200},
			TextColor:       clay.Color{R: 220, G: 230, B: 255, A: 255},
		},
		Focused: ButtonState{
			BackgroundColor: clay.Color{R: 30, G: 150, B: 255, A: 255},
			TextColor:       clay.Color{R: 255, G: 255, B: 255, A: 255},
		},
	}
}

// CreateButtonConfig cria uma configuração personalizada para botão com estado
func CreateButtonConfig(normalBg, focusedBg, normalText, focusedText clay.Color) ButtonConfig {
	config := DefaultButtonConfig()
	config.Normal.BackgroundColor = normalBg
	config.Normal.TextColor = normalText
	config.Focused.BackgroundColor = focusedBg
	config.Focused.TextColor = focusedText
	return config
}

// Configurações predefinidas para diferentes tipos de botões
func PrimaryButtonConfig() ButtonConfig {
	return CreateButtonConfig(
		clay.Color{R: 50, G: 120, B: 200, A: 200},  // Azul normal
		clay.Color{R: 30, G: 150, B: 255, A: 255},  // Azul vibrante focado
		clay.Color{R: 220, G: 230, B: 255, A: 255}, // Texto azul claro
		clay.Color{R: 255, G: 255, B: 255, A: 255}, // Texto branco focado
	)
}

func DangerButtonConfig() ButtonConfig {
	return CreateButtonConfig(
		clay.Color{R: 200, G: 60, B: 60, A: 180},   // Vermelho normal
		clay.Color{R: 255, G: 80, B: 80, A: 255},   // Vermelho vibrante focado
		clay.Color{R: 255, G: 200, B: 200, A: 255}, // Texto vermelho claro
		clay.Color{R: 255, G: 255, B: 255, A: 255}, // Texto branco focado
	)
}

func SecondaryButtonConfig() ButtonConfig {
	return CreateButtonConfig(
		clay.Color{R: 120, G: 60, B: 200, A: 180},  // Roxo normal
		clay.Color{R: 150, G: 80, B: 255, A: 255},  // Roxo vibrante focado
		clay.Color{R: 220, G: 200, B: 255, A: 255}, // Texto roxo claro
		clay.Color{R: 255, G: 255, B: 255, A: 255}, // Texto branco focado
	)
}

func DefaultTextConfig() TextConfig {
	return TextConfig{
		FontSize:  16,
		TextColor: clay.Color{R: 255, G: 255, B: 255, A: 255},
	}
}

// CreateContainer cria um container básico
func (cls *ClayLayoutSystem) CreateContainer(id string, config ContainerConfig, children func()) {
	if !cls.enabled || !cls.isActive {
		log.Printf("Clay not enabled or not active, skipping CreateContainer: %s", id)
		return
	}

	log.Printf("Creating container: %s", id)
	clay.UI()(clay.ElementDeclaration{
		Id: clay.ID(id),
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  config.Sizing.Width,
				Height: config.Sizing.Height,
			},
			Padding:         config.Padding,
			ChildGap:        config.ChildGap,
			LayoutDirection: config.LayoutDirection,
		},
		BackgroundColor: config.BackgroundColor,
	}, children)
	log.Printf("Container created successfully: %s", id)
}

// CreateText cria um elemento de texto
func (cls *ClayLayoutSystem) CreateText(text string, config TextConfig) {
	if !cls.enabled || !cls.isActive {
		log.Println("Clay not enabled or not active, skipping CreateText")
		return
	}

	// Usar configuração simplificada de texto
	textConfig := clay.TextElementConfig{
		FontSize:  config.FontSize,
		TextColor: config.TextColor,
	}

	log.Printf("Creating text: %s", text)
	clay.Text(text, &textConfig)
	log.Println("Text created successfully")
}

// CreateButton cria um botão que gerencia seus próprios estados
func (cls *ClayLayoutSystem) CreateButton(id string, text string, config ButtonConfig, isFocused bool, onClick func()) {
	if !cls.enabled || !cls.isActive {
		log.Printf("Clay not enabled or not active, skipping CreateStatefulButton: %s", id)
		return
	}

	// Determinar estado atual baseado no foco
	var currentState ButtonState
	if isFocused {
		currentState = config.Focused
	} else {
		currentState = config.Normal
	}

	log.Printf("Creating new button: %s (focused: %t)", id, isFocused)
	clay.UI()(clay.ElementDeclaration{
		Id: clay.ID(id),
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  config.Sizing.Width,
				Height: config.Sizing.Height,
			},
			Padding: config.Padding,
			ChildAlignment: clay.ChildAlignment{
				X: clay.ALIGN_X_CENTER,
				Y: clay.ALIGN_Y_CENTER,
			},
		},
		CornerRadius:    clay.CornerRadiusAll(config.CornerRadius),
		BackgroundColor: currentState.BackgroundColor,
	}, func() {
		// Texto do botão centralizado com cor do estado atual
		cls.CreateText(text, TextConfig{
			FontSize:  config.TextSize,
			TextColor: currentState.TextColor,
		})
	})
	log.Printf("New button created: %s", id)
}

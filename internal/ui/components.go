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

type ButtonConfig struct {
	Sizing          clay.Sizing
	Padding         clay.Padding
	BackgroundColor clay.Color
	TextColor       clay.Color
	TextSize        uint16
	CornerRadius    float32
}

// Funções helper para criar configurações comuns

func DefaultContainerConfig() ContainerConfig {
	return ContainerConfig{
		Sizing: clay.Sizing{
			Width:  clay.SizingGrow(0),
			Height: clay.SizingGrow(0),
		},
		Padding:         clay.PaddingAll(0),
		ChildGap:        0,
		LayoutDirection: clay.TOP_TO_BOTTOM,
		BackgroundColor: clay.Color{R: 0, G: 0, B: 0, A: 0},
	}
}

func DefaultButtonConfig() ButtonConfig {
	return ButtonConfig{
		Sizing: clay.Sizing{
			Width:  clay.SizingFixed(220),
			Height: clay.SizingFixed(45),
		},
		Padding:         clay.Padding{Left: 20, Right: 20, Top: 12, Bottom: 12},
		BackgroundColor: clay.Color{R: 70, G: 130, B: 220, A: 200}, // Azul moderno com transparência
		TextColor:       clay.Color{R: 255, G: 255, B: 255, A: 255},
		TextSize:        16,
		CornerRadius:    12, // Bordas mais arredondadas
	}
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

// CreateButton cria um botão interativo
func (cls *ClayLayoutSystem) CreateButton(id string, text string, config ButtonConfig, onClick func()) {
	if !cls.enabled || !cls.isActive {
		log.Printf("Clay not enabled or not active, skipping CreateButton: %s", id)
		return
	}

	log.Printf("Creating button: %s", id)
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
		BackgroundColor: config.BackgroundColor,
	}, func() {
		// Texto do botão centralizado
		cls.CreateText(text, TextConfig{
			FontSize:  config.TextSize,
			TextColor: config.TextColor,
		})
	})
	log.Printf("Button created successfully: %s", id)
}

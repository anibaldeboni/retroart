package theme

import (
	"fmt"
	"log"

	"retroart-sdl2/internal/renderer"

	"github.com/TotallyGamerJet/clay"
	"github.com/veandco/go-sdl2/ttf"
)

// ColorPalette define a paleta de cores do design system
type ColorPalette struct {
	// Cores primárias
	Primary       clay.Color
	PrimaryHover  clay.Color
	PrimaryActive clay.Color

	// Cores secundárias
	Secondary       clay.Color
	SecondaryHover  clay.Color
	SecondaryActive clay.Color

	// Cores semânticas
	Success      clay.Color
	SuccessHover clay.Color
	Warning      clay.Color
	WarningHover clay.Color
	Danger       clay.Color
	DangerHover  clay.Color
	Info         clay.Color
	InfoHover    clay.Color

	// Cores neutras
	Background       clay.Color
	Surface          clay.Color
	SurfaceSecondary clay.Color
	SurfaceTertiary  clay.Color
	Border           clay.Color

	// Cores de texto
	TextPrimary     clay.Color
	TextSecondary   clay.Color
	TextMuted       clay.Color
	TextOnPrimary   clay.Color
	TextOnSecondary clay.Color
	TextOnDanger    clay.Color
	TextOnSuccess   clay.Color
	TextPlaceholder clay.Color

	// Cores de input/interação
	InputBackground        clay.Color
	InputBackgroundFocused clay.Color
	InputBorder            clay.Color
	InputBorderFocused     clay.Color

	// Cores específicas para componentes
	CheckboxBackground clay.Color
	CheckboxSelected   clay.Color
	KeyboardBackground clay.Color
	Overlay            clay.Color
}

// Typography define os tamanhos de tipografia
type Typography struct {
	XSmall uint16
	Small  uint16
	Base   uint16
	Large  uint16
	XLarge uint16
}

// Spacing define os espaçamentos padronizados
type Spacing struct {
	XS uint16
	SM uint16
	MD uint16
	LG uint16
	XL uint16
}

// BorderRadius define os raios de borda padronizados
type BorderRadius struct {
	Small  float32
	Medium float32
	Large  float32
	XLarge float32
}

type BorderWidth struct {
	XSmall uint16
	Small  uint16
	Medium uint16
	Large  uint16
	XLarge uint16
}

// Elevation define configurações de elevação/sombra
type Elevation struct {
	Low    uint8
	Medium uint8
	High   uint8
}

// DesignSystem contém todas as configurações de design
type DesignSystem struct {
	Colors     ColorPalette
	Typography Typography
	Spacing    Spacing
	Border     Border
	Elevation  Elevation
}

type Border struct {
	Radius BorderRadius
	Width  BorderWidth
	Color  clay.Color
}

// DefaultDesignSystem retorna o design system padrão baseado no tema atual da aplicação
func DefaultDesignSystem() DesignSystem {
	return DesignSystem{
		Colors: ColorPalette{
			// Cores primárias (azul macOS)
			Primary:       clay.Color{R: 0, G: 122, B: 255, A: 255},
			PrimaryHover:  clay.Color{R: 10, G: 132, B: 255, A: 255},
			PrimaryActive: clay.Color{R: 0, G: 112, B: 245, A: 255},

			// Cores secundárias (roxo/índigo macOS)
			Secondary:       clay.Color{R: 88, G: 86, B: 214, A: 255},
			SecondaryHover:  clay.Color{R: 98, G: 96, B: 224, A: 255},
			SecondaryActive: clay.Color{R: 78, G: 76, B: 204, A: 255},

			// Cores semânticas (baseadas no macOS)
			Success:      clay.Color{R: 52, G: 199, B: 89, A: 255},
			SuccessHover: clay.Color{R: 62, G: 209, B: 99, A: 255},
			Warning:      clay.Color{R: 255, G: 159, B: 10, A: 255},
			WarningHover: clay.Color{R: 255, G: 169, B: 20, A: 255},
			Danger:       clay.Color{R: 255, G: 69, B: 58, A: 255},
			DangerHover:  clay.Color{R: 255, G: 79, B: 68, A: 255},
			Info:         clay.Color{R: 90, G: 200, B: 250, A: 255},
			InfoHover:    clay.Color{R: 100, G: 210, B: 255, A: 255},

			// Cores neutras (macOS dark style)
			Background:       clay.Color{R: 30, G: 30, B: 30, A: 255},
			Surface:          clay.Color{R: 44, G: 44, B: 46, A: 255},
			SurfaceSecondary: clay.Color{R: 58, G: 58, B: 60, A: 255},
			SurfaceTertiary:  clay.Color{R: 72, G: 72, B: 74, A: 255},
			Border:           clay.Color{R: 99, G: 99, B: 102, A: 255},

			// Cores de texto (macOS dark theme)
			TextPrimary:     clay.Color{R: 255, G: 255, B: 255, A: 255},
			TextSecondary:   clay.Color{R: 235, G: 235, B: 245, A: 153},
			TextMuted:       clay.Color{R: 235, G: 235, B: 245, A: 102},
			TextOnPrimary:   clay.Color{R: 255, G: 255, B: 255, A: 255},
			TextOnSecondary: clay.Color{R: 255, G: 255, B: 255, A: 255},
			TextOnDanger:    clay.Color{R: 255, G: 255, B: 255, A: 255},
			TextOnSuccess:   clay.Color{R: 255, G: 255, B: 255, A: 255},
			TextPlaceholder: clay.Color{R: 235, G: 235, B: 245, A: 76},

			// Cores de input/interação
			InputBackground:        clay.Color{R: 28, G: 28, B: 30, A: 255},
			InputBackgroundFocused: clay.Color{R: 44, G: 44, B: 46, A: 255},
			InputBorder:            clay.Color{R: 72, G: 72, B: 74, A: 255},
			InputBorderFocused:     clay.Color{R: 0, G: 122, B: 255, A: 255},

			// Cores específicas para components
			CheckboxBackground: clay.Color{R: 28, G: 28, B: 30, A: 255},
			CheckboxSelected:   clay.Color{R: 52, G: 199, B: 89, A: 255},
			KeyboardBackground: clay.Color{R: 44, G: 44, B: 46, A: 220},
			Overlay:            clay.Color{R: 0, G: 0, B: 0, A: 128},
		},
		Typography: Typography{
			XSmall: 12,
			Small:  14,
			Base:   16,
			Large:  18,
			XLarge: 24,
		},
		Spacing: Spacing{
			XS: 4,
			SM: 8,
			MD: 12,
			LG: 20,
			XL: 32,
		},
		Border: Border{
			Radius: BorderRadius{
				Small:  4,
				Medium: 8,
				Large:  12,
				XLarge: 16,
				// Small:  0,
				// Medium: 0,
				// Large:  0,
				// XLarge: 0,
			},
			Width: BorderWidth{
				// XSmall: 1,
				// Small:  2,
				// Medium: 3,
				// Large:  4,
				// XLarge: 5,
				XSmall: 0,
				Small:  0,
				Medium: 0,
				Large:  0,
				XLarge: 0,
			},
		},

		Elevation: Elevation{
			Low:    50,
			Medium: 100,
			High:   200,
		},
	}
}

// GetFontIdForSize maps font sizes to Clay FontId indices based on typography system
// Clay expects FontId to be the index in the fonts array, not the actual size
func GetFontIdForSize(fontSize uint16) uint16 {
	sizes := GetAllTypographySizes()

	// Find exact match first
	for i, size := range sizes {
		if size == fontSize {
			return uint16(i)
		}
	}

	// If no exact match, find closest size
	var closestIndex int
	var minDiff uint16 = 1000

	for i, size := range sizes {
		var diff uint16
		if size > fontSize {
			diff = size - fontSize
		} else {
			diff = fontSize - size
		}

		if diff < minDiff {
			minDiff = diff
			closestIndex = i
		}
	}

	// Safety check: ensure we don't return an index that could be out of bounds
	if closestIndex >= len(sizes) {
		log.Printf("Warning: fontId %d is out of bounds for %d available fonts, using 0", closestIndex, len(sizes))
		return 0
	}

	log.Printf("Font size %d mapped to fontId %d (actual size %d)", fontSize, closestIndex, sizes[closestIndex])
	return uint16(closestIndex)
}

// FontSystem gerencia o carregamento de fontes baseado na tipografia
type FontSystem struct {
	fonts []renderer.Font
}

// NewFontSystem cria um novo sistema de fontes
func NewFontSystem() *FontSystem {
	return &FontSystem{
		fonts: make([]renderer.Font, 0),
	}
}

// GetAllTypographySizes retorna todos os tamanhos de fonte definidos na tipografia
func GetAllTypographySizes() []uint16 {
	ds := DefaultDesignSystem()
	return []uint16{
		ds.Typography.XSmall,
		ds.Typography.Small,
		ds.Typography.Base,
		ds.Typography.Large,
		ds.Typography.XLarge,
	}
}

// InitializeFonts inicializa o sistema de fontes carregando todos os tamanhos da tipografia
func (fs *FontSystem) InitializeFonts() error {
	// Use actual typography sizes from design system
	typographySizes := GetAllTypographySizes()
	typographyInts := make([]int, len(typographySizes))
	for i, size := range typographySizes {
		typographyInts[i] = int(size)
	}

	fs.fonts = make([]renderer.Font, len(typographyInts))

	for i, size := range typographyInts {
		font, err := fs.loadFontWithSize(size)
		if err != nil {
			return fmt.Errorf("failed to load font size %d: %v", size, err)
		}

		clayFont := renderer.Font{FontId: uint32(i), Font: font}
		fs.fonts[i] = clayFont
		log.Printf("Successfully loaded font size %d at index %d", size, i)
	}

	log.Printf("Font system initialized with %d fonts: %v", len(typographyInts), typographyInts)
	return nil
}

// loadFontWithSize carrega uma fonte com tamanho específico
func (fs *FontSystem) loadFontWithSize(size int) (*ttf.Font, error) {
	// Lista de possíveis caminhos de fontes no sistema
	fontPaths := []string{
		"assets/DejaVuSansCondensed.ttf",
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
		"/usr/share/fonts/TTF/DejaVuSans.ttf",
		"/System/Library/Fonts/Helvetica.ttc",
		"/usr/share/fonts/liberation/LiberationSans-Regular.ttf",
	}

	for _, fontPath := range fontPaths {
		font, err := ttf.OpenFont(fontPath, size)
		if err == nil {
			log.Printf("Successfully loaded font from: %s (size %d)", fontPath, size)
			return font, nil
		}
	}

	return nil, fmt.Errorf("could not load any font for size %d", size)
}

// GetFonts returns a pointer to the internal clayFonts slice for Clay's MeasureText function
// This ensures Clay gets a stable pointer that won't be garbage collected
func (fs *FontSystem) GetFonts() *[]renderer.Font {
	if len(fs.fonts) == 0 {
		log.Printf("Warning: GetClayFonts called but no fonts initialized")
		return nil
	}

	log.Printf("GetClayFonts returning %d fonts", len(fs.fonts))
	return &fs.fonts
}

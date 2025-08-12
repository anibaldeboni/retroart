package theme

import "github.com/TotallyGamerJet/clay"

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
	Border           clay.Color

	// Cores de texto
	TextPrimary   clay.Color
	TextSecondary clay.Color
	TextMuted     clay.Color
	TextOnPrimary clay.Color
	TextOnDanger  clay.Color
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

// Elevation define configurações de elevação/sombra
type Elevation struct {
	Low    uint8
	Medium uint8
	High   uint8
}

// DesignSystem contém todas as configurações de design
type DesignSystem struct {
	Colors       ColorPalette
	Typography   Typography
	Spacing      Spacing
	BorderRadius BorderRadius
	Elevation    Elevation
}

// DefaultDesignSystem retorna o design system padrão baseado no tema atual da aplicação
func DefaultDesignSystem() DesignSystem {
	return DesignSystem{
		Colors: ColorPalette{
			// Cores primárias (azul)
			Primary:       clay.Color{R: 50, G: 120, B: 200, A: 200},
			PrimaryHover:  clay.Color{R: 30, G: 150, B: 255, A: 255},
			PrimaryActive: clay.Color{R: 20, G: 100, B: 180, A: 255},

			// Cores secundárias (roxo)
			Secondary:       clay.Color{R: 120, G: 60, B: 200, A: 180},
			SecondaryHover:  clay.Color{R: 150, G: 80, B: 255, A: 255},
			SecondaryActive: clay.Color{R: 100, G: 40, B: 180, A: 255},

			// Cores semânticas
			Success:      clay.Color{R: 40, G: 120, B: 80, A: 255},
			SuccessHover: clay.Color{R: 30, G: 180, B: 120, A: 255},
			Warning:      clay.Color{R: 255, G: 193, B: 7, A: 255},
			WarningHover: clay.Color{R: 255, G: 207, B: 50, A: 255},
			Danger:       clay.Color{R: 200, G: 60, B: 60, A: 180},
			DangerHover:  clay.Color{R: 255, G: 80, B: 80, A: 255},
			Info:         clay.Color{R: 60, G: 120, B: 200, A: 255},
			InfoHover:    clay.Color{R: 80, G: 140, B: 220, A: 255},

			// Cores neutras
			Background:       clay.Color{R: 40, G: 42, B: 54, A: 255},
			Surface:          clay.Color{R: 60, G: 63, B: 75, A: 180},
			SurfaceSecondary: clay.Color{R: 25, G: 30, B: 40, A: 240},
			Border:           clay.Color{R: 60, G: 70, B: 85, A: 255},

			// Cores de texto
			TextPrimary:   clay.Color{R: 255, G: 255, B: 255, A: 255},
			TextSecondary: clay.Color{R: 230, G: 230, B: 230, A: 255},
			TextMuted:     clay.Color{R: 200, G: 200, B: 200, A: 255},
			TextOnPrimary: clay.Color{R: 220, G: 230, B: 255, A: 255},
			TextOnDanger:  clay.Color{R: 255, G: 200, B: 200, A: 255},
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
		BorderRadius: BorderRadius{
			Small:  4,
			Medium: 8,
			Large:  12,
			XLarge: 16,
		},
		Elevation: Elevation{
			Low:    50,
			Medium: 100,
			High:   200,
		},
	}
}

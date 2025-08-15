package theme

import "github.com/TotallyGamerJet/clay"

// ContainerStyle define estilos para containers
type ContainerStyle struct {
	Image           clay.ImageElementConfig
	BackgroundColor clay.Color
	Padding         clay.Padding
	CornerRadius    float32
	Border          clay.BorderElementConfig
}

// GetMainContainerStyle retorna o estilo para o container principal
func (ds DesignSystem) GetMainContainerStyle() ContainerStyle {
	return ContainerStyle{
		BackgroundColor: ds.Colors.Background,
		Padding:         clay.Padding{Left: 0, Right: 0, Top: 0, Bottom: 0},
		CornerRadius:    0,
	}
}

// GetContentContainerStyle retorna o estilo para containers de conteúdo
func (ds DesignSystem) GetContentContainerStyle() ContainerStyle {
	return ContainerStyle{
		BackgroundColor: ds.Colors.Surface,
		Padding:         clay.Padding{Left: ds.Spacing.LG, Right: ds.Spacing.LG, Top: ds.Spacing.LG, Bottom: ds.Spacing.LG},
		CornerRadius:    ds.Border.Radius.Large,
		Border: clay.BorderElementConfig{
			Width: clay.BorderWidth{Left: ds.Border.Width.Small, Right: ds.Border.Width.Small, Top: ds.Border.Width.Small, Bottom: ds.Border.Width.Small},
			Color: ds.Colors.Border,
		},
	}
}

package theme

import "github.com/TotallyGamerJet/clay"

// CheckboxListStyle contém configurações para checkbox lists
type CheckboxListStyle struct {
	FontSize        uint16
	Padding         clay.Padding
	ChildGap        uint16
	BackgroundColor clay.Color
	ScrollOffset    int
	Checkbox        Checkbox
	CornerRadius    float32

	// Cores para diferentes estados dos itens
	ItemNormalBg     clay.Color
	ItemSelectedBg   clay.Color
	ItemFocusedBg    clay.Color
	ItemNormalText   clay.Color
	ItemSelectedText clay.Color
	ItemFocusedText  clay.Color
}

type CheckboxScrollIndicator struct {
	Size       uint16
	UpSymbol   string
	DownSymbol string
}

// GetCheckboxListStyle retorna a configuração de estilo para checkbox lists
func (ds DesignSystem) GetCheckboxListStyle() CheckboxListStyle {
	return CheckboxListStyle{
		Padding:         clay.Padding{Left: ds.Spacing.SM, Right: ds.Spacing.SM, Top: ds.Spacing.MD, Bottom: ds.Spacing.MD},
		ChildGap:        ds.Spacing.SM,
		BackgroundColor: ds.Colors.SurfaceSecondary,
		ScrollOffset:    0,
		FontSize:        ds.Typography.Large,
		Checkbox: Checkbox{
			Size:         22,
			CornerRadius: ds.Border.Radius.Small,
			Background:   ds.Colors.CheckboxBackground,
			Color: CheckboxColor{
				Normal:   ds.Colors.SurfaceTertiary,
				Selected: ds.Colors.CheckboxSelected,
				Mark:     ds.Colors.TextPrimary,
			},
			Mark: CheckboxMark{
				Symbol: "◣",
				Size:   ds.Typography.Base,
			},
			ScrollIndicator: CheckboxScrollIndicator{
				Size:       ds.Typography.Base,
				UpSymbol:   "▲",
				DownSymbol: "▼",
			},
		},
		CornerRadius: ds.Border.Radius.Large,

		// Estados dos itens
		ItemNormalBg:     clay.Color{R: 0, G: 0, B: 0, A: 0}, // Transparente
		ItemSelectedBg:   ds.Colors.Success,
		ItemFocusedBg:    ds.Colors.Info,
		ItemNormalText:   ds.Colors.TextSecondary,
		ItemSelectedText: ds.Colors.TextOnSuccess,
		ItemFocusedText:  ds.Colors.TextPrimary,
	}
}

package theme

import "github.com/TotallyGamerJet/clay"

type Checkbox struct {
	Size            float32
	Border          clay.BorderElementConfig
	CornerRadius    float32
	Background      clay.Color
	Mark            CheckboxMark
	ScrollIndicator CheckboxScrollIndicator
	Color           CheckboxColor
}

type CheckboxMark struct {
	Symbol string
	Size   uint16
}

type CheckboxColor struct {
	Normal   clay.Color
	Selected clay.Color
	Mark     clay.Color
}

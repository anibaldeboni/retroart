package ui

import (
	"github.com/TotallyGamerJet/clay"
	"github.com/veandco/go-sdl2/sdl"
)

// renderRectangle renderiza um retângulo com suporte a cantos arredondados
func (cls *ClayLayoutSystem) renderRectangle(renderer *sdl.Renderer, command *clay.RenderCommand) error {
	config := &command.RenderData.Rectangle
	boundingBox := command.BoundingBox

	// Definir cor de fundo
	renderer.SetDrawColor(
		uint8(config.BackgroundColor.R),
		uint8(config.BackgroundColor.G),
		uint8(config.BackgroundColor.B),
		uint8(config.BackgroundColor.A),
	)

	// Verificar se há cantos arredondados
	cornerRadius := config.CornerRadius
	hasRoundedCorners := cornerRadius.TopLeft > 0 || cornerRadius.TopRight > 0 ||
		cornerRadius.BottomLeft > 0 || cornerRadius.BottomRight > 0

	if hasRoundedCorners {
		// Renderizar retângulo com cantos arredondados
		return cls.renderRoundedRectangle(renderer, boundingBox, cornerRadius, config.BackgroundColor)
	} else {
		// Renderizar retângulo normal
		rect := sdl.Rect{
			X: int32(boundingBox.X),
			Y: int32(boundingBox.Y),
			W: int32(boundingBox.Width),
			H: int32(boundingBox.Height),
		}
		return renderer.FillRect(&rect)
	}
}

// renderRoundedRectangle renderiza um retângulo com cantos arredondados
func (cls *ClayLayoutSystem) renderRoundedRectangle(renderer *sdl.Renderer, boundingBox clay.BoundingBox, cornerRadius clay.CornerRadius, color clay.Color) error {
	x := int32(boundingBox.X)
	y := int32(boundingBox.Y)
	w := int32(boundingBox.Width)
	h := int32(boundingBox.Height)

	// Usar o maior raio de canto para simplificar
	radius := int32(cornerRadius.TopLeft)
	if cornerRadius.TopRight > cornerRadius.TopLeft {
		radius = int32(cornerRadius.TopRight)
	}
	if cornerRadius.BottomLeft > float32(radius) {
		radius = int32(cornerRadius.BottomLeft)
	}
	if cornerRadius.BottomRight > float32(radius) {
		radius = int32(cornerRadius.BottomRight)
	}

	// Limitar o raio ao tamanho do retângulo
	maxRadius := w / 2
	if h/2 < maxRadius {
		maxRadius = h / 2
	}
	if radius > maxRadius {
		radius = maxRadius
	}

	// Renderizar retângulo central
	if radius > 0 {
		centerRect := sdl.Rect{
			X: x + radius,
			Y: y,
			W: w - 2*radius,
			H: h,
		}
		renderer.FillRect(&centerRect)

		leftRect := sdl.Rect{
			X: x,
			Y: y + radius,
			W: radius,
			H: h - 2*radius,
		}
		renderer.FillRect(&leftRect)

		rightRect := sdl.Rect{
			X: x + w - radius,
			Y: y + radius,
			W: radius,
			H: h - 2*radius,
		}
		renderer.FillRect(&rightRect)

		// Renderizar cantos arredondados
		cls.renderFilledCircle(renderer, x+radius, y+radius, radius)
		cls.renderFilledCircle(renderer, x+w-radius-1, y+radius, radius)
		cls.renderFilledCircle(renderer, x+radius, y+h-radius-1, radius)
		cls.renderFilledCircle(renderer, x+w-radius-1, y+h-radius-1, radius)
	} else {
		rect := sdl.Rect{X: x, Y: y, W: w, H: h}
		renderer.FillRect(&rect)
	}

	return nil
}

// renderFilledCircle renderiza um círculo preenchido
func (cls *ClayLayoutSystem) renderFilledCircle(renderer *sdl.Renderer, centerX, centerY, radius int32) {
	r, g, b, a, _ := renderer.GetDrawColor()

	subPixelRadius := float32(radius) + 0.5

	for y := -radius - 1; y <= radius+1; y++ {
		for x := -radius - 1; x <= radius+1; x++ {
			coverage := cls.calculatePixelCoverage(float32(x), float32(y), subPixelRadius)

			if coverage > 0 {
				alpha := uint8(float32(a) * coverage)
				renderer.SetDrawColor(r, g, b, alpha)
				renderer.DrawPoint(centerX+x, centerY+y)
			}
		}
	}

	renderer.SetDrawColor(r, g, b, a)
}

// calculatePixelCoverage calcula cobertura de pixel
func (cls *ClayLayoutSystem) calculatePixelCoverage(x, y, radius float32) float32 {
	subPixels := 4
	coveredSubPixels := 0
	subPixelSize := 1.0 / float32(subPixels)

	for sy := range subPixels {
		for sx := range subPixels {
			subX := x + (float32(sx)+0.5)*subPixelSize - 0.5
			subY := y + (float32(sy)+0.5)*subPixelSize - 0.5

			distance := subX*subX + subY*subY
			if distance <= radius*radius {
				coveredSubPixels++
			}
		}
	}

	return float32(coveredSubPixels) / float32(subPixels*subPixels)
}

// renderText renderiza texto
func (cls *ClayLayoutSystem) renderText(renderer *sdl.Renderer, command *clay.RenderCommand) error {
	if cls.font == nil {
		return nil
	}

	config := &command.RenderData.Text
	boundingBox := command.BoundingBox

	text := config.StringContents.String()
	if text == "" {
		return nil
	}

	color := sdl.Color{
		R: uint8(config.TextColor.R),
		G: uint8(config.TextColor.G),
		B: uint8(config.TextColor.B),
		A: uint8(config.TextColor.A),
	}

	surface, err := cls.font.RenderUTF8Blended(text, color)
	if err != nil {
		return err
	}
	defer surface.Free()

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return err
	}
	defer texture.Destroy()

	// Calcular posição centralizada
	textWidth := surface.W
	textHeight := surface.H
	containerWidth := int32(boundingBox.Width)
	containerHeight := int32(boundingBox.Height)

	offsetX := (containerWidth - textWidth) / 2
	offsetY := (containerHeight - textHeight) / 2

	if offsetX < 0 {
		offsetX = 0
	}
	if offsetY < 0 {
		offsetY = 0
	}

	destRect := sdl.Rect{
		X: int32(boundingBox.X) + offsetX,
		Y: int32(boundingBox.Y) + offsetY,
		W: textWidth,
		H: textHeight,
	}

	return renderer.Copy(texture, nil, &destRect)
}

// renderBorder renderiza borda
func (cls *ClayLayoutSystem) renderBorder(renderer *sdl.Renderer, command *clay.RenderCommand) error {
	config := &command.RenderData.Border
	boundingBox := command.BoundingBox

	renderer.SetDrawColor(
		uint8(config.Color.R),
		uint8(config.Color.G),
		uint8(config.Color.B),
		uint8(config.Color.A),
	)

	// Renderizar bordas
	if config.Width.Top > 0 {
		rect := sdl.Rect{
			X: int32(boundingBox.X),
			Y: int32(boundingBox.Y),
			W: int32(boundingBox.Width),
			H: int32(config.Width.Top),
		}
		renderer.FillRect(&rect)
	}

	if config.Width.Bottom > 0 {
		rect := sdl.Rect{
			X: int32(boundingBox.X),
			Y: int32(boundingBox.Y + boundingBox.Height - float32(config.Width.Bottom)),
			W: int32(boundingBox.Width),
			H: int32(config.Width.Bottom),
		}
		renderer.FillRect(&rect)
	}

	if config.Width.Left > 0 {
		rect := sdl.Rect{
			X: int32(boundingBox.X),
			Y: int32(boundingBox.Y),
			W: int32(config.Width.Left),
			H: int32(boundingBox.Height),
		}
		renderer.FillRect(&rect)
	}

	if config.Width.Right > 0 {
		rect := sdl.Rect{
			X: int32(boundingBox.X + boundingBox.Width - float32(config.Width.Right)),
			Y: int32(boundingBox.Y),
			W: int32(config.Width.Right),
			H: int32(boundingBox.Height),
		}
		renderer.FillRect(&rect)
	}

	return nil
}

// renderScissorStart inicia recorte
func (cls *ClayLayoutSystem) renderScissorStart(renderer *sdl.Renderer, command *clay.RenderCommand) error {
	boundingBox := command.BoundingBox

	clipRect := sdl.Rect{
		X: int32(boundingBox.X),
		Y: int32(boundingBox.Y),
		W: int32(boundingBox.Width),
		H: int32(boundingBox.Height),
	}

	return renderer.SetClipRect(&clipRect)
}

// renderScissorEnd termina recorte
func (cls *ClayLayoutSystem) renderScissorEnd(renderer *sdl.Renderer) error {
	return renderer.SetClipRect(nil)
}

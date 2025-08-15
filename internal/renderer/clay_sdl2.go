package renderer

import (
	"fmt"
	"log/slog"
	"strings"
	"unsafe"

	"github.com/TotallyGamerJet/clay"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Font struct {
	FontId uint32
	Font   *ttf.Font
}

func MeasureText(text clay.StringSlice, config *clay.TextElementConfig, userData unsafe.Pointer) clay.Dimensions {
	fonts := *(*[]Font)(userData)
	font := fonts[config.FontId].Font
	chars := strings.Clone(text.String())
	width, height, err := font.SizeUTF8(chars)
	if err != nil {
		panic(fmt.Errorf("sdl2: failed to measure text: %w", err))
	}
	return clay.Dimensions{
		Width:  float32(width),
		Height: float32(height),
	}
}

// Custom SDL2 renderer that avoids RenderGeometry for TrimUI compatibility
func ClayRender(renderer *sdl.Renderer, renderCommands clay.RenderCommandArray, fonts []Font) error {
	for renderCommand := range renderCommands.Iter() {
		boundingBox := renderCommand.BoundingBox
		switch renderCommand.CommandType {
		case clay.RENDER_COMMAND_TYPE_RECTANGLE:
			config := &renderCommand.RenderData.Rectangle
			color := config.BackgroundColor
			if err := renderer.SetDrawColor(uint8(color.R), uint8(color.G), uint8(color.B), uint8(color.A)); err != nil {
				return err
			}
			rect := sdl.FRect{
				X: boundingBox.X,
				Y: boundingBox.Y,
				W: boundingBox.Width,
				H: boundingBox.Height,
			}
			if config.CornerRadius.TopLeft > 0 {
				if err := renderFillRoundedRectPrimitive(renderer, rect, config.CornerRadius.TopLeft, color); err != nil {
					return err
				}
			} else {
				if err := renderer.FillRectF(&rect); err != nil {
					return err
				}
			}
		case clay.RENDER_COMMAND_TYPE_TEXT:
			config := &renderCommand.RenderData.Text
			cloned := strings.Clone(config.StringContents.String())
			font := fonts[config.FontId].Font
			surface, err := font.RenderUTF8Blended(cloned, sdl.Color{
				R: uint8(config.TextColor.R),
				G: uint8(config.TextColor.G),
				B: uint8(config.TextColor.B),
				A: uint8(config.TextColor.A),
			})
			if err != nil {
				return err
			}
			texture, err := renderer.CreateTextureFromSurface(surface)
			if err != nil {
				return err
			}
			destination := sdl.Rect{
				X: int32(boundingBox.X),
				Y: int32(boundingBox.Y),
				W: int32(boundingBox.Width),
				H: int32(boundingBox.Height),
			}
			if err := renderer.Copy(texture, nil, &destination); err != nil {
				return err
			}
			surface.Free()
			if err := texture.Destroy(); err != nil {
				return err
			}
		case clay.RENDER_COMMAND_TYPE_IMAGE:
			config := &renderCommand.RenderData.Image
			texture, err := renderer.CreateTextureFromSurface((*sdl.Surface)(config.ImageData.(unsafe.Pointer)))
			if err != nil {
				return err
			}
			destination := sdl.Rect{
				X: int32(boundingBox.X),
				Y: int32(boundingBox.Y),
				W: int32(boundingBox.Width),
				H: int32(boundingBox.Height),
			}
			if err := renderer.Copy(texture, nil, &destination); err != nil {
				return err
			}
			if err := texture.Destroy(); err != nil {
				return err
			}
		case clay.RENDER_COMMAND_TYPE_BORDER:
			config := &renderCommand.RenderData.Border
			color := config.Color
			if err := renderer.SetDrawColor(uint8(color.R), uint8(color.G), uint8(color.B), uint8(color.A)); err != nil {
				return err
			}
			if err := renderBorderPrimitive(renderer, boundingBox, config); err != nil {
				return err
			}
		case clay.RENDER_COMMAND_TYPE_SCISSOR_START:
			rect := sdl.Rect{
				X: int32(boundingBox.X),
				Y: int32(boundingBox.Y),
				W: int32(boundingBox.Width),
				H: int32(boundingBox.Height),
			}
			if err := renderer.SetClipRect(&rect); err != nil {
				return err
			}
		case clay.RENDER_COMMAND_TYPE_SCISSOR_END:
			renderer.SetClipRect(nil)
		default:
			slog.Warn("Unknown render command type", "type", renderCommand.CommandType)
		}
	}
	return nil
}

// Custom rounded rectangle implementation using SDL2 primitives only
func renderFillRoundedRectPrimitive(renderer *sdl.Renderer, rect sdl.FRect, cornerRadius float32, color clay.Color) error {
	// Convert to int32 for SDL calls
	x := int32(rect.X)
	y := int32(rect.Y)
	w := int32(rect.W)
	h := int32(rect.H)
	radius := int32(cornerRadius)

	// Limit radius to rectangle size
	maxRadius := min(h/2, w/2)
	if radius > maxRadius {
		radius = maxRadius
	}

	if radius <= 0 {
		return renderer.FillRectF(&rect)
	}

	// Set drawing color
	if err := renderer.SetDrawColor(uint8(color.R), uint8(color.G), uint8(color.B), uint8(color.A)); err != nil {
		return err
	}

	// Draw the main body rectangles using the original simple approach
	// Center rectangle (full width, full height)
	centerRect := sdl.Rect{
		X: x + radius,
		Y: y,
		W: w - 2*radius,
		H: h,
	}
	if err := renderer.FillRect(&centerRect); err != nil {
		return err
	}

	// Left rectangle (covers left area minus corners)
	leftRect := sdl.Rect{
		X: x,
		Y: y + radius,
		W: radius,
		H: h - 2*radius,
	}
	if err := renderer.FillRect(&leftRect); err != nil {
		return err
	}

	// Right rectangle (covers right area minus corners)
	rightRect := sdl.Rect{
		X: x + w - radius,
		Y: y + radius,
		W: radius,
		H: h - 2*radius,
	}
	if err := renderer.FillRect(&rightRect); err != nil {
		return err
	}

	// Draw full filled circles at corners (like the original implementation)
	if err := renderFilledCircle(renderer, x+radius, y+radius, radius); err != nil {
		return err
	}
	if err := renderFilledCircle(renderer, x+w-radius-1, y+radius, radius); err != nil {
		return err
	}
	if err := renderFilledCircle(renderer, x+radius, y+h-radius-1, radius); err != nil {
		return err
	}
	if err := renderFilledCircle(renderer, x+w-radius-1, y+h-radius-1, radius); err != nil {
		return err
	}

	return nil
}

// renderFilledCircle renders a filled circle with anti-aliasing using sub-pixel sampling
// This is the original simple implementation that worked well
func renderFilledCircle(renderer *sdl.Renderer, centerX, centerY, radius int32) error {
	// Get current drawing color
	r, g, b, a, err := renderer.GetDrawColor()
	if err != nil {
		return err
	}

	subPixelRadius := float32(radius) + 0.5

	for y := -radius - 1; y <= radius+1; y++ {
		for x := -radius - 1; x <= radius+1; x++ {
			coverage := calculatePixelCoverage(float32(x), float32(y), subPixelRadius)

			if coverage > 0 {
				alpha := uint8(float32(a) * coverage)
				if err := renderer.SetDrawColor(r, g, b, alpha); err != nil {
					return err
				}
				if err := renderer.DrawPoint(centerX+x, centerY+y); err != nil {
					return err
				}
			}
		}
	}

	// Restore original color
	return renderer.SetDrawColor(r, g, b, a)
}

// calculatePixelCoverage calculates pixel coverage for anti-aliasing
func calculatePixelCoverage(x, y, radius float32) float32 {
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

// renderBorderPrimitive renders borders using SDL2 primitives only, TrimUI compatible
func renderBorderPrimitive(renderer *sdl.Renderer, boundingBox clay.BoundingBox, config *clay.BorderRenderData) error {
	if boundingBox.Width <= 0 || boundingBox.Height <= 0 {
		return nil
	}

	maxRadius := min(boundingBox.Width, boundingBox.Height) / 2.0

	// Convert border configuration
	leftWidth := int32(config.Width.Left)
	rightWidth := int32(config.Width.Right)
	topWidth := int32(config.Width.Top)
	bottomWidth := int32(config.Width.Bottom)

	// Convert corner radii
	topLeftRadius := min(config.CornerRadius.TopLeft, maxRadius)
	topRightRadius := min(config.CornerRadius.TopRight, maxRadius)
	bottomLeftRadius := min(config.CornerRadius.BottomLeft, maxRadius)
	bottomRightRadius := min(config.CornerRadius.BottomRight, maxRadius)

	// Convert coordinates
	x := int32(boundingBox.X)
	y := int32(boundingBox.Y)
	w := int32(boundingBox.Width)
	h := int32(boundingBox.Height)

	// Top border (horizontal rectangle avoiding corners)
	if topWidth > 0 {
		topRect := sdl.Rect{
			X: x + int32(topLeftRadius),
			Y: y,
			W: w - int32(topLeftRadius) - int32(topRightRadius),
			H: topWidth,
		}
		if topRect.W > 0 {
			if err := renderer.FillRect(&topRect); err != nil {
				return err
			}
		}
	}

	// Bottom border (horizontal rectangle avoiding corners)
	if bottomWidth > 0 {
		bottomRect := sdl.Rect{
			X: x + int32(bottomLeftRadius),
			Y: y + h - bottomWidth,
			W: w - int32(bottomLeftRadius) - int32(bottomRightRadius),
			H: bottomWidth,
		}
		if bottomRect.W > 0 {
			if err := renderer.FillRect(&bottomRect); err != nil {
				return err
			}
		}
	}

	// Left border (vertical rectangle avoiding corners)
	if leftWidth > 0 {
		leftRect := sdl.Rect{
			X: x,
			Y: y + int32(topLeftRadius),
			W: leftWidth,
			H: h - int32(topLeftRadius) - int32(bottomLeftRadius),
		}
		if leftRect.H > 0 {
			if err := renderer.FillRect(&leftRect); err != nil {
				return err
			}
		}
	}

	// Right border (vertical rectangle avoiding corners)
	if rightWidth > 0 {
		rightRect := sdl.Rect{
			X: x + w - rightWidth,
			Y: y + int32(topRightRadius),
			W: rightWidth,
			H: h - int32(topRightRadius) - int32(bottomRightRadius),
		}
		if rightRect.H > 0 {
			if err := renderer.FillRect(&rightRect); err != nil {
				return err
			}
		}
	}

	// Render corner borders using our rounded corner approach
	// Top-left corner
	if topLeftRadius > 0 && (leftWidth > 0 || topWidth > 0) {
		if err := renderCornerBorderPrimitive(renderer,
			float32(x), float32(y), topLeftRadius,
			float32(leftWidth), float32(topWidth), 0); err != nil {
			return err
		}
	}

	// Top-right corner
	if topRightRadius > 0 && (rightWidth > 0 || topWidth > 0) {
		if err := renderCornerBorderPrimitive(renderer,
			float32(x+w), float32(y), topRightRadius,
			float32(rightWidth), float32(topWidth), 1); err != nil {
			return err
		}
	}

	// Bottom-right corner
	if bottomRightRadius > 0 && (rightWidth > 0 || bottomWidth > 0) {
		if err := renderCornerBorderPrimitive(renderer,
			float32(x+w), float32(y+h), bottomRightRadius,
			float32(rightWidth), float32(bottomWidth), 2); err != nil {
			return err
		}
	}

	// Bottom-left corner
	if bottomLeftRadius > 0 && (leftWidth > 0 || bottomWidth > 0) {
		if err := renderCornerBorderPrimitive(renderer,
			float32(x), float32(y+h), bottomLeftRadius,
			float32(leftWidth), float32(bottomWidth), 3); err != nil {
			return err
		}
	}

	return nil
}

// renderCornerBorderPrimitive renders a corner border using a ring approach
func renderCornerBorderPrimitive(renderer *sdl.Renderer, cornerX, cornerY, radius, borderWidth1, borderWidth2 float32, cornerIndex int) error {
	if radius <= 0 || (borderWidth1 <= 0 && borderWidth2 <= 0) {
		return nil
	}

	// Calculate circle centers based on corner position
	var centerX, centerY float32
	switch cornerIndex {
	case 0: // Top-left
		centerX = cornerX + radius
		centerY = cornerY + radius
	case 1: // Top-right
		centerX = cornerX - radius
		centerY = cornerY + radius
	case 2: // Bottom-right
		centerX = cornerX - radius
		centerY = cornerY - radius
	case 3: // Bottom-left
		centerX = cornerX + radius
		centerY = cornerY - radius
	}

	// Use the maximum border width for the corner
	maxBorderWidth := max(borderWidth1, borderWidth2)

	// Render border ring instead of filled circles
	return renderBorderRing(renderer, int32(centerX), int32(centerY), int32(radius), int32(maxBorderWidth), cornerIndex)
}

// renderBorderRing renders a ring border (hollow circle) using pixel-by-pixel approach
// Only renders the quarter that belongs to the specified corner
func renderBorderRing(renderer *sdl.Renderer, centerX, centerY, outerRadius, borderWidth int32, cornerIndex int) error {
	if outerRadius <= 0 || borderWidth <= 0 {
		return nil
	}

	innerRadius := outerRadius - borderWidth
	if innerRadius < 0 {
		innerRadius = 0
	}

	// Get current drawing color
	r, g, b, a, err := renderer.GetDrawColor()
	if err != nil {
		return err
	}

	// Render ring by checking each pixel, but only in the correct quadrant
	for y := -outerRadius - 1; y <= outerRadius+1; y++ {
		for x := -outerRadius - 1; x <= outerRadius+1; x++ {
			// Check if this pixel belongs to the correct corner quadrant
			var inCorrectQuadrant bool
			switch cornerIndex {
			case 0: // Top-left: x <= 0, y <= 0
				inCorrectQuadrant = x <= 0 && y <= 0
			case 1: // Top-right: x >= 0, y <= 0
				inCorrectQuadrant = x >= 0 && y <= 0
			case 2: // Bottom-right: x >= 0, y >= 0
				inCorrectQuadrant = x >= 0 && y >= 0
			case 3: // Bottom-left: x <= 0, y >= 0
				inCorrectQuadrant = x <= 0 && y >= 0
			}

			if !inCorrectQuadrant {
				continue
			}

			distance := float32(x*x + y*y)
			outerRadiusSquared := float32(outerRadius * outerRadius)
			innerRadiusSquared := float32(innerRadius * innerRadius)

			// Check if pixel is within the border ring
			if distance <= outerRadiusSquared && distance >= innerRadiusSquared {
				// Calculate coverage for anti-aliasing
				outerCoverage := calculatePixelCoverage(float32(x), float32(y), float32(outerRadius)+0.5)
				innerCoverage := calculatePixelCoverage(float32(x), float32(y), float32(innerRadius)-0.5)
				
				// Border coverage is the difference
				coverage := outerCoverage - innerCoverage
				if coverage < 0 {
					coverage = 0
				}

				if coverage > 0 {
					alpha := uint8(float32(a) * coverage)
					if err := renderer.SetDrawColor(r, g, b, alpha); err != nil {
						return err
					}
					if err := renderer.DrawPoint(centerX+x, centerY+y); err != nil {
						return err
					}
				}
			}
		}
	}

	// Restore original color
	return renderer.SetDrawColor(r, g, b, a)
}

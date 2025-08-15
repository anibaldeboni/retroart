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
			// Skip border rendering to avoid RenderGeometry issues
			// Borders are disabled in our theme anyway
			slog.Debug("Skipping border rendering for TrimUI compatibility")
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

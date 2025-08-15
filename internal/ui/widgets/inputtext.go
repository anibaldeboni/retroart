package widgets

import (
	"log"
	"retroart-sdl2/internal/input"
	"retroart-sdl2/internal/theme"

	"github.com/TotallyGamerJet/clay"
)

type InputText struct {
	ID          string
	Text        string
	Placeholder string
	MaxLength   int
	Width       clay.SizingAxis
	Height      clay.SizingAxis
	CursorPos   int
	Config      theme.InputTextStyle
	OnChange    func(text string)
	OnSubmit    func(text string)
	focused     bool
	enabled     bool
	showCursor  bool
	keyboard    *VirtualKeyboard
}

// NewInputText cria um novo campo de entrada de texto
func NewInputText(id, placeholder string, maxLength int, width, height clay.SizingAxis, onChange, onSubmit func(string)) *InputText {
	inputText := &InputText{
		ID:          id,
		Text:        "",
		Placeholder: placeholder,
		Width:       width,
		Height:      height,
		MaxLength:   maxLength,
		CursorPos:   0,
		Config:      theme.GetInputTextStyle(),
		OnChange:    onChange,
		OnSubmit:    onSubmit,
		enabled:     true,
		showCursor:  true,
	}

	// Criar teclado virtual associado
	inputText.keyboard = NewVirtualKeyboard(id+"_keyboard", inputText.onKeyboardInput)

	return inputText
}

// Interface Focusable implementation
func (it *InputText) GetID() string {
	return it.ID
}

func (it *InputText) IsFocused() bool {
	return it.focused
}

func (it *InputText) OnFocusChanged(focused bool) {
	it.focused = focused
	if !focused {
		it.CloseKeyboard()
	}
}

func (it *InputText) CanFocus() bool {
	return it.enabled
}

func (it *InputText) HandleInput(inputType input.InputType) bool {
	// Se o teclado estiver ativo, delegar para ele
	if it.keyboard != nil && it.keyboard.IsVisible() {
		return it.keyboard.HandleInput(inputType)
	}

	// Processar input do campo de texto
	switch inputType {
	case input.InputConfirm:
		it.OpenKeyboard()
		return true
	case input.InputLeft:
		it.MoveCursorLeft()
		return true
	case input.InputRight:
		it.MoveCursorRight()
		return true
	case input.InputBack:
		it.Backspace()
		return true
	}

	return false
}

// Métodos de manipulação de texto
func (it *InputText) SetText(text string) {
	if len(text) > it.MaxLength {
		text = text[:it.MaxLength]
	}
	it.Text = text
	it.CursorPos = len(text)
	if it.OnChange != nil {
		it.OnChange(it.Text)
	}
}

func (it *InputText) InsertText(text string) {
	if len(it.Text)+len(text) > it.MaxLength {
		availableSpace := it.MaxLength - len(it.Text)
		if availableSpace <= 0 {
			return
		}
		text = text[:availableSpace]
	}

	left := it.Text[:it.CursorPos]
	right := it.Text[it.CursorPos:]
	it.Text = left + text + right
	it.CursorPos += len(text)

	if it.OnChange != nil {
		it.OnChange(it.Text)
	}
}

func (it *InputText) Backspace() {
	if it.CursorPos > 0 {
		left := it.Text[:it.CursorPos-1]
		right := it.Text[it.CursorPos:]
		it.Text = left + right
		it.CursorPos--

		if it.OnChange != nil {
			it.OnChange(it.Text)
		}
	}
}

func (it *InputText) MoveCursorLeft() {
	if it.CursorPos > 0 {
		it.CursorPos--
	}
}

func (it *InputText) MoveCursorRight() {
	if it.CursorPos < len(it.Text) {
		it.CursorPos++
	}
}

func (it *InputText) Clear() {
	it.Text = ""
	it.CursorPos = 0
	if it.OnChange != nil {
		it.OnChange(it.Text)
	}
}

// Métodos do teclado virtual
func (it *InputText) OpenKeyboard() {
	if it.keyboard != nil {
		it.keyboard.Show()
	}
}

func (it *InputText) CloseKeyboard() {
	if it.keyboard != nil {
		log.Printf("InputText: Closing virtual keyboard for '%s'", it.ID)
		it.keyboard.Hide()
	}
}

func (it *InputText) onKeyboardInput(action KeyboardAction, value string) {
	switch action {
	case KeyboardActionCharacter:
		it.InsertText(value)
	case KeyboardActionBackspace:
		it.Backspace()
	case KeyboardActionSpace:
		it.InsertText(" ")
	case KeyboardActionEnter:
		if it.OnSubmit != nil {
			it.OnSubmit(it.Text)
		}
		it.CloseKeyboard()
	case KeyboardActionCancel:
		it.CloseKeyboard()
	case KeyboardActionClear:
		it.Clear()
	}
}

// Render renderiza o campo de texto
func (it *InputText) Render() {
	// Determinar estado atual baseado no foco
	var backgroundColor, textColor, borderColor clay.Color
	if it.focused {
		backgroundColor = it.Config.FocusedBackgroundColor
		textColor = it.Config.FocusedTextColor
		borderColor = it.Config.FocusedBorderColor
	} else {
		backgroundColor = it.Config.BackgroundColor
		textColor = it.Config.TextColor
		borderColor = it.Config.BorderColor
	}

	log.Printf("InputText: Rendering '%s' (focused: %t)", it.ID, it.focused)

	clay.UI()(clay.ElementDeclaration{
		Id: clay.ID(it.ID),
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  it.Width,
				Height: it.Height,
			},
			Padding: it.Config.Padding,
			ChildAlignment: clay.ChildAlignment{
				X: clay.ALIGN_X_LEFT,
				Y: clay.ALIGN_Y_CENTER,
			},
		},
		CornerRadius:    clay.CornerRadiusAll(it.Config.CornerRadius - float32(it.Config.BorderWidth)),
		BackgroundColor: backgroundColor,
		Border: clay.BorderElementConfig{
			Width: clay.BorderWidth{Left: it.Config.BorderWidth, Right: it.Config.BorderWidth, Top: it.Config.BorderWidth, Bottom: it.Config.BorderWidth},
			Color: borderColor,
		},
	}, func() {
		// Renderizar texto ou placeholder
		displayText := it.getDisplayText()
		if displayText == "" && it.Placeholder != "" {
			// Renderizar placeholder
			Text(it.Placeholder, it.Config.TextSize, it.Config.PlaceholderColor)
		} else {
			// Renderizar texto atual com cursor
			Text(displayText, it.Config.TextSize, textColor)
		}
	})

	// Renderizar teclado virtual se estiver visível
	if it.keyboard != nil && it.keyboard.IsVisible() {
		it.keyboard.Render()
	}

	log.Printf("InputText: Rendered '%s'", it.ID)
}

func (it *InputText) getDisplayText() string {
	if it.focused && it.showCursor {
		// Adicionar cursor visual na posição atual
		left := it.Text[:it.CursorPos]
		right := it.Text[it.CursorPos:]
		return left + "|" + right
	}
	return it.Text
}

// SetEnabled habilita/desabilita o campo de texto
func (it *InputText) SetEnabled(enabled bool) {
	it.enabled = enabled
}

// IsKeyboardVisible retorna se o teclado virtual está visível
func (it *InputText) IsKeyboardVisible() bool {
	return it.keyboard != nil && it.keyboard.IsVisible()
}

// GetKeyboard retorna o teclado virtual associado
func (it *InputText) GetKeyboard() *VirtualKeyboard {
	return it.keyboard
}

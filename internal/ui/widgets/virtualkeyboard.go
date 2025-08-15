package widgets

import (
	"log"
	"retroart-sdl2/internal/input"
	"retroart-sdl2/internal/theme"

	"github.com/TotallyGamerJet/clay"
)

// KeyboardAction define as ações que o teclado virtual pode executar
type KeyboardAction int

const (
	KeyboardActionCharacter KeyboardAction = iota
	KeyboardActionBackspace
	KeyboardActionSpace
	KeyboardActionEnter
	KeyboardActionCancel
	KeyboardActionClear
	KeyboardActionShift
	KeyboardActionSymbols
)

var keys = [][]Key{
	// Linha de números com símbolos correspondentes
	{
		{Display: "1", Value: "1", SymbolValue: "!", Action: KeyboardActionCharacter, Width: 0},
		{Display: "2", Value: "2", SymbolValue: "@", Action: KeyboardActionCharacter, Width: 0},
		{Display: "3", Value: "3", SymbolValue: "#", Action: KeyboardActionCharacter, Width: 0},
		{Display: "4", Value: "4", SymbolValue: "$", Action: KeyboardActionCharacter, Width: 0},
		{Display: "5", Value: "5", SymbolValue: "%", Action: KeyboardActionCharacter, Width: 0},
		{Display: "6", Value: "6", SymbolValue: "¨", Action: KeyboardActionCharacter, Width: 0},
		{Display: "7", Value: "7", SymbolValue: "&", Action: KeyboardActionCharacter, Width: 0},
		{Display: "8", Value: "8", SymbolValue: "*", Action: KeyboardActionCharacter, Width: 0},
		{Display: "9", Value: "9", SymbolValue: "(", Action: KeyboardActionCharacter, Width: 0},
		{Display: "0", Value: "0", SymbolValue: ")", Action: KeyboardActionCharacter, Width: 0},
		{Display: "-", Value: "-", SymbolValue: "_", Action: KeyboardActionCharacter, Width: 0},
		{Display: "=", Value: "=", SymbolValue: "+", Action: KeyboardActionCharacter, Width: 0},
	},
	// Primeira linha de letras
	{
		{Display: "Q", Value: "q", SymbolValue: "", Action: KeyboardActionCharacter, Width: 0},
		{Display: "W", Value: "w", SymbolValue: "", Action: KeyboardActionCharacter, Width: 0},
		{Display: "E", Value: "e", SymbolValue: "", Action: KeyboardActionCharacter, Width: 0},
		{Display: "R", Value: "r", SymbolValue: "", Action: KeyboardActionCharacter, Width: 0},
		{Display: "T", Value: "t", SymbolValue: "", Action: KeyboardActionCharacter, Width: 0},
		{Display: "Y", Value: "y", SymbolValue: "", Action: KeyboardActionCharacter, Width: 0},
		{Display: "U", Value: "u", SymbolValue: "", Action: KeyboardActionCharacter, Width: 0},
		{Display: "I", Value: "i", SymbolValue: "", Action: KeyboardActionCharacter, Width: 0},
		{Display: "O", Value: "o", SymbolValue: "", Action: KeyboardActionCharacter, Width: 0},
		{Display: "P", Value: "p", SymbolValue: "", Action: KeyboardActionCharacter, Width: 0},
		{Display: "[", Value: "[", SymbolValue: "{", Action: KeyboardActionCharacter, Width: 0},
		{Display: "]", Value: "]", SymbolValue: "}", Action: KeyboardActionCharacter, Width: 0},
		{Display: "\\", Value: "\\", SymbolValue: "|", Action: KeyboardActionCharacter, Width: 0},
	},
	// Segunda linha de letras
	{
		{Display: "A", Value: "a", SymbolValue: "", Action: KeyboardActionCharacter, Width: 0},
		{Display: "S", Value: "s", SymbolValue: "", Action: KeyboardActionCharacter, Width: 0},
		{Display: "D", Value: "d", SymbolValue: "", Action: KeyboardActionCharacter, Width: 0},
		{Display: "F", Value: "f", SymbolValue: "", Action: KeyboardActionCharacter, Width: 0},
		{Display: "G", Value: "g", SymbolValue: "", Action: KeyboardActionCharacter, Width: 0},
		{Display: "H", Value: "h", SymbolValue: "", Action: KeyboardActionCharacter, Width: 0},
		{Display: "J", Value: "j", SymbolValue: "", Action: KeyboardActionCharacter, Width: 0},
		{Display: "K", Value: "k", SymbolValue: "", Action: KeyboardActionCharacter, Width: 0},
		{Display: "L", Value: "l", SymbolValue: "", Action: KeyboardActionCharacter, Width: 0},
		{Display: ";", Value: ";", SymbolValue: ":", Action: KeyboardActionCharacter, Width: 0},
	},
	// Terceira linha de letras
	{

		{Display: "Z", Value: "z", SymbolValue: "", Action: KeyboardActionCharacter, Width: 0},
		{Display: "X", Value: "x", SymbolValue: "", Action: KeyboardActionCharacter, Width: 0},
		{Display: "C", Value: "c", SymbolValue: "", Action: KeyboardActionCharacter, Width: 0},
		{Display: "V", Value: "v", SymbolValue: "", Action: KeyboardActionCharacter, Width: 0},
		{Display: "B", Value: "b", SymbolValue: "", Action: KeyboardActionCharacter, Width: 0},
		{Display: "N", Value: "n", SymbolValue: "", Action: KeyboardActionCharacter, Width: 0},
		{Display: "M", Value: "m", SymbolValue: "", Action: KeyboardActionCharacter, Width: 0},
		{Display: ",", Value: ",", SymbolValue: "<", Action: KeyboardActionCharacter, Width: 0},
		{Display: ".", Value: ".", SymbolValue: ">", Action: KeyboardActionCharacter, Width: 0},
		{Display: "/", Value: "/", SymbolValue: "?", Action: KeyboardActionCharacter, Width: 0},
	},
	// Linha de ações
	{
		{Display: "Shift", Value: "", SymbolValue: "", Action: KeyboardActionShift, Width: 0},
		{Display: "Sym", Value: "", SymbolValue: "", Action: KeyboardActionSymbols, Width: 0},
		{Display: "Space", Value: " ", SymbolValue: " ", Action: KeyboardActionSpace, Width: 0},
		{Display: "Back", Value: "", SymbolValue: "", Action: KeyboardActionBackspace, Width: 0},
		{Display: "Enter", Value: "", SymbolValue: "", Action: KeyboardActionEnter, Width: 0},
		{Display: "Cancel", Value: "", SymbolValue: "", Action: KeyboardActionCancel, Width: 0},
	},
}

// Key representa uma tecla individual do teclado virtual
type Key struct {
	Display     string
	Value       string
	SymbolValue string // Value when in symbols mode
	Action      KeyboardAction
	Width       float32 // Para teclas especiais que podem ser mais largas
}

// VirtualKeyboard representa o teclado virtual na tela
type VirtualKeyboard struct {
	ID          string
	visible     bool
	currentRow  int
	currentCol  int
	keys        [][]Key
	config      theme.VirtualKeyboardStyle
	onInput     func(KeyboardAction, string)
	upperCase   bool
	symbolsMode bool
}

// NewVirtualKeyboard cria um novo teclado virtual
func NewVirtualKeyboard(id string, onInput func(KeyboardAction, string)) *VirtualKeyboard {
	vk := &VirtualKeyboard{
		ID:          id,
		visible:     false,
		currentRow:  0,
		currentCol:  0,
		config:      theme.GetVirtualKeyboardStyle(),
		onInput:     onInput,
		upperCase:   false,
		symbolsMode: false,
		keys:        keys,
	}

	return vk
}

// Show exibe o teclado virtual
func (vk *VirtualKeyboard) Show() {
	vk.visible = true
	vk.currentRow = 0
	vk.currentCol = 0
	log.Printf("VirtualKeyboard: Showing keyboard '%s'", vk.ID)
}

// Hide esconde o teclado virtual
func (vk *VirtualKeyboard) Hide() {
	vk.visible = false
	log.Printf("VirtualKeyboard: Hiding keyboard '%s'", vk.ID)
}

// IsVisible retorna se o teclado está visível
func (vk *VirtualKeyboard) IsVisible() bool {
	return vk.visible
}

// HandleInput processa input do teclado virtual
func (vk *VirtualKeyboard) HandleInput(inputType input.InputType) bool {
	if !vk.visible {
		return false
	}

	switch inputType {
	case input.InputUp:
		vk.navigateUp()
		return true
	case input.InputDown:
		vk.navigateDown()
		return true
	case input.InputLeft:
		vk.navigateLeft()
		return true
	case input.InputRight:
		vk.navigateRight()
		return true
	case input.InputConfirm:
		vk.activateCurrentKey()
		return true
	case input.InputBack:
		vk.onInput(KeyboardActionCancel, "")
		return true
	}

	return false
}

// Métodos de navegação
func (vk *VirtualKeyboard) navigateUp() {
	if vk.currentRow > 0 {
		vk.currentRow--
		vk.clampCurrentCol()
	}
}

func (vk *VirtualKeyboard) navigateDown() {
	if vk.currentRow < len(vk.keys)-1 {
		vk.currentRow++
		vk.clampCurrentCol()
	}
}

func (vk *VirtualKeyboard) navigateLeft() {
	if vk.currentCol > 0 {
		vk.currentCol--
	} else {
		// Wrap to end of row
		vk.currentCol = len(vk.keys[vk.currentRow]) - 1
	}
}

func (vk *VirtualKeyboard) navigateRight() {
	if vk.currentCol < len(vk.keys[vk.currentRow])-1 {
		vk.currentCol++
	} else {
		// Wrap to beginning of row
		vk.currentCol = 0
	}
}

func (vk *VirtualKeyboard) clampCurrentCol() {
	if vk.currentCol >= len(vk.keys[vk.currentRow]) {
		vk.currentCol = len(vk.keys[vk.currentRow]) - 1
	}
}

func (vk *VirtualKeyboard) activateCurrentKey() {
	if vk.currentRow < len(vk.keys) && vk.currentCol < len(vk.keys[vk.currentRow]) {
		key := vk.keys[vk.currentRow][vk.currentCol]

		switch key.Action {
		case KeyboardActionCharacter:
			var value string
			// Check for symbols mode first, then uppercase
			if vk.symbolsMode && key.SymbolValue != "" {
				value = key.SymbolValue
			} else if vk.upperCase {
				value = key.Display // Use uppercase display value
			} else {
				value = key.Value // Use lowercase value
			}
			vk.onInput(key.Action, value)
		case KeyboardActionShift:
			vk.ToggleCase()
			log.Printf("VirtualKeyboard: Toggled case, upperCase is now %v", vk.upperCase)
		case KeyboardActionSymbols:
			vk.ToggleSymbols()
			log.Printf("VirtualKeyboard: Toggled symbols, symbolsMode is now %v", vk.symbolsMode)
		default:
			vk.onInput(key.Action, key.Value)
		}

		log.Printf("VirtualKeyboard: Activated key '%s'", key.Display)
	}
}

// Render renderiza o teclado virtual
func (vk *VirtualKeyboard) Render() {
	if !vk.visible {
		return
	}
	// Container principal do teclado
	clay.UI()(clay.ElementDeclaration{
		Id: clay.ID(vk.ID + "_container"),
		Floating: clay.FloatingElementConfig{
			AttachTo: clay.ATTACH_TO_PARENT,
			Offset:   clay.Vector2{X: 0, Y: -50},
			AttachPoints: clay.FloatingAttachPoints{
				Parent:  clay.ATTACH_POINT_CENTER_CENTER,
				Element: clay.ATTACH_POINT_CENTER_TOP,
			},
		},
		Layout: clay.LayoutConfig{
			Padding:         vk.config.Padding,
			ChildGap:        vk.config.KeySpacing,
			LayoutDirection: clay.TOP_TO_BOTTOM,
			ChildAlignment: clay.ChildAlignment{
				X: clay.ALIGN_X_CENTER,
				Y: clay.ALIGN_Y_CENTER,
			},
		},
		CornerRadius:    clay.CornerRadiusAll(vk.config.CornerRadius),
		BackgroundColor: vk.config.BackgroundColor,
	}, func() {

		for rowIndex, row := range vk.keys {
			vk.renderKeyRow(rowIndex, row)
		}
	})
}

func (vk *VirtualKeyboard) renderKeyRow(rowIndex int, row []Key) {
	clay.UI()(clay.ElementDeclaration{
		Id: clay.ID(vk.ID + "_row_" + string(rune(rowIndex+'0'))),
		Layout: clay.LayoutConfig{
			ChildGap:        vk.config.KeySpacing,
			LayoutDirection: clay.LEFT_TO_RIGHT,
			ChildAlignment: clay.ChildAlignment{
				X: clay.ALIGN_X_CENTER,
				Y: clay.ALIGN_Y_CENTER,
			},
		},
	}, func() {

		for colIndex, key := range row {
			vk.renderKey(rowIndex, colIndex, key)
		}
	})
}

func (vk *VirtualKeyboard) renderKey(rowIndex, colIndex int, key Key) {
	isFocused := rowIndex == vk.currentRow && colIndex == vk.currentCol
	isShiftKey := key.Action == KeyboardActionShift
	isSymbolsKey := key.Action == KeyboardActionSymbols

	// Determinar cores baseadas no estado
	var backgroundColor, textColor clay.Color
	if isFocused {
		backgroundColor = vk.config.KeyButtonStyle.FocusedBackgroundColor
		textColor = vk.config.KeyButtonStyle.FocusedTextColor
	} else if isShiftKey && vk.upperCase {
		// Shift key is active (uppercase mode) - show it highlighted
		backgroundColor = vk.config.KeyButtonStyle.FocusedBackgroundColor
		textColor = vk.config.KeyButtonStyle.FocusedTextColor
	} else if isSymbolsKey && vk.symbolsMode {
		// Symbols key is active (symbols mode) - show it highlighted
		backgroundColor = vk.config.KeyButtonStyle.FocusedBackgroundColor
		textColor = vk.config.KeyButtonStyle.FocusedTextColor
	} else {
		backgroundColor = vk.config.KeyButtonStyle.BackgroundColor
		textColor = vk.config.KeyButtonStyle.TextColor
	}

	keyWidth := vk.config.KeyButtonStyle.Width
	if key.Width > 0 {
		keyWidth = key.Width
	}

	clay.UI()(clay.ElementDeclaration{
		Id: clay.ID(vk.ID + "_key_" + string(rune(rowIndex+'0')) + "_" + string(rune(colIndex+'0'))),
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  clay.SizingFit(keyWidth, 0),
				Height: clay.SizingFixed(vk.config.KeyButtonStyle.Height),
			},
			Padding: vk.config.KeyButtonStyle.Padding,
			ChildAlignment: clay.ChildAlignment{
				X: clay.ALIGN_X_CENTER,
				Y: clay.ALIGN_Y_CENTER,
			},
		},
		CornerRadius:    clay.CornerRadiusAll(vk.config.KeyButtonStyle.CornerRadius),
		BackgroundColor: backgroundColor,
	}, func() {

		displayText := key.Display

		// Determine what to display based on current mode
		if key.Action == KeyboardActionCharacter {
			if vk.symbolsMode && key.SymbolValue != "" {
				displayText = key.SymbolValue
			} else if vk.upperCase {
				displayText = key.Display // Already uppercase
			} else {
				displayText = key.Value // Lowercase
			}
		}

		Text(displayText, vk.config.KeyButtonStyle.FontSize, textColor)
	})
}

// ToggleCase alterna entre maiúsculas e minúsculas
func (vk *VirtualKeyboard) ToggleCase() {
	vk.upperCase = !vk.upperCase
}

// ToggleSymbols alterna entre números e símbolos
func (vk *VirtualKeyboard) ToggleSymbols() {
	vk.symbolsMode = !vk.symbolsMode
}

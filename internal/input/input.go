package input

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

// InputEvent representa um evento de input unificado
type InputEvent struct {
	Type    InputType
	Pressed bool
}

// InputType define os tipos de input suportados
type InputType int

const (
	InputUp InputType = iota
	InputDown
	InputLeft
	InputRight
	InputConfirm // A button / Enter
	InputBack    // B button / Escape
	InputMenu    // Start button
	InputSelect  // Select button
	InputX       // X button
	InputY       // Y button
)

// InputHandler define a interface para processadores de input
type InputHandler interface {
	Process() bool
	GetInputType() InputType
}

// InputProcessor processa diferentes tipos de input usando estratégias
type InputProcessor struct {
	channel              chan InputEvent
	lastDirectionalEvent map[any]uint64
	directionalThrottle  uint64
}

// NewInputProcessor cria um novo processador de input
func NewInputProcessor(channel chan InputEvent) *InputProcessor {
	return &InputProcessor{
		channel:              channel,
		lastDirectionalEvent: make(map[any]uint64),
		directionalThrottle:  150, // 150ms throttle
	}
}

// SendEvent envia um evento para o canal com fallback seguro
func (p *InputProcessor) SendEvent(inputType InputType) {
	select {
	case p.channel <- InputEvent{Type: inputType, Pressed: true}:
	default:
		log.Println("Input channel is full, ignoring event")
	}
}

// ProcessDirectionalInput processa input direcional com throttling
func (p *InputProcessor) ProcessDirectionalInput(key any, inputType InputType, isPressed bool) {
	if !isPressed {
		return
	}

	currentTime := sdl.GetTicks64()
	lastEventTime := p.lastDirectionalEvent[key]

	if lastEventTime == 0 || (currentTime-lastEventTime) >= p.directionalThrottle {
		p.SendEvent(inputType)
		p.lastDirectionalEvent[key] = currentTime
	}
}

// ProcessActionInput processa input de ação (edge-triggered)
func (p *InputProcessor) ProcessActionInput(inputType InputType, isPressed, wasPressed bool) {
	if isPressed && !wasPressed {
		p.SendEvent(inputType)
	}
}

var inputCh = make(chan InputEvent, 10)
var processor *InputProcessor

func Initialize() <-chan InputEvent {
	processor = NewInputProcessor(inputCh)
	go listenForKeyboardEvents()
	go listenForControllerEvents()
	return inputCh
}

// KeyboardHandler processa eventos de teclado
type KeyboardHandler struct {
	processor       *InputProcessor
	keyMappings     map[sdl.Scancode]InputType
	directionalKeys map[sdl.Scancode]bool
	previousState   []uint8
}

// NewKeyboardHandler cria um novo handler de teclado
func NewKeyboardHandler(processor *InputProcessor) *KeyboardHandler {
	return &KeyboardHandler{
		processor: processor,
		keyMappings: map[sdl.Scancode]InputType{
			sdl.SCANCODE_UP:     InputUp,
			sdl.SCANCODE_DOWN:   InputDown,
			sdl.SCANCODE_LEFT:   InputLeft,
			sdl.SCANCODE_RIGHT:  InputRight,
			sdl.SCANCODE_RETURN: InputConfirm,
			sdl.SCANCODE_SPACE:  InputConfirm,
			sdl.SCANCODE_ESCAPE: InputBack,
			sdl.SCANCODE_A:      InputConfirm, // Para TrimUI
			sdl.SCANCODE_B:      InputBack,    // Para TrimUI
			sdl.SCANCODE_X:      InputX,
			sdl.SCANCODE_Y:      InputY,
		},
		directionalKeys: map[sdl.Scancode]bool{
			sdl.SCANCODE_UP:    true,
			sdl.SCANCODE_DOWN:  true,
			sdl.SCANCODE_LEFT:  true,
			sdl.SCANCODE_RIGHT: true,
		},
		previousState: make([]uint8, sdl.NUM_SCANCODES),
	}
}

// ProcessInput processa input do teclado
func (h *KeyboardHandler) ProcessInput() {
	currentKeyState := sdl.GetKeyboardState()

	for scancode, inputType := range h.keyMappings {
		isPressed := currentKeyState[scancode] == 1
		wasPressed := h.previousState[scancode] == 1

		if h.directionalKeys[scancode] {
			h.processor.ProcessDirectionalInput(scancode, inputType, isPressed)
		} else {
			h.processor.ProcessActionInput(inputType, isPressed, wasPressed)
		}
	}

	copy(h.previousState, currentKeyState)
}

func listenForKeyboardEvents() {
	handler := NewKeyboardHandler(processor)

	for {
		handler.ProcessInput()
		sdl.Delay(16) // ~60 FPS para polling de input
	}
}

// ControllerHandler processa eventos de game controller
type ControllerHandler struct {
	processor           *InputProcessor
	controller          *sdl.GameController
	buttonMappings      map[sdl.GameControllerButton]InputType
	directionalButtons  map[sdl.GameControllerButton]bool
	previousButtonState map[sdl.GameControllerButton]bool
}

// NewControllerHandler cria um novo handler de controller
func NewControllerHandler(processor *InputProcessor, controller *sdl.GameController) *ControllerHandler {
	handler := &ControllerHandler{
		processor:  processor,
		controller: controller,
		buttonMappings: map[sdl.GameControllerButton]InputType{
			sdl.CONTROLLER_BUTTON_DPAD_UP:    InputUp,
			sdl.CONTROLLER_BUTTON_DPAD_DOWN:  InputDown,
			sdl.CONTROLLER_BUTTON_DPAD_LEFT:  InputLeft,
			sdl.CONTROLLER_BUTTON_DPAD_RIGHT: InputRight,
			sdl.CONTROLLER_BUTTON_A:          InputConfirm,
			sdl.CONTROLLER_BUTTON_B:          InputBack,
			sdl.CONTROLLER_BUTTON_X:          InputX,
			sdl.CONTROLLER_BUTTON_Y:          InputY,
			sdl.CONTROLLER_BUTTON_START:      InputMenu,
			sdl.CONTROLLER_BUTTON_BACK:       InputSelect,
		},
		directionalButtons: map[sdl.GameControllerButton]bool{
			sdl.CONTROLLER_BUTTON_DPAD_UP:    true,
			sdl.CONTROLLER_BUTTON_DPAD_DOWN:  true,
			sdl.CONTROLLER_BUTTON_DPAD_LEFT:  true,
			sdl.CONTROLLER_BUTTON_DPAD_RIGHT: true,
		},
		previousButtonState: make(map[sdl.GameControllerButton]bool),
	}

	return handler
}

// ProcessInput processa input do controller
func (h *ControllerHandler) ProcessInput() {
	// Processar botões
	for button, inputType := range h.buttonMappings {
		isPressed := h.controller.Button(button) == 1
		wasPressed := h.previousButtonState[button]

		if h.directionalButtons[button] {
			h.processor.ProcessDirectionalInput(button, inputType, isPressed)
		} else {
			h.processor.ProcessActionInput(inputType, isPressed, wasPressed)
		}

		h.previousButtonState[button] = isPressed
	}
}

// listenForControllerEvents processa eventos de Game Controller
func listenForControllerEvents() {
	controller := openController()
	defer func() {
		if controller != nil {
			controller.Close()
		}
	}()

	if controller == nil {
		return
	}

	handler := NewControllerHandler(processor, controller)

	for {
		handler.ProcessInput()
		sdl.Delay(100) // ~60 FPS para polling de controller
	}
}

// openController procura e abre o primeiro Game Controller disponível
func openController() *sdl.GameController {
	// Primeiro, verificar se há joysticks conectados
	numJoysticks := sdl.NumJoysticks()
	if numJoysticks == 0 {
		return nil
	}

	// Procurar por Game Controllers
	for i := range numJoysticks {
		if sdl.IsGameController(i) {
			controller := sdl.GameControllerOpen(i)
			if controller != nil {
				name := controller.Name()
				// Log do controller encontrado (útil para debug no TrimUI)
				println("Controller encontrado:", name)
				return controller
			}
		}
	}

	return nil
}

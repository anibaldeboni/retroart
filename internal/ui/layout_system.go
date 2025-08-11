package ui

import (
	"log"

	"github.com/TotallyGamerJet/clay"
	"github.com/veandco/go-sdl2/sdl"
)

// ClayLayoutSystem gerencia o sistema de layout Clay
type ClayLayoutSystem struct {
	renderer *sdl.Renderer
	enabled  bool
	isActive bool
}

// NewClayLayoutSystem cria um novo sistema de layout Clay
func NewClayLayoutSystem(renderer *sdl.Renderer) *ClayLayoutSystem {
	cls := &ClayLayoutSystem{
		renderer: renderer,
		enabled:  false,
		isActive: false,
	}

	// Verificar se Clay foi inicializado globalmente
	if globalClayContext != nil {
		cls.enabled = true
		log.Println("ClayLayoutSystem created successfully")
	} else {
		log.Println("Clay not initialized globally - ClayLayoutSystem disabled")
	}

	return cls
}

// BeginLayout inicia um novo frame de layout seguindo o padrão do videodemo.go
func (cls *ClayLayoutSystem) BeginLayout() {
	if !cls.enabled {
		log.Println("Clay not enabled, skipping BeginLayout")
		return
	}

	// Verificar se o contexto atual existe
	currentContext := clay.GetCurrentContext()
	if currentContext == nil {
		log.Println("Clay current context is nil, trying to set global context")
		if globalClayContext != nil {
			clay.SetCurrentContext(globalClayContext)
			currentContext = clay.GetCurrentContext()
			if currentContext == nil {
				log.Println("Failed to set current context, aborting BeginLayout")
				return
			}
		} else {
			log.Println("Global Clay context is nil, aborting BeginLayout")
			return
		}
	}

	// Resetar arena para o offset correto como no clay.h:2155
	globalClayArena.NextAllocation = arenaResetOffset
	log.Printf("Arena reset to offset: %d", arenaResetOffset)

	cls.isActive = true
	log.Println("Calling clay.BeginLayout()")
	clay.BeginLayout()
	log.Println("clay.BeginLayout() completed")
}

// EndLayout finaliza o layout e retorna os comandos de renderização
func (cls *ClayLayoutSystem) EndLayout() clay.RenderCommandArray {
	if !cls.enabled || !cls.isActive {
		log.Println("Clay not enabled or not active, returning empty RenderCommandArray")
		return clay.RenderCommandArray{}
	}
	cls.isActive = false
	log.Println("Calling clay.EndLayout()")
	commands := clay.EndLayout()
	log.Printf("clay.EndLayout() completed, got %d commands", commands.Length)
	return commands
}

// Render executa o ciclo completo de renderização
func (cls *ClayLayoutSystem) Render() {
	if !cls.enabled {
		return
	}

	commands := cls.EndLayout()

	err := cls.RenderClayCommands(commands)
	if err != nil {
		log.Printf("Error rendering Clay commands: %v", err)
	}
}

// RenderClayCommands renderiza os comandos Clay no renderer SDL2
func (cls *ClayLayoutSystem) RenderClayCommands(commands clay.RenderCommandArray) error {
	if !cls.enabled {
		return nil
	}

	log.Printf("RenderClayCommands called with %d commands", commands.Length)

	// Iterar pelos comandos de renderização
	for i := int32(0); i < commands.Length; i++ {
		command := clay.RenderCommandArray_Get(&commands, i)

		switch command.CommandType {
		case clay.RENDER_COMMAND_TYPE_RECTANGLE:
			config := &command.RenderData.Rectangle
			log.Printf("Clay Rectangle Command %d: BoundingBox X=%.2f Y=%.2f W=%.2f H=%.2f, Color R=%.0f G=%.0f B=%.0f A=%.0f",
				i, command.BoundingBox.X, command.BoundingBox.Y, command.BoundingBox.Width, command.BoundingBox.Height,
				config.BackgroundColor.R, config.BackgroundColor.G, config.BackgroundColor.B, config.BackgroundColor.A)
			err := cls.renderRectangle(command)
			if err != nil {
				log.Printf("Error rendering rectangle: %v", err)
				return err
			}
		case clay.RENDER_COMMAND_TYPE_TEXT:
			err := cls.renderText(command)
			if err != nil {
				log.Printf("Error rendering text: %v", err)
				return err
			}
		case clay.RENDER_COMMAND_TYPE_BORDER:
			log.Printf("Rendering border command")
			err := cls.renderBorder(command)
			if err != nil {
				return err
			}
		case clay.RENDER_COMMAND_TYPE_SCISSOR_START:
			err := cls.renderScissorStart(command)
			if err != nil {
				return err
			}
		case clay.RENDER_COMMAND_TYPE_SCISSOR_END:
			err := cls.renderScissorEnd()
			if err != nil {
				return err
			}
		default:
			log.Printf("Command not implemented: %v", command.CommandType)
		}
	}

	return nil
}

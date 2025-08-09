package main

import (
	"fmt"
	"os"

	"retroart-sdl2/internal/app"
)

func main() {
	application := app.New()
	if err := application.Init(); err != nil {
		fmt.Printf("Erro ao inicializar aplicação: %v\n", err)
		os.Exit(1)
	}
	defer application.Cleanup()

	application.Run()
}

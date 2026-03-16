// cmd/main.go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/franciscozamorau/osmi-gateway/internal/config"
	"github.com/franciscozamorau/osmi-gateway/internal/server"
)

func main() {
	// 1. Cargar configuración
	cfg := config.Load()

	// 2. Crear servidor
	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Error al crear servidor: %v", err)
	}

	// 3. Manejo de señales para shutdown graceful
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Apagando servidor...")
		if err := srv.Stop(); err != nil {
			log.Printf("Error al apagar: %v", err)
		}
	}()

	// 4. Iniciar servidor
	log.Printf("Gateway iniciado en puerto %s", cfg.HTTPPort)
	if err := srv.Start(); err != nil {
		log.Fatalf("Error en servidor: %v", err)
	}
}

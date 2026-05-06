package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"solar-backend/internal/handlers"
)

func main() {
	if err := godotenv.Load(); err != nil {
		// Jangan hentikan program jika .env tidak ada (misal saat di produksi/Vercel)
		// tapi beri peringatan di log.
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan env sistem")
	}

	router := chi.NewRouter()
	router.Use(cors.AllowAll().Handler)
	router.Post("/api/v1/simulation", handlers.SimulateHandler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("server listening on :8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}
	log.Println("server stopped")
}

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"orches-saving/tracing"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func main() {
	if err := loadConfig(); err != nil {
		log.Fatalf("error from loading config: %v", err)
	}

	cleanup, err := initTracing()
	if err != nil {
		log.Fatalf("error from init tracer: %v", err)
	}
	defer cleanup(context.Background())

	if err := startServer(); err != nil {
		log.Fatalf("start server failed: %v", err)
	}
}

func loadConfig() error {
	// env
	if err := godotenv.Load("config/.env"); err != nil {
		log.Println("error from reading .env file")
	}

	// config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if configFile, err := os.ReadFile("./config/config.yml"); err != nil {
		log.Fatalf("error from reading config file: %v", err)
	} else {
		expandedConfig := os.ExpandEnv(string(configFile))
		viper.ReadConfig(bytes.NewBufferString(expandedConfig))
	}

	// log
	if configJSON, err := json.MarshalIndent(viper.AllSettings(), "", "  "); err != nil {
		log.Fatalf("error from marshalling config: %v", err)
	} else {
		log.Printf("config:\n%s", string(configJSON))
	}

	return nil
}

func initTracing() (func(context.Context) error, error) {
	serviceName := viper.GetString("app.name")
	jaegerEndpoint := viper.GetString("tracing.endpoint")

	cleanup, err := tracing.InitTracerAbstraced(
		serviceName,
		viper.GetString("app.version"),
		viper.GetString("app.environment"),
		jaegerEndpoint,
	)
	if err != nil {
		return nil, fmt.Errorf("error from init tracer: %w", err)
	}
	log.Printf("init tracer serviceName: %s, jaegerEndpoint: %s", serviceName, jaegerEndpoint)
	return cleanup, nil
}

func startServer() error {
	port := viper.GetInt("app.port")

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/test", testhandler)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      http.DefaultServeMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("start server at port %d", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen and serve failed: %v", err)
		}
	}()

	return gracefulShutdown(server)
}

func gracefulShutdown(server *http.Server) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server force to shutdown: %w", err)
	}
	log.Println("server stopped gracefully")
	return nil
}

// test
func testhandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("test success"))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

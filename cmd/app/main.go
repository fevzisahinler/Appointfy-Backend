package main

import (
	"fmt"
	"github.com/fevzisahinler/Appointfy-Backend/config"
	"github.com/fevzisahinler/Appointfy-Backend/db"
	"github.com/fevzisahinler/Appointfy-Backend/http/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	loc, err := time.LoadLocation("Europe/Istanbul")
	if err != nil {
		fmt.Printf("Timezone couldn't load: %v\n", err)
	}
	time.Local = loc

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v", err)
	}
	app := fiber.New()
	app.Use(cors.New())

	if err := db.ConnectDatabase(cfg); err != nil {
		fmt.Println("Database connection failed")
	}

	routes.AuthRoutes(app)

	port := ":4000"
	go func() {
		if err := app.Listen(port); err != nil {
			fmt.Println("Server failed to start")
		}
	}()
	waitForShutdown()

}

func waitForShutdown() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	fmt.Println("Shutting down server...")
}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/The-Flash/code-racer/internal/runtime_manager"
	"github.com/docker/docker/client"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var (
	manifestPtr   = flag.String("f", "", "Path to manifest file")
	mountPointPtr = flag.String("m", "", "Path to mount point")
)

func main() {
	flag.Parse()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal("could not connect to docker ", err)
	}
	defer cli.Close()

	go runtime_manager.RunScaler(cli, *manifestPtr, *mountPointPtr)

	port := os.Getenv("PORT")
	app := fiber.New()

	if port == "" {
		port = "8000"
	}

	app.Listen(fmt.Sprintf(":%s", port))

}

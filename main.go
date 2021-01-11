package main

import (
	c "github.com/aldinokemal/go-whatsapp-web/config"
	r "github.com/aldinokemal/go-whatsapp-web/routers"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	_ = os.Setenv("TZ", "Asia/Jakarta")

	c.DBOpen() // database
	c.SetupEnv()

	app := r.Routers()
	_ = app.Run(c.AppPort)
}

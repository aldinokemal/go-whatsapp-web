package main

import (
	c "github.com/aldinokemal/go-whatsapp-web/config"
	r "github.com/aldinokemal/go-whatsapp-web/routers"
	"os"
)

func main() {
	_ = os.Setenv("TZ", "Asia/Jakarta")

	c.DBOpen() // database
	c.SetupEnv()

	app := r.Routers()
	_ = app.Run(":8000")
}

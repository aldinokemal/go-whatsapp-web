package config

import (
	"github.com/Rhymen/go-whatsapp"
	"os"
	"time"
)

var (
	AppPort string

	PathQrCode    string
	PathWaImage string
	PathWaSession string
	BasicAuthUser string
	BasicAuthPswd string

	WhatsappConfig whatsapp.Options
)

func SetupEnv() {
	AppPort = ":" + os.Getenv("APP_PORT")

	PathQrCode = "statics/images/qrcode/"
	PathWaImage = "statics/images/wa_images/"
	PathWaSession = "storage/session/"

	BasicAuthUser = os.Getenv("BASIC_AUTH_USER")
	BasicAuthPswd = os.Getenv("BASIC_AUTH_PSWD")

	//return
	WhatsappConfig = whatsapp.Options{
		Timeout:         5 * time.Second,
		Handler:         nil,
		ShortClientName: "X-WA",
		LongClientName:  "Whatsapp By Aldino Kemal",
		ClientVersion:   "1.0",
		Store:           nil,
	}
}

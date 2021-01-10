package config

import (
	"github.com/Rhymen/go-whatsapp"
	"time"
)

var (
	PathQrCode    string
	PathWaSession string
	BasicAuthUser string
	BasicAuthPswd string

	WhatsappConfig whatsapp.Options
)

func SetupEnv() {
	PathQrCode = "statics/images/qrcode/"
	PathWaSession = "storage/session/"

	BasicAuthUser = "dev"
	BasicAuthPswd = "2104"

	//return
	WhatsappConfig = whatsapp.Options{
		Timeout:         5 * time.Second,
		Handler:         nil,
		ShortClientName: "X-WA",
		LongClientName:  "Whatsapp By Kemal",
		ClientVersion:   "1.0",
		Store:           nil,
	}
}

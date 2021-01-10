package config

var (
	PathQrCode    string
	PathWaSession string
)

func SetupEnv() {
	PathQrCode = "statics/"
	PathWaSession = "storage/session/"
}

package helpers

import (
	"encoding/gob"
	"github.com/Rhymen/go-whatsapp"
	c "github.com/aldinokemal/go-whatsapp-web/config"
	"os"
)

func ReadSession(phone string) (whatsapp.Session, error) {
	session := whatsapp.Session{}
	//fmt.Println(os.TempDir() + "whatsappSession.gob")
	file, err := os.Open(c.PathWaSession + phone + "Session.gob")
	if err != nil {
		return session, err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&session)
	if err != nil {
		return session, err
	}
	return session, nil
}

func WriteSession(session whatsapp.Session, phone string) error {
	file, err := os.Create(c.PathWaSession + phone + "Session.gob")
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(session)
	if err != nil {
		return err
	}
	return nil
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

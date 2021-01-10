package handler

import (
	"fmt"
	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Rhymen/go-whatsapp"
	c "github.com/aldinokemal/go-whatsapp-web/config"
	h "github.com/aldinokemal/go-whatsapp-web/helpers"
	"github.com/aldinokemal/go-whatsapp-web/structs"
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"net/http"
	"os"
	"strings"
	"time"
)

var qr chan string
var qrName string

func Authenticated(g *gin.Context) {
	var validation structs.ValidateLogin
	if err := g.ShouldBind(&validation); err != nil {
		h.RespondJSON(g, http.StatusBadRequest, strings.Split(err.Error(), "\n"), "Parameter tidak valid")
	} else {
		x := c.TableAccount{AccPhone: validation.Phone}
		data := x.FindByPhone()
		if data.AccID != 0 {
			if h.FileExists(c.PathWaSession + data.AccSessionName) {
				results := map[string]string{
					"message": "This number already logged in",
				}
				h.RespondJSON(g, http.StatusOK, results)
				return
			} else {
				_ = x.DelByID()
			}
		}

		wac, err := whatsapp.NewConn(5 * time.Second) // Create connection to whatsapp (belum masuk ke proseds login/cek sesion)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error creating connection: %v\n", err)
			h.RespondJSON(g, http.StatusBadRequest, nil, fmt.Sprintf("error creating connection: %v\n", err.Error()))
			return
		} else {
			qrName = "qr_" + validation.Phone + time.Now().String() + ".png"
			sessionName := validation.Phone + "Session.gob"
			err = LoginViaWeb(wac, g, validation.Phone)
			if err != nil {
				fmt.Println("terjadi kesalahan: ", err.Error())
				//if h.FileExists("statics/" + qrName) {
				//
				//} else {
				_, _ = fmt.Fprintf(os.Stderr, "error logging in: %v\n", err)
				fmt.Println(c.PathWaSession + sessionName)
				err = os.Remove(c.PathWaSession + sessionName)
				if err != nil {
					fmt.Println(err.Error())
				}

				Authenticated(g)
				//}
			} else {
				results := map[string]string{
					"message": "Your QR is generated, please scan it",
					"image":   qrName,
				}
				h.RespondJSON(g, http.StatusOK, results)
				return
			}
		}

	}
}

func Login(wac *whatsapp.Conn) error {
	//load saved session
	session, err := h.ReadSession("1298379123")
	if err == nil {
		//restore session
		session, err = wac.RestoreWithSession(session)
		if err != nil {
			return fmt.Errorf("restoring failed: %v\n", err)
		}
	} else {
		//no saved session -> regular login
		qr := make(chan string)
		go func() {
			terminal := qrcodeTerminal.New()
			terminal.Get(<-qr).Print()
		}()
		session, err = wac.Login(qr)
		if err != nil {
			return fmt.Errorf("error during login: %v\n", err)
		}
	}

	//save session
	err = h.WriteSession(session, "1298379123")
	if err != nil {
		return fmt.Errorf("error saving session: %v\n", err)
	}
	return nil
}

func LoginViaWeb(wac *whatsapp.Conn, g *gin.Context, phone string) error {
	//load saved session
	session, err := h.ReadSession(phone)
	fmt.Println("current session ", session)
	if err == nil {
		//restore session
		session, err = wac.RestoreWithSession(session)
		if err != nil {
			return fmt.Errorf("restoring failed: %v\n", err)
		} else {
			// Save Session
			err = h.WriteSession(session, phone)
			if err != nil {
				return fmt.Errorf("error saving session: %v\n", err)
			}
		}
	} else {
		fmt.Println("generate session (png)", session)
		//no saved session -> regular login
		qr = make(chan string)
		go func() {
			fmt.Println("generate session images", session)
			fmt.Println(c.PathQrCode + qrName)
			err = qrcode.WriteFile(<-qr, qrcode.Medium, 512, c.PathQrCode+qrName)
			if err != nil {
				fmt.Println("salah saat generate qr: ", err.Error())
			} else {
				fmt.Println("Stop looping", session)
			}
		}()

		go func() {
			session, err = wac.Login(qr)
			if err != nil {
				_ = fmt.Errorf("error during login: %v\n", err)
			}

			// Save Session
			err = h.WriteSession(session, phone)
			if err != nil {
				_ = fmt.Errorf("error saving session: %v\n", err)
			} else {
				account := c.TableAccount{
					AccPhone:       phone,
					AccSessionName: phone + "Session.gob",
				}
				account.InsertAccount()
			}
		}()
	}

	return nil
}

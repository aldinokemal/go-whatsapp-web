package handler

import (
	"fmt"
	"github.com/Rhymen/go-whatsapp"
	h "github.com/aldinokemal/go-whatsapp-web/helpers"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
	"time"
)

type ValidationSendWA struct {
	To      string `binding:"required" json:"to" form:"to"`
	Message string `binding:"required" json:"message" form:"message"`
}

func SendMessage(c *gin.Context) {
	var validation ValidationSendWA
	if err := c.ShouldBind(&validation); err != nil {
		h.RespondJSON(c, http.StatusBadRequest, strings.Split(err.Error(), "\n"), "Parameter tidak valid")
	} else {
		//create new WhatsApp connection
		wac, err := whatsapp.NewConn(5 * time.Second)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error creating connection: %v\n", err)
			return
		}

		err = LoginViaWeb(wac, c, "qrcodewa.png")
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error logging in: %v\n", err)
			fmt.Println(os.TempDir() + "whatsappSession.gob")
			err = os.Remove(os.TempDir() + "whatsappSession.gob")
			if err != nil {
				fmt.Println(err.Error())
			}

			SendMessage(c)
		}

		<-time.After(3 * time.Second)

		msg := whatsapp.TextMessage{
			Info: whatsapp.MessageInfo{
				RemoteJid: fmt.Sprintf("%s@s.whatsapp.net", validation.To),
			},
			Text: validation.Message,
		}

		msgId, err := wac.Send(msg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Terjadi kesalahan")
			fmt.Fprintf(os.Stderr, "error sending message: %v", err)
			os.Exit(1)
		} else {
			fmt.Println("Message Sent -> ID : " + msgId)
			c.JSON(http.StatusOK, "message sent")
		}
	}
}

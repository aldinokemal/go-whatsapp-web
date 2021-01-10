package handler

import (
	"fmt"
	"github.com/Rhymen/go-whatsapp"
	c "github.com/aldinokemal/go-whatsapp-web/config"
	h "github.com/aldinokemal/go-whatsapp-web/helpers"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
	"time"
)

type ValidationSendWA struct {
	From    string `binding:"required" json:"from" form:"from"`
	To      string `binding:"required" json:"to" form:"to"`
	Message string `binding:"required" json:"message" form:"message"`
}

func SendMessage(g *gin.Context) {
	var validation ValidationSendWA
	if err := g.ShouldBind(&validation); err != nil {
		h.RespondJSON(g, http.StatusBadRequest, strings.Split(err.Error(), "\n"), "Parameter tidak valid")
	} else {
		//create new WhatsApp connection
		wac, err := whatsapp.NewConnWithOptions(&c.WhatsappConfig)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error creating connection: %v\n", err)
			h.RespondJSON(g, http.StatusBadRequest, nil, fmt.Sprintf("error creating connection: %v\n", err.Error()))
			return
		}

		sessionName := validation.From + "Session.gob"
		err = LoginViaWeb(wac, validation.From)
		if err != nil {
			err = os.Remove(c.PathWaSession + sessionName)
			if err != nil {
				fmt.Println(err.Error())
			}

			SendMessage(g)
		} else {
			<-time.After(3 * time.Second)

			msg := whatsapp.TextMessage{
				Info: whatsapp.MessageInfo{
					RemoteJid: fmt.Sprintf("%s@s.whatsapp.net", validation.To),
				},
				Text: validation.Message,
			}

			msgId, err := wac.Send(msg)
			if err != nil {
				h.RespondJSON(g, http.StatusInternalServerError, err.Error(), "terjadi kesalahan")
				return
			} else {
				h.RespondJSON(g, http.StatusOK, msgId, "message sent")
			}
		}
	}
}

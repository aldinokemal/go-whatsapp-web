package handler

import (
	c "github.com/aldinokemal/go-whatsapp-web/config"
	h "github.com/aldinokemal/go-whatsapp-web/helpers"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Home(g *gin.Context) {
	g.HTML(http.StatusOK, "home.html", gin.H{
		"title": "WhatsApp Web by Aldino Kemal",
		"ba_u":  c.BasicAuthUser,
		"ba_p":  c.BasicAuthPswd,
	})
}

func GetAccount(g *gin.Context) {
	var (
		dataAccount c.TableAccount
		response    []map[string]interface{}
	)
	data := dataAccount.FindAll()
	for _, d := range data {
		x := map[string]interface{}{
			"id":       d.AccID,
			"phone":    d.AccPhone,
			"phone_id": strings.Split(d.AccWaID.String, "@")[0],
			"created":  d.AccCreatedAt.Time.Format("January 2, 2006 15:04:05"),
		}
		response = append(response, x)
	}

	h.RespondJSON(g, http.StatusOK, response, "data ditemukan")
}

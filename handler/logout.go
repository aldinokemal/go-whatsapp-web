package handler

import (
	h "github.com/aldinokemal/go-whatsapp-web/helpers"
	"github.com/aldinokemal/go-whatsapp-web/structs"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthLogout(g *gin.Context) {
	var validation structs.ValidateLogin
	if err := g.ShouldBind(&validation); err != nil {
		h.RespondJSON(g, http.StatusBadRequest, strings.Split(err.Error(), "\n"), "Parameter tidak valid")
	} else {

	}
}

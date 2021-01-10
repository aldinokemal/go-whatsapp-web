package handler

import (
	"fmt"
	c "github.com/aldinokemal/go-whatsapp-web/config"
	h "github.com/aldinokemal/go-whatsapp-web/helpers"
	"github.com/aldinokemal/go-whatsapp-web/structs"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

func AuthLogout(g *gin.Context) {
	var validation structs.ValidateLogin
	if err := g.ShouldBind(&validation); err != nil {
		h.RespondJSON(g, http.StatusBadRequest, strings.Split(err.Error(), "\n"), "Parameter tidak valid")
	} else {
		x := c.TableAccount{AccPhone: validation.Phone}
		data := x.FindByPhone()
		if data.AccID != 0 {
			BasePath, _ := os.Getwd()
			fmt.Println(data)
			fmt.Print(BasePath + "/" + c.PathWaSession + data.AccSessionName.String)
			_ = x.DelByPhone()
			_ = os.Remove(c.PathWaSession + data.AccSessionName.String)
		}

		h.RespondJSON(g, http.StatusOK, nil, "Logout success")
	}
}

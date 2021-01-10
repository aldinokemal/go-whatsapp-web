package middleware

import (
	"fmt"
	"github.com/Rhymen/go-whatsapp"
	"github.com/aldinokemal/go-whatsapp-web/handler"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

func Authenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		wac, err := whatsapp.NewConn(5 * time.Second)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error creating connection: %v\n", err)
			return
		}
		err = handler.Login(wac)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error logging in: %v\n", err)
			fmt.Println(os.TempDir() + "whatsappSession.gob")
			err = os.Remove(os.TempDir() + "whatsappSession.gob")
			if err != nil {
				fmt.Println(err.Error())
			}

			Authenticated()
			c.Abort()
			return
		}
	}
}

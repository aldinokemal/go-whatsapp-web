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

type ValidationSendImage struct {
	AppID   string `binding:"required" json:"from" form:"from"`
	To      string `binding:"required" json:"to" form:"to"`
	Caption string `json:"caption" form:"caption"`
}

func SendImage(g *gin.Context) {
	var validation ValidationSendImage
	if err := g.ShouldBind(&validation); err != nil {
		h.RespondJSON(g, http.StatusBadRequest, strings.Split(err.Error(), "\n"), "parameter tidak valid")
		return
	} else {
		// Cek apakah upload gambar tersedia
		file, err := g.FormFile("image")
		if err != nil {
			h.RespondJSON(g, http.StatusBadRequest, file, err.Error())
			return
		} else {
			allowedMime := map[string]bool{
				"image/jpeg": true,
				"image/jpg":  true,
				"image/gif":  true,
				"image/png":  true,
			}

			if file.Size > 1000000 {
				h.RespondJSON(g, http.StatusBadRequest, file, "file max 1MB")
				return
			}
			valid := allowedMime[file.Header.Get("Content-Type")] && file.Size < 500000

			if valid {
				imagePath := c.PathWaImage + time.Now().Format("20060102_150405_") + file.Filename
				err := g.SaveUploadedFile(file, imagePath)
				if err != nil {
					h.RespondJSON(g, http.StatusBadRequest, nil, err.Error())
					return
				} else {
					// Process Send Images
					x := c.TableAccount{AccPhone: validation.AppID}
					data := x.FindByPhone()
					if data.AccID != 0 {
						if h.FileExists(c.PathWaSession + data.AccSessionName.String) {
							//create new WhatsApp connection
							wac, err := whatsapp.NewConnWithOptions(&c.WhatsappConfig)
							if err != nil {
								_, _ = fmt.Fprintf(os.Stderr, "error creating connection: %v\n", err)
								h.RespondJSON(g, http.StatusBadRequest, nil, fmt.Sprintf("error creating connection: %v\n", err.Error()))
								return
							}

							sessionName := validation.AppID + "Session.gob"
							err = LoginViaWeb(wac, validation.AppID)
							if err != nil {
								err = os.Remove(c.PathWaSession + sessionName)
								if err != nil {
									fmt.Println(err.Error())
								}
							} else {
								<-time.After(3 * time.Second)

								img, err := os.Open(imagePath)
								if err != nil {
									_, _ = fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
									os.Exit(1)
								}

								msg := whatsapp.ImageMessage{
									Info: whatsapp.MessageInfo{
										RemoteJid: fmt.Sprintf("%s@s.whatsapp.net", validation.To),
									},
									Type:    "image/jpeg",
									Caption: validation.Caption,
									Content: img,
								}

								msgId, err := wac.Send(msg)
								if err != nil {
									h.RespondJSON(g, http.StatusInternalServerError, err.Error(), "terjadi kesalahan")
									return
								} else {
									// Hapus gambar yang dikirim
									err = os.Remove(imagePath)
									if err != nil {
										fmt.Println(err.Error())
									}

									h.RespondJSON(g, http.StatusOK, msgId, "message sent")
									return
								}
							}
						} else {
							_ = x.DelByPhone()
							h.RespondJSON(g, http.StatusInternalServerError, nil, "mohon untuk login ulang")
							return
						}
					} else {
						h.RespondJSON(g, http.StatusInternalServerError, nil, "nomor ini belum login")
						return
					}
				}
			} else {
				h.RespondJSON(g, http.StatusInternalServerError, nil, "validasi image tidak valid")
				return
			}
		}

	}
}
